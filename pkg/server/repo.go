package server

import (
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
	"time"
)

type Repo struct {
	dao *db.DAO
}

func (r *Repo) CreateUser(email string, encryptedPassword string) (*User, error) {
	customerUUID := uuid.New().String()

	sqlStatement := "INSERT INTO users (uuid, email, encrypted_password, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.dao.Conn.Exec(sqlStatement, customerUUID, email, encryptedPassword, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	user, err := r.GetUserWithUUID(customerUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repo) GetUserWithUUID(customerUUID string) (*User, error) {
	sqlStatement := "SELECT id,uuid,email FROM users WHERE uuid=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, customerUUID)
	var user User
	err := row.Scan(&user.id, &user.uuid, &user.email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) CreateSession(customerId int, expiration time.Time) (*Session, error) {
	sessionUUID := uuid.New().String()

	sqlStatement := "INSERT INTO sessions (uuid, expiration) VALUES ($1, $2)"
	_, err := r.dao.Conn.Exec(sqlStatement, sessionUUID, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	session, err := r.GetSessionWithUUID(sessionUUID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repo) GetSessionWithUUID(sessionUUID string) (*Session, error) {
	sqlStatement := "SELECT id,uuid,expiration FROM sessions WHERE uuid=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, sessionUUID)
	var session Session
	err := row.Scan(&session.id, &session.uuid, &session.expiration)
	if err != nil {
		return nil, err
	}

	return &session, nil
}