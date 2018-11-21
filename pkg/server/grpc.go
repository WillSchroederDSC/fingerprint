package server

import "github.com/willschroeder/fingerprint/pkg/proto"
import "context"

type Server struct{}

func (s *Server) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user := &proto.User{Uuid: "1111", Email:request.Email}
	session := &proto.Session{Uuid:"1111", Token:"token", Json:"{}"}
	return &proto.CreateUserResponse{User:user, Session:session}, nil
}

func (s *Server) GetUser(context.Context, *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user := &proto.User{Uuid: "1111", Email:"fake@email.com"}
	return &proto.GetUserResponse{User:user}, nil
}

func (s *Server) CreateGuestUser(context.Context, *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	panic("implement me")
}

func (s *Server) CreatePasswordResetToken(context.Context, *proto.CreatePasswordResetTokenResponse) (*proto.CreatePasswordResetTokenResponse, error) {
	panic("implement me")
}

func (s *Server) UpdateUserPassword(context.Context, *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	panic("implement me")
}

func (s *Server) CreateSession(context.Context, *proto.CreatePasswordResetTokenRequest) (*proto.CreatePasswordResetTokenResponse, error) {
	panic("implement me")
}

func (s *Server) CreateSessionRevoke(context.Context, *proto.CreateSessionRevokeRequest) (*proto.CreateSessionRevokeResponse, error) {
	panic("implement me")
}

func (s *Server) GetSession(context.Context, *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	panic("implement me")
}