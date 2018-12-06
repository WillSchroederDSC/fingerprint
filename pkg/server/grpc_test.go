package server

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
	"time"
)

var testRepo *db.Repo
var testDAO *db.DAO
var testServer *GRPCServer

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	testDAO = db.ConnectToDatabase()
	defer testDAO.Conn.Close()
	testRepo = &db.Repo{Dao: testDAO}
	testServer = NewGRPCServer(testRepo, testDAO)
	code := m.Run()
	os.Exit(code)
}

func createEncryptedPassword() string {
	pw := gofakeit.Password(true, true, true, true, true, 32)
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

func createTestUser(isGuest bool) *db.User {
	tx, _ := testDAO.Conn.Begin()
	user, _ := testRepo.CreateUser(tx, gofakeit.Email(), createEncryptedPassword(), isGuest)
	tx.Commit()
	return user
}

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
			Uuid: user.Uuid,
		},
	}

	res, _ := testServer.GetUser(context.Background(), req)
	print(res.User.Email)
}
