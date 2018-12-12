package server

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/services"
	"github.com/willschroeder/fingerprint/pkg/util"
	"log"
)
import "context"

type GRPCServer struct {
	db *sql.DB
}

func NewGRPCServer(db *sql.DB) *GRPCServer {
	return &GRPCServer{db: db}
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	usersService := services.NewUserService(s.db)

	user, session, err := usersService.CreateUser(request)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	return &proto.CreateUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) CreateGuestUser(_ context.Context, request *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	usersService := services.NewUserService(s.db)

	user, session, err := usersService.CreateGuestUser(request)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	return &proto.CreateGuestUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) GetUser(_ context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	repo := db.NewRepo(s.db)

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
	usersService := services.NewUserService(s.db)
	exp, err := util.ConvertTimestampToTime(request.Expiration)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	passwordReset, err := usersService.CreatePasswordResetToken(request.Email, exp)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	return &proto.CreatePasswordResetTokenResponse{PasswordResetToken: passwordReset.Token}, nil
}

func (s *GRPCServer) UpdateUserPassword(_ context.Context, request *proto.UpdateUserPasswordRequest) (*proto.UpdateUserPasswordResponse, error) {
	usersService := services.NewUserService(s.db)
	err := usersService.UpdateUserPassword(request.Email, request.PasswordResetToken, request.Password, request.PasswordConfirmation)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	// TODO Return actual status
	return &proto.UpdateUserPasswordResponse{Status: proto.UpdateUserPasswordResponse_SUCCESSFUL}, nil
}

func (s *GRPCServer) CreateSession(_ context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	usersService := services.NewUserService(s.db)
	user, err := usersService.ValidateEmailAndPassword(request.Email, request.Password)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	sessionService := services.NewSessionService(db.NewRepo(s.db))
	session, err := sessionService.CreateSession(user.Uuid, request.ScopeGroupings)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	// TODO Return actual status
	return &proto.CreateSessionResponse{Status: proto.CreateSessionResponse_SUCCESSFUL, Session: &proto.Session{Uuid: session.Uuid, Token: session.Token, Json: json}}, nil
}

func (s *GRPCServer) GetSession(_ context.Context, request *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	repo := db.NewRepo(s.db)
	session, err := repo.GetSessionWithToken(request.Token)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	json, err := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	return &proto.GetSessionResponse{Session: &proto.Session{Uuid: session.Uuid, Token: session.Token, Json: json}}, nil
}

func (s *GRPCServer) DeleteSession(_ context.Context, request *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	repo := db.NewRepo(s.db)

	switch representation := request.Representation.(type) {
	case *proto.DeleteSessionRequest_Uuid:
		err := repo.DeleteSessionWithUUID(representation.Uuid)
		if err != nil {
			return nil, logAndUnwrapError(err)
		}
		return &proto.DeleteSessionResponse{}, nil
	case *proto.DeleteSessionRequest_Token:
		err := repo.DeleteSessionWithToken(representation.Token)
		if err != nil {
			return nil, logAndUnwrapError(err)
		}
		return &proto.DeleteSessionResponse{}, nil
	}

	return &proto.DeleteSessionResponse{}, nil
}

func (s *GRPCServer) DeleteUser(_ context.Context, request *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	usersService := services.NewUserService(s.db)
	err := usersService.DeleteUser(request.Email, request.Password)
	if err != nil {
		return nil, logAndUnwrapError(err)
	}

	return &proto.DeleteUserResponse{}, nil
}

func logAndUnwrapError(err error) error {
	log.Printf("FATAL: %+v\n", err)
	return errors.Cause(err)
}
