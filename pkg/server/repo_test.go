package server

import (
	"github.com/brianvoe/gofakeit"
	"github.com/willschroeder/fingerprint/pkg/db"
	"os"
	"testing"
)

var repo *Repo
var dao *db.DAO

func TestMain(m *testing.M) {
	dao = db.ConnectToDatabase()
	defer dao.Conn.Close()
	repo = &Repo{dao: dao}
	code := m.Run()
	os.Exit(code)
}

func createTestUser() *User {
	return repo.CreateUser(gofakeit.Email())
}

func TestCreateUser(t *testing.T) {
	email := gofakeit.Email()
	user := repo.CreateUser(email)
	if user.email != email {
		t.Errorf("User not created with test email")
	}
}

func TestGetUser(t *testing.T) {
	testUser := createTestUser()
	gotUser := repo.GetUser(testUser.uuid)
	if gotUser == nil || gotUser.email == "" {
		t.Errorf("Not able to get test user")
	}
}