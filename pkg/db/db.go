package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "fingerprint_development"
)

type DB struct {
	Connection *sql.DB
}

func ConnectToDatabase() (*DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// Ensures Connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &DB{Connection: db}
}

func (db *DB) CreateUser() {
	sqlStatement := "INSERT INTO users (uuid, email) VALUES ('1112', 'test@test.io')"
	_, err := db.Connection.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}