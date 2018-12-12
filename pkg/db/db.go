package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	"log"
	"path/filepath"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "fingerprint_development"
)


func NewTransaction(DB *sql.DB) (*sql.Tx, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't build transaction")
	}

	return tx, nil
}

func ConnectToDatabase() *sql.DB {
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
	return conn
}

func MigrateUp() {
	absPath, _ := filepath.Abs("./migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: absPath,
	}

	conn := ConnectToDatabase()
	HandleClose(conn)

	n, err := migrate.Exec(conn, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Applied %d migrations\n", n)
}

func MigrateDown() {
	absPath, _ := filepath.Abs("./migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: absPath,
	}

	conn := ConnectToDatabase()
	HandleClose(conn)

	n, err := migrate.Exec(conn, "postgres", migrations, migrate.Down)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rolled back %d migrations\n", n)
}

func HandleClose(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Failed to close db connection\n", err)
	}
}

func HandleRollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Println("Failed to rollback Transaction\n", err)
	}
}
