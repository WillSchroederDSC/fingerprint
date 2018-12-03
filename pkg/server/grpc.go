package server

import (
	"errors"
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
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
		panic(err)
	}

	user, err := s.builder.buildUser(tx, request.Password, request.PasswordConfirmation, request.Email)

	sessionUUID := uuid.New()
	sessionToken, json, furthestExpiration, err := s.builder.buildToken(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		panic(err)
	}


	session, err := s.builder.buildSession(tx, sessionUUID, user.id, sessionToken, furthestExpiration)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = s.builder.buildScopeGroupings(tx, request.ScopeGroupings, session.id)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return &proto.CreateUserResponse{User:user.ConvertToProtobuff(), Session:session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) GetUser(_ context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	switch ident := request.Identifier.(type) {
	case *proto.GetUserRequest_Email:
		user, err := s.repo.GetUserWithEmail(ident.Email)
		if err != nil {
			panic(err)
		}
		return &proto.GetUserResponse{User:user.ConvertToProtobuff()}, nil
	case *proto.GetUserRequest_Uuid:
		user, err := s.repo.GetUserWithUUID(ident.Uuid)
		if err != nil {
			panic(err)
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

func (s *GRPCServer) CreateSession(_ context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	user, err := s.repo.GetUserWithEmail(request.Email)
	if err != nil {
		panic(err)
	}

	hash, err := BuildPasswordHash(request.Password)
	if err != nil {
		panic(err)
	}

	if hash != user.encryptedPassword {
		return nil, errors.New("incorrect password")
	}

	tx, err :=  s.dao.Conn.Begin()
	if err != nil {
		panic(err)
	}

	sessionUUID := uuid.New()
	sessionToken, json, furthestExpiration, err := s.builder.buildToken(user, sessionUUID, request.ScopeGroupings)
	if err != nil {
		tx.Rollback()
		panic(err)
	}


	session, err := s.builder.buildSession(tx, sessionUUID, user.id, sessionToken, furthestExpiration)
	if err != nil {
		panic(err)
	}

	_, err = s.builder.buildScopeGroupings(tx, request.ScopeGroupings, session.id)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return &proto.CreateSessionResponse{Session: &proto.Session{Uuid:session.uuid, Token:sessionToken, Json:json}}, nil
}

func (s *GRPCServer) CreateSessionRevoke(context.Context, *proto.CreateSessionRevokeRequest) (*proto.CreateSessionRevokeResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) GetSession(_ context.Context, request *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	session, err := s.repo.GetSessionWithToken(request.Token)
	if err != nil {
		panic(err)
	}

	json := session_representations.DecodeTokenToJson(session.token)

	return &proto.GetSessionResponse{Session:&proto.Session{Uuid:session.uuid, Token:session.token, Json:json}}, nil
}