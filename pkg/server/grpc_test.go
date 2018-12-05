package server

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	oneHour, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * time.Duration(1)))
	twoHour, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * time.Duration(2)))

	req := &proto.CreateUserRequest{
		Email:                gofakeit.Email(),
		Password:             "test",
		PasswordConfirmation: "test",
		ScopeGroupings: []*proto.ScopeGrouping{
			{
				Scopes:     []string{"read"},
				Expiration: oneHour,
			},
			{
				Scopes:     []string{"write"},
				Expiration: twoHour,
			},
		},
	}
	res, _ := testServer.CreateUser(context.Background(), req)
	print(res.Session.Token)
}

func TestGetUserWithUUID(t *testing.T) {
	user := createTestUser(false)
	req := &proto.GetUserRequest{
		Identifier: &proto.GetUserRequest_Uuid{
			Uuid: user.uuid,
		},
	}

	res, _ := testServer.GetUser(context.Background(), req)
	print(res.User.Email)
}
