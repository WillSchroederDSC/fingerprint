package server

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"

)
import "context"

type GRPCServer struct {
	repo    *db.Repo
	dao     *db.DAO
	actions *Actions
}

func NewGRPCServer(dao *db.DAO) *GRPCServer {
	return &GRPCServer{db.NewRepo(dao), dao, NewActions(dao)}
}

func (s *GRPCServer) CreateUser(_ context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	tx, err := s.dao.Conn.Begin()
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		return nil, errors.Cause(err)
	}

	user, err := s.actions.buildUser(tx, request.Email, request.Password, request.PasswordConfirmation)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sessionUUID := uuid.New().String()
	sessionToken, json, err := s.actions.buildSessionRepresentation(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	session, err := s.actions.buildSession(tx, sessionUUID, user.Uuid, sessionToken)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = s.actions.buildScopeGroupings(tx, request.ScopeGroupings, session.Uuid)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return &proto.CreateUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) GetUser(_ context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	switch ident := request.Identifier.(type) {
	case *proto.GetUserRequest_Email:
		user, err := s.repo.GetUserWithEmail(ident.Email)
		if err != nil {
			panic(err)
		}
		return &proto.GetUserResponse{User: user.ConvertToProtobuff()}, nil
	case *proto.GetUserRequest_Uuid:
		user, err := s.repo.GetUserWithUUID(ident.Uuid)
		if err != nil {
			panic(err)
		}
		return &proto.GetUserResponse{User: user.ConvertToProtobuff()}, nil
	}

	return nil, errors.New("unknown user identifier")
}

func (s *GRPCServer) CreateGuestUser(_ context.Context, request *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	tx, err := s.dao.Conn.Begin()
	if err != nil {
		panic(err)
	}

	user, err := s.actions.buildGuestUser(tx, request.Email)

	sessionUUID := uuid.New().String()
	sessionToken, json, err := s.actions.buildSessionRepresentation(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	session, err := s.actions.buildSession(tx, sessionUUID, user.Uuid, sessionToken)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = s.actions.buildScopeGroupings(tx, request.ScopeGroupings, session.Uuid)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return &proto.CreateGuestUserResponse{User: user.ConvertToProtobuff(), Session: session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) CreatePasswordResetToken(_ context.Context, request *proto.CreatePasswordResetTokenRequest) (*proto.CreatePasswordResetTokenResponse, error) {
	user, err := s.repo.GetUserWithEmail(request.Email)
	if err != nil {
		panic(err)
	}

	err = s.repo.DeleteAllPasswordResetTokensForUser(user.Uuid)
	if err != nil {
		panic(err)
	}

	exp, err := ptypes.Timestamp(request.Expiration)
	if err != nil {
		panic(err)
	}

	resetToken, err := s.repo.CreatePasswordResetToken(user.Uuid, exp)
	if err != nil {
		panic(err)
	}

	return &proto.CreatePasswordResetTokenResponse{PasswordResetToken: resetToken.Token}, nil
}

func (s *GRPCServer) UpdateUserPassword(_ context.Context, request *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	err := s.actions.updateUserPassword(request.Email, request.PasswordResetToken, request.Password, request.PasswordConfirmation)
	if err != nil {
		panic(err)
	}

	return &proto.ResetUserPasswordResponse{Status: proto.ResetUserPasswordResponse_SUCCESSFUL}, nil
}

func (s *GRPCServer) CreateSession(_ context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	user, err := s.repo.GetUserWithEmail(request.Email)
	if err != nil {
		panic(err)
	}

	hash, err := BuildPasswordHash(request.Password)
	if err != nil {
		panic(err)
	}

	if hash != user.EncryptedPassword {
		return nil, errors.New("incorrect password")
	}

	tx, err := s.dao.Conn.Begin()
	if err != nil {
		panic(err)
	}

	sessionUUID := uuid.New().String()
	sessionToken, json, err := s.actions.buildSessionRepresentation(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	session, err := s.actions.buildSession(tx, sessionUUID, user.Uuid, sessionToken)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = s.actions.buildScopeGroupings(tx, request.ScopeGroupings, session.Uuid)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return &proto.CreateSessionResponse{Session: &proto.Session{Uuid: session.Uuid, Token: sessionToken, Json: json}}, nil
}

func (s *GRPCServer) GetSession(_ context.Context, request *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	session, err := s.repo.GetSessionWithToken(request.Token)
	if err != nil {
		panic(err)
	}

	json, err := session_representations.DecodeTokenToJson(session.Token)
	if err != nil {
		panic(err)
	}


	return &proto.GetSessionResponse{Session: &proto.Session{Uuid: session.Uuid, Token: session.Token, Json: json}}, nil
}

func (s *GRPCServer) DeleteSession(_ context.Context, request *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	switch representation := request.Representation.(type) {
	case *proto.DeleteSessionRequest_Uuid:
		err := s.repo.DeleteSessionWithUUID(representation.Uuid)
		if err != nil {
			panic(err)
		}
		return &proto.DeleteSessionResponse{Successful: true}, nil
	case *proto.DeleteSessionRequest_Token:
		err := s.repo.DeleteSessionWithToken(representation.Token)
		if err != nil {
			panic(err)
		}
		return &proto.DeleteSessionResponse{Successful: true}, nil
	}

	return &proto.DeleteSessionResponse{Successful: false}, nil
}
