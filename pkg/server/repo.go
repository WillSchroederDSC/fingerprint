package server

import "github.com/willschroeder/fingerprint/pkg/db"

func CreateUser(db *db.DAO) {
	sqlStatement := "INSERT INTO users (uuid, email) VALUES ($1, $2)"
	_, err := db.Conn.Exec(sqlStatement, "f47ac10b-58cc-0372-8567-0e02b2c3d479", "test@test.com")
	if err != nil {
		panic(err)
	}
}

func GetUser(db *db.DAO) {
	sqlStatement := `SELECT email FROM users WHERE uuid=$1;`
	var email string

	row := db.Conn.QueryRow(sqlStatement, "f47ac10b-58cc-0372-8567-0e02b2c3d479")
	err := row.Scan(&email)
	if err != nil {
		panic(err)
	}

	println(email)
}
