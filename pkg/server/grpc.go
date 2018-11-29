package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"golang.org/x/crypto/bcrypt"
)
import "context"

type GRPCServer struct{
	repo *Repo
	dao *db.DAO
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if request.Password != request.PasswordConfirmation {
		return nil, errors.New("password and confirmation don't match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := s.repo.CreateUser(request.Email, string(hash))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	session, err := s.repo.CreateSession(user.id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create Scopes
	scopeGroupings := make([]*ScopeGrouping, len(request.ScopeGroupings))
	for i, sg := range request.ScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		scopeGrouping, err := s.repo.CreateScopeGrouping(session.id, sg.Scopes, exp)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		scopeGroupings[i] = scopeGrouping
	}

	return &proto.CreateUserResponse{User:user.ConvertToProtobuff(), Session:session.ConvertToProtobuff("111", "{}")}, nil
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