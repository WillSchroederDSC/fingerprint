package server

import (
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
)

type Repo struct {
	dao *db.DAO
}

func (r *Repo) CreateUser(email string) *User {
	customerUUID := uuid.New().String()
	sqlStatement := "INSERT INTO users (uuid, email) VALUES ($1, $2)"
	_, err := r.dao.Conn.Exec(sqlStatement, customerUUID, email)
	if err != nil {
		panic(err)
	}

	return r.GetUser(customerUUID)
}

func (r *Repo) GetUser(customerUUID string) *User {
	sqlStatement := "SELECT uuid,email FROM users WHERE uuid=$1"
	var user User

	row := r.dao.Conn.QueryRow(sqlStatement, customerUUID)
	err := row.Scan(&user.uuid, &user.email)
	if err != nil {
		panic(err)
	}

	return &user
}
