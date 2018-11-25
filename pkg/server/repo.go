package server

import "github.com/willschroeder/fingerprint/pkg/db"

type Repo struct {
	dao *db.DAO
}

func (r *Repo) CreateUser() {
	sqlStatement := "INSERT INTO users (uuid, email) VALUES ($1, $2)"
	_, err := r.dao.Conn.Exec(sqlStatement, "f47ac10b-58cc-0372-8567-0e02b2c3d479", "test@test.com")
	if err != nil {
		panic(err)
	}
}

func (r *Repo) GetUser() {
	sqlStatement := "SELECT email FROM users WHERE uuid=$1"
	var user User

	row := r.dao.Conn.QueryRow(sqlStatement, "f47ac10b-58cc-0372-8567-0e02b2c3d479")
	err := row.Scan(&user.email)
	if err != nil {
		panic(err)
	}

	println(user.email)
}
