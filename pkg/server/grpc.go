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

	user, err := s.builder.buildUser(tx, request.Email, request.Password, request.PasswordConfirmation)
	if err != nil {
		tx.Rollback()
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

func (s *GRPCServer) CreateGuestUser(_ context.Context, request *proto.CreateGuestUserRequest) (*proto.CreateGuestUserResponse, error) {
	tx, err :=  s.dao.Conn.Begin()
	if err != nil {
		panic(err)
	}

	user, err := s.builder.buildGuestUser(tx,request.Email)

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

	return &proto.CreateGuestUserResponse{User:user.ConvertToProtobuff(), Session:session.ConvertToProtobuff(json)}, nil
}

func (s *GRPCServer) CreatePasswordResetToken(_ context.Context, request *proto.CreatePasswordResetTokenRequest) (*proto.CreatePasswordResetTokenResponse, error) {
	//token, err := s.repo.UpdateUserPasswordResetToken(request.Email)
	//if err != nil {
	//	panic(err)
	//}
	//
	//return &proto.CreatePasswordResetTokenResponse{PasswordResetToken:token}, nil
	panic("soon")
}

func (s *GRPCServer) UpdateUserPassword(_ context.Context, request *proto.ResetUserPasswordRequest) (*proto.ResetUserPasswordResponse, error) {
	err := s.builder.updateUserPassword(request.Email,request.PasswordResetToken,request.Password,request.PasswordConfirmation)
	if err != nil {
		panic(err)
	}

	return &proto.ResetUserPasswordResponse{Status:proto.ResetUserPasswordResponse_SUCCESSFUL}, nil
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
	return &proto.CreateSessionResponse{Session: &proto.Session{Uuid:session.uuid, Token:sessionToken, Json:json}}, nil
}

func (s *GRPCServer) GetSession(_ context.Context, request *proto.GetSessionRequest) (*proto.GetSessionResponse, error) {
	session, err := s.repo.GetSessionWithToken(request.Token)
	if err != nil {
		panic(err)
	}

	json := session_representations.DecodeTokenToJson(session.token)

	return &proto.GetSessionResponse{Session:&proto.Session{Uuid:session.uuid, Token:session.token, Json:json}}, nil
}

func (s *GRPCServer) DeleteSession(_ context.Context, request *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	//successful, err := s.repo.DeleteSessionWithUUID(request.Uuid)
	//if err != nil {
	//	panic(err)
	//}
	//
	return &proto.DeleteSessionResponse{Successful:false}, nil
}