package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "fingerprint_development"
)

type DAO struct {
	Conn *sql.DB
}

func ConnectToDatabase() *DAO {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Ensures Conn
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &DAO{Conn: conn}
}

func MigrateUp() {

}

func MigrateDown() {

}
