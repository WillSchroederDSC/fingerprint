package server

import (
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"golang.org/x/crypto/bcrypt"
	"log"
	"errors"
	"time"
)
import "context"

type GRPCSerer struct{
	repo *Repo
	dao *db.DAO
	logger *log.Logger
}

func (s *GRPCSerer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if request.Password != request.PasswordConfirmation {
		return nil, errors.New("password and confirmation don't match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(request.Email, string(hash))
	if err != nil {
		return nil, err
	}

	session, err := s.repo.CreateSession(user.id, time.Now().UTC())
	if err != nil {
		return nil, err
	}
		
	return &proto.CreateUserResponse{User:user.ConvertToProtobuff(), Session:session.ConvertToProtobuff("111", "{}")}, nil
}

func (s *GRPCSerer) GetUser(context.Context, *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user := &proto.User{Uuid: "1111", Email:"fake@email.com"}
	return &proto.GetUserResponse{User:user}, nil
}

func (s *GRPCSerer) CreateGuestUser(context.Context, *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	panic("implement me")
}

func (s *GRPCSerer) CreatePasswordResetToken(context.Context, *proto.CreatePasswordResetTokenResponse) (*proto.CreatePasswordResetTokenResponse, error) {
	panic("implement me")
}

func (s *GRPCSerer) UpdateUserPassword(context.Context, *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	panic("implement me")
}

func (s *GRPCSerer) CreateSession(context.Context, *proto.CreatePasswordResetTokenRequest) (*proto.CreatePasswordResetTokenResponse, error) {
	panic("implement me")
}

func (s *GRPCSerer) CreateSessionRevoke(context.Context, *proto.CreateSessionRevokeRequest) (*proto.CreateSessionRevokeResponse, error) {
	panic("implement me")
}

func (s *GRPCSerer) GetSession(context.Context, *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	panic("implement me")
}