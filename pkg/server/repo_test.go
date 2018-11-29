package server

import (
	"github.com/brianvoe/gofakeit"
	"github.com/willschroeder/fingerprint/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

var repo *Repo
var dao *db.DAO
var server *GRPCServer

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	dao = db.ConnectToDatabase()
	defer dao.Conn.Close()
	repo = &Repo{dao: dao}
	server = &GRPCServer{repo:repo, dao:dao}
	code := m.Run()
	os.Exit(code)
}

func createEncryptedPassword() string {
	pw := gofakeit.Password(true, true, true, true, true, 32)
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

func createTestUser() *User {
	user, _ := repo.CreateUser(gofakeit.Email(), createEncryptedPassword())
	return user
}

func TestRepoCreateUser(t *testing.T) {
	email := gofakeit.Email()
	user, err := repo.CreateUser(email, createEncryptedPassword())
	if err != nil {
		t.Fatal(err)
	}
	if user.email != email {
		t.Errorf("User not created with test email")
	}
}

func TestRepoGetUser(t *testing.T) {
	testUser := createTestUser()
	gotUser, err := repo.GetUserWithUUID(testUser.uuid)
	if err != nil {
		t.Fatal(err)
	}
	if gotUser == nil || gotUser.email == "" {
		t.Errorf("Not able to get test user")
	}
}