package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
)
import "context"

type GRPCServer struct{
	repo *Repo
	dao *db.DAO
	builder *Builder
}

func NewGRPCServer(repo *Repo, dao *db.DAO) *GRPCServer {
	return &GRPCServer{repo, dao, &Builder{repo:repo, dao:dao}}
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	tx, err :=  s.dao.Conn.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := s.builder.buildUser(tx, request.Password, request.PasswordConfirmation, request.Email)

	sessionUUID := uuid.New()
	sessionToken, json, furthestExpiration, err := s.builder.buildToken(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return nil, err
	}


	session, err := s.builder.buildSession(tx, sessionUUID, user.id, sessionToken, furthestExpiration)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = s.builder.buildScopeGroupings(tx, request.ScopeGroupings, session.id)
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