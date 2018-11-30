package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)
import "context"

type GRPCServer struct{
	repo *Repo
	dao *db.DAO
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user, err := s.buildUser(request.Password, request.PasswordConfirmation, request.Email)

	sessionUUID := uuid.New()
	sessionToken, json, furthestExpiration, err := s.buildToken(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	session, err := s.buildSession(sessionUUID, user.id, sessionToken, furthestExpiration)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = s.buildScopeGroupings(request.ScopeGroupings, session.id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.CreateUserResponse{User:user.ConvertToProtobuff(), Session:session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) GetUser(context.Context, *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user := &proto.User{Uuid: "1111", Email:"fake@email.com"}
	return &proto.GetUserResponse{User:user}, nil
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

func (s *GRPCServer) buildUser(password string, passwordConfirmation string, email string) (*User, error) {
	if password != passwordConfirmation {
		return nil, errors.New("password and confirmation don't match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := s.repo.CreateUser(email, string(hash))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (s *GRPCServer) buildSession(newSessionUUID uuid.UUID, userID int, sessionToken string, furthestExpiration time.Time) (*Session, error) {
	session, err := s.repo.CreateSession(newSessionUUID, userID, sessionToken, furthestExpiration)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return session, nil
}

func (s *GRPCServer) buildScopeGroupings(protoScopeGroupings []*proto.ScopeGrouping, sessionID int) ([]*ScopeGrouping, error) {
	scopeGroupings := make([]*ScopeGrouping, len(protoScopeGroupings))
	for i, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		scopeGrouping, err := s.repo.CreateScopeGrouping(sessionID, sg.Scopes, exp)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		scopeGroupings[i] = scopeGrouping
	}

	return scopeGroupings, nil
}

func (s *GRPCServer) buildToken(user *User, sessionUUID uuid.UUID, protoScopeGroupings []*proto.ScopeGrouping) (tokenStr string, json string, furthestExpiration time.Time, err error) {
	tf := token.NewTokenFactory(user.uuid, sessionUUID.String())
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