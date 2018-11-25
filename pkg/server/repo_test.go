package server

import (
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


func TestGetUser(t *testing.T) {
	repo.GetUser()
}