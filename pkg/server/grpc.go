package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
	"golang.org/x/crypto/bcrypt"
	"time"
)
import "context"

type GRPCServer struct{
	repo *Repo
	dao *db.DAO
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	tx, err :=  s.dao.Conn.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := s.buildUser(tx, request.Password, request.PasswordConfirmation, request.Email)

	sessionUUID := uuid.New()
	sessionToken, json, furthestExpiration, err := s.buildToken(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return nil, err
	}


	session, err := s.buildSession(tx, sessionUUID, user.id, sessionToken, furthestExpiration)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = s.buildScopeGroupings(tx, request.ScopeGroupings, session.id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.CreateUserResponse{User:user.ConvertToProtobuff(), Session:session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) GetUser(_ context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	tx, err :=  s.dao.Conn.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	switch ident := request.Identifier.(type) {
	case *proto.GetUserRequest_Email:
		panic("implement me")
	case *proto.GetUserRequest_Uuid:
		user, err := s.repo.GetUserWithUUID(tx, ident.Uuid)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &proto.GetUserResponse{User:user.ConvertToProtobuff()}, nil
	}

	return nil, errors.New("unknown user identifier")
}

func (s *GRPCServer) CreateGuestUser(context.Context, *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) CreatePasswordResetToken(context.Context, *proto.CreatePasswordResetTokenResponse) (*proto.CreatePasswordResetTokenResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) UpdateUserPassword(context.Context, *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) CreateSession(context.Context, *proto.CreatePasswordResetTokenRequest) (*proto.CreatePasswordResetTokenResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) CreateSessionRevoke(context.Context, *proto.CreateSessionRevokeRequest) (*proto.CreateSessionRevokeResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) GetSession(context.Context, *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	panic("implement me")
}

// Builder Functions

func (s *GRPCServer) buildUser(tx *sql.Tx, password string, passwordConfirmation string, email string) (*User, error) {
	if password != passwordConfirmation {
		return nil, errors.New("password and confirmation don't match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := s.repo.CreateUser(tx, email, string(hash))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

// move builders elsewhere
func (s *GRPCServer) buildSession(tx *sql.Tx, newSessionUUID uuid.UUID, userID int, sessionToken string, furthestExpiration time.Time) (*Session, error) {
	session, err := s.repo.CreateSession(tx, newSessionUUID, userID, sessionToken, furthestExpiration)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return session, nil
}

func (s *GRPCServer) buildScopeGroupings(tx *sql.Tx, protoScopeGroupings []*proto.ScopeGrouping, sessionID int) ([]*ScopeGrouping, error) {
	scopeGroupings := make([]*ScopeGrouping, len(protoScopeGroupings))
	for i, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		scopeGrouping, err := s.repo.CreateScopeGrouping(tx, sessionID, sg.Scopes, exp)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		scopeGroupings[i] = scopeGrouping
	}

	return scopeGroupings, nil
}

func (s *GRPCServer) buildToken(user *User, sessionUUID uuid.UUID, protoScopeGroupings []*proto.ScopeGrouping) (tokenStr string, json string, furthestExpiration time.Time, err error) {
	tf := session_representations.NewTokenFactory(user.uuid, sessionUUID.String())
	for _, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			fmt.Println(err)
			return "", "", time.Now(), err
		}
		tf.AddScopeGrouping(sg.Scopes, exp)
	}

	sess, err := tf.GenerateSession()
	if err != nil {
		fmt.Println(err)
		return "", "", time.Now(), err
	}

	return  sess.Token, sess.Json, sess.FurthestExpiration, nil
}