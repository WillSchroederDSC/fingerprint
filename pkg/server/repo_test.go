package server

import (
	"github.com/brianvoe/gofakeit"
	"github.com/willschroeder/fingerprint/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

var testRepo *Repo
var testDAO *db.DAO
var testServer *GRPCServer

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	testDAO = db.ConnectToDatabase()
	defer testDAO.Conn.Close()
	testRepo = &Repo{dao: testDAO}
	testServer = NewGRPCServer(testRepo, testDAO)
	code := m.Run()
	os.Exit(code)
}

func createEncryptedPassword() string {
	pw := gofakeit.Password(true, true, true, true, true, 32)
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

func createTestUser(isGuest bool) *User {
	tx, _ := testDAO.Conn.Begin()
	user, _ := testRepo.CreateUser(tx, gofakeit.Email(), createEncryptedPassword(), isGuest)
	tx.Commit()
	return user
}

func TestRepoCreateUser(t *testing.T) {
	email := gofakeit.Email()
	tx, _ := testDAO.Conn.Begin()
	user, err := testRepo.CreateUser(tx, email, createEncryptedPassword(), false)
	if err != nil {
		t.Fatal(err)
	}
	tx.Commit()
	if user.email != email {
		t.Errorf("User not created with test email")
	}
}

func TestRepoGetUser(t *testing.T) {
	testUser := createTestUser(false)
	tx, _ := testDAO.Conn.Begin()
	gotUser, err := testRepo.GetUserWithUUIDUsingTx(tx, testUser.uuid)
	if err != nil {
		t.Fatal(err)
	}
	tx.Commit()
	if gotUser == nil || gotUser.email == "" {
		t.Errorf("Not able to get test user")
	}
}
