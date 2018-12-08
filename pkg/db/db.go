package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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
	DB *sql.DB
}

func (dao *DAO) NewTransaction() (*sql.Tx, error) {
	tx, err := dao.DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't build transaction")
	}

	return tx, nil
}

func ConnectToDatabase() *DAO {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Ensures DB
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &DAO{DB: conn}
}

func MigrateUp() {

}

func MigrateDown() {

}

func HandleRollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Println("Failed to rollback Transaction\n", err)
	}
}
