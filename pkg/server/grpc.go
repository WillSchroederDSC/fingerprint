package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/db/services"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
	"github.com/willschroeder/fingerprint/pkg/util"
)
import "context"

type GRPCServer struct {
	dao     *db.DAO
}

func NewGRPCServer(dao *db.DAO) *GRPCServer {
	return &GRPCServer{dao: dao}
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	usersService := services.NewUserService(s.dao)

	user, session, err := usersService.CreateUser(request)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	json, err  := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	return &proto.CreateUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) CreateGuestUser(_ context.Context, request *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	usersService := services.NewUserService(s.dao)

	user, session, err := usersService.CreateGuestUser(request)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	json, err  := services.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
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
		return nil, PrintAndUnwrapError(err)
	}

	passwordReset, err := usersService.CreatePasswordResetToken(request.Email, exp)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	return &proto.CreatePasswordResetTokenResponse{PasswordResetToken: passwordReset.Token}, nil
}

func (s *GRPCServer) UpdateUserPassword(_ context.Context, request *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	usersService := services.NewUserService(s.dao)
	err := usersService.UpdateUserPassword(request.Email, request.PasswordResetToken, request.Password, request.PasswordConfirmation)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	// TODO Return actual status
	return &proto.ResetUserPasswordResponse{Status: proto.ResetUserPasswordResponse_SUCCESSFUL}, nil
}

func (s *GRPCServer) CreateSession(_ context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	//repo := db.NewRepo(s.dao.DB)
	//
	//user, err := repo.GetUserWithEmail(request.Email)
	//if err != nil {
	//	return nil, PrintAndUnwrapError(err)
	//}
	//
	//hash, err := services.BuildPasswordHash(request.Password)
	//if err != nil {
	//	return nil, PrintAndUnwrapError(err)
	//}
	//
	//// TODO this check should be moved to the user service
	//if hash != user.EncryptedPassword {
	//	return nil, errors.New("incorrect password")
	//}
	//
	//sessionService := services.NewSessionService(db.NewRepo(s.dao.DB))
	//sessionService.CreateSession(user.Uuid)
	//
	//tx, err := s.dao.DB.Begin()
	//if err != nil {
	//	panic(err)
	//}
	//
	//sessionUUID := uuid.New().String()
	//sessionToken, json, err := s.actions.buildSessionRepresentation(user, sessionUUID, request.ScopeGroupings)
	//if err != nil {
	//	tx.Rollback()
	//	panic(err)
	//}
	//
	//session, err := s.actions.buildSession(tx, sessionUUID, user.Uuid, sessionToken)
	//if err != nil {
	//	tx.Rollback()
	//	panic(err)
	//}
	//
	//_, err = s.actions.buildScopeGroupings(tx, request.ScopeGroupings, session.Uuid)
	//if err != nil {
	//	tx.Rollback()
	//	panic(err)
	//}
	//
	//err = tx.Commit()
	//if err != nil {
	//	panic(err)
	//}
	//return &proto.CreateSessionResponse{Session: &proto.Session{Uuid: session.Uuid, Token: sessionToken, Json: json}}, nil
	panic("not yet")
}

func (s *GRPCServer) GetSession(_ context.Context, request *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	repo := db.NewRepo(s.dao.DB)
	session, err := repo.GetSessionWithToken(request.Token)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	json, err := session_representations.DecodeTokenToJson(session.Token)
	if err != nil {
		return nil, PrintAndUnwrapError(err)
	}

	return &proto.GetSessionResponse{Session: &proto.Session{Uuid: session.Uuid, Token: session.Token, Json: json}}, nil
}

func (s *GRPCServer) DeleteSession(_ context.Context, request *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	repo := db.NewRepo(s.dao.DB)

	switch representation := request.Representation.(type) {
	case *proto.DeleteSessionRequest_Uuid:
		err := repo.DeleteSessionWithUUID(representation.Uuid)
		if err != nil {
			return nil, PrintAndUnwrapError(err)
		}
		return &proto.DeleteSessionResponse{Successful: true}, nil
	case *proto.DeleteSessionRequest_Token:
		err := repo.DeleteSessionWithToken(representation.Token)
		if err != nil {
			return nil, PrintAndUnwrapError(err)
		}
		return &proto.DeleteSessionResponse{Successful: true}, nil
	}

	return &proto.DeleteSessionResponse{Successful: false}, nil
}

func PrintAndUnwrapError(err error) error {
	fmt.Printf("FATAL: %+v\n", err)
	return errors.Cause(err)
}