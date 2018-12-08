package server

import (
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/services"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/util"
	"log"
)
import "context"

type GRPCServer struct {
	dao *db.DAO
}

func NewGRPCServer(dao *db.DAO) *GRPCServer {
	return &GRPCServer{dao: dao}
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	usersService := services.NewUserService(s.dao)

	user, session, err := usersService.CreateUser(request)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	return &proto.CreateUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) CreateGuestUser(_ context.Context, request *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	usersService := services.NewUserService(s.dao)

	user, session, err := usersService.CreateGuestUser(request)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	return &proto.CreateGuestUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) GetUser(_ context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	repo := db.NewRepo(s.dao.DB)

	switch ident := request.Identifier.(type) {
	case *proto.GetUserRequest_Email:
		user, err := repo.GetUserWithEmail(ident.Email)
		if err != nil {
			panic(err)
		}
		return &proto.GetUserResponse{User: user.ConvertToProtobuff()}, nil
	case *proto.GetUserRequest_Uuid:
		user, err := repo.GetUserWithUUID(ident.Uuid)
		if err != nil {
			panic(err)
		}
		return &proto.GetUserResponse{User: user.ConvertToProtobuff()}, nil
	}

	return nil, errors.New("unknown user identifier type")
}

func (s *GRPCServer) CreatePasswordResetToken(_ context.Context, request *proto.CreatePasswordResetTokenRequest) (*proto.CreatePasswordResetTokenResponse, error) {
	usersService := services.NewUserService(s.dao)
	exp, err := util.ConvertTimestampToTime(request.Expiration)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	passwordReset, err := usersService.CreatePasswordResetToken(request.Email, exp)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	return &proto.CreatePasswordResetTokenResponse{PasswordResetToken: passwordReset.Token}, nil
}

func (s *GRPCServer) UpdateUserPassword(_ context.Context, request *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	usersService := services.NewUserService(s.dao)
	err := usersService.UpdateUserPassword(request.Email, request.PasswordResetToken, request.Password, request.PasswordConfirmation)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	// TODO Return actual status
	return &proto.ResetUserPasswordResponse{Status: proto.ResetUserPasswordResponse_SUCCESSFUL}, nil
}

func (s *GRPCServer) CreateSession(_ context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	usersService := services.NewUserService(s.dao)
	user, err := usersService.ValidateEmailAndPassword(request.Email, request.Password)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	sessionService := services.NewSessionService(db.NewRepo(s.dao.DB))
	session, err := sessionService.CreateSession(user.Uuid, request.ScopeGroupings)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	// TODO Return actual status
	return &proto.CreateSessionResponse{Status: proto.CreateSessionResponse_SUCCESSFUL, Session: &proto.Session{Uuid: session.Uuid, Token: session.Token, Json: json}}, nil
}

func (s *GRPCServer) GetSession(_ context.Context, request *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	repo := db.NewRepo(s.dao.DB)
	session, err := repo.GetSessionWithToken(request.Token)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	return &proto.GetSessionResponse{Session: &proto.Session{Uuid: session.Uuid, Token: session.Token, Json: json}}, nil
}

func (s *GRPCServer) DeleteSession(_ context.Context, request *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	repo := db.NewRepo(s.dao.DB)

	switch representation := request.Representation.(type) {
	case *proto.DeleteSessionRequest_Uuid:
		err := repo.DeleteSessionWithUUID(representation.Uuid)
		if err != nil {
			return nil, LogAndUnwrapError(err)
		}
		return &proto.DeleteSessionResponse{}, nil
	case *proto.DeleteSessionRequest_Token:
		err := repo.DeleteSessionWithToken(representation.Token)
		if err != nil {
			return nil, LogAndUnwrapError(err)
		}
		return &proto.DeleteSessionResponse{}, nil
	}

	return &proto.DeleteSessionResponse{}, nil
}

func (s *GRPCServer) DeleteUser(_ context.Context, request *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	usersService := services.NewUserService(s.dao)
	err := usersService.DeleteUser(request.Email, request.Password)
	if err != nil {
		return nil, LogAndUnwrapError(err)
	}

	return &proto.DeleteUserResponse{}, nil
}

func LogAndUnwrapError(err error) error {
	log.Printf("FATAL: %+v\n", err)
	return errors.Cause(err)
}
