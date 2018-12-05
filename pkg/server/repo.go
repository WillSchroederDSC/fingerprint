package server

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/willschroeder/fingerprint/pkg/db"
	"time"
)

type Repo struct {
	dao *db.DAO
}

func (r *Repo) CreateUser(tx *sql.Tx, email string, encryptedPassword string, isGuest bool) (*User, error) {
	userUUID := uuid.New().String()

	sqlStatement := "INSERT INTO users (uuid, email, encrypted_password, is_guest, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(sqlStatement, userUUID, email, encryptedPassword, isGuest, time.Now().UTC())
	if err != nil {
		panic(err)
	}

	user, err := r.GetUserWithUUIDUsingTx(tx, userUUID)
	if err != nil {
		panic(err)
	}

	return user, nil
}

func (r *Repo) GetUserWithUUID(userUUID string) (*User, error) {
	sqlStatement := "SELECT id,uuid,email,encrypted_password,is_guest FROM users WHERE uuid=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, userUUID)
	var user User
	err := row.Scan(&user.id, &user.uuid, &user.email, &user.encryptedPassword, &user.isGuest)
	if err != nil {
		panic(err)
	}

	return &user, nil
}

func (r *Repo) UpdateUserPassword(email string, encryptedPassword string) error {
	sqlStatement := "UPDATE users SET encrypted_password=$1WHERE email=$3"
	_, err := r.dao.Conn.Exec(sqlStatement, email)
	if err != nil {
		panic(err)
	}

	return nil
}

func (r *Repo) GetUserWithEmail(email string) (*User, error) {
	sqlStatement := "SELECT id,uuid,email,encrypted_password,is_guest FROM users WHERE email=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, email)
	var user User
	err := row.Scan(&user.id, &user.uuid, &user.email, &user.encryptedPassword, &user.isGuest)
	if err != nil {
		panic(err)
	}

	return &user, nil
}

func (r *Repo) GetUserWithUUIDUsingTx(tx *sql.Tx, userUUID string) (*User, error) {
	sqlStatement := "SELECT id,uuid,email,encrypted_password,is_guest FROM users WHERE uuid=$1"

	row := tx.QueryRow(sqlStatement, userUUID)
	var user User
	err := row.Scan(&user.id, &user.uuid, &user.email, &user.encryptedPassword, &user.isGuest)
	if err != nil {
		panic(err)
	}

	return &user, nil
}

func (r *Repo) CreateSession(tx *sql.Tx, newSessionUUID uuid.UUID, userId int, token string, expiration time.Time) (*Session, error) {
	sessionUUID := newSessionUUID.String()

	sqlStatement := "INSERT INTO sessions (uuid, user_id, token, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(sqlStatement, sessionUUID, userId, token, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		panic(err)
	}

	session, err := r.GetSessionWithUUIDUsingTx(tx, sessionUUID)
	if err != nil {
		panic(err)
	}

	return session, nil
}

func (r *Repo) GetSessionWithUUIDUsingTx(tx *sql.Tx, sessionUUID string) (*Session, error) {
	sqlStatement := "SELECT id,uuid,token,expiration FROM sessions WHERE uuid=$1"

	row := tx.QueryRow(sqlStatement, sessionUUID)
	var session Session
	err := row.Scan(&session.id, &session.uuid, &session.token, &session.expiration)
	if err != nil {
		panic(err)
	}

	return &session, nil
}

func (r *Repo) GetSessionWithToken(token string) (*Session, error) {
	sqlStatement := "SELECT id,uuid,token,expiration FROM sessions WHERE token=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, token)
	var session Session
	err := row.Scan(&session.id, &session.uuid, &session.token, &session.expiration)
	if err != nil {
		panic(err)
	}

	return &session, nil
}

func (r *Repo) CreateScopeGrouping(tx *sql.Tx, sessionId int, scopes []string, expiration time.Time) (*ScopeGrouping, error) {
	groupingUUID := uuid.New().String()

	sqlStatement := "INSERT INTO scope_groupings (uuid, session_id, scopes, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(sqlStatement, groupingUUID, sessionId, pq.Array(scopes), expiration, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	grouping, err := r.GetScopeGroupingWithUUID(tx, groupingUUID)
	if err != nil {
		panic(err)
	}

	return grouping, nil
}

func (r *Repo) GetScopeGroupingWithUUID(tx *sql.Tx, groupingUUID string) (*ScopeGrouping, error) {
	sqlStatement := "SELECT id,uuid,scopes,expiration FROM scope_groupings WHERE uuid=$1"
	row := tx.QueryRow(sqlStatement, groupingUUID)
	var sg ScopeGrouping
	err := row.Scan(&sg.id, &sg.uuid, pq.Array(&sg.scopes), &sg.expiration)
	if err != nil {
		panic(err)
	}

	return &sg, nil
}

func (r *Repo) DeleteSessionWithUUID(sessionUUID string) error {
	sqlStatement := "DELETE FROM sessions WHERE uuid=$1"
	_, err := r.dao.Conn.Exec(sqlStatement, sessionUUID)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Repo) CreatePasswordResetToken(userId int, expiration time.Time) (*PasswordResetToken, error) {
	resetUUID := uuid.New().String()
	newResetToken := db.RandomString(16)
	sqlStatement := "INSERT INTO password_reset_tokens (uuid, user_id, token, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.dao.Conn.Exec(sqlStatement, resetUUID, userId, newResetToken, expiration, time.Now().UTC())
	if err != nil {
		panic(err)
	}

	prt, err := r.GetPasswordResetToken(resetUUID)
	if err != nil {
		panic(err)
	}

	return prt, nil
}

func (r *Repo) GetPasswordResetToken(resetUUID string) (*PasswordResetToken, error) {
	sqlStatement := "SELECT uuid,user_id,token,expiration FROM password_reset_tokens WHERE uuid=$1"
	row := r.dao.Conn.QueryRow(sqlStatement, resetUUID)
	var t PasswordResetToken
	err := row.Scan(&t.uuid, &t.userId, &t.token, &t.expiration)
	if err != nil {
		panic(err)
	}

	return &t, nil
}

func (r *Repo) DeletePasswordResetToken(resetUUID string) error {
	sqlStatement := "DELETE FROM password_reset_tokens WHERE uuid=$1"
	_, err := r.dao.Conn.Exec(sqlStatement, resetUUID)
	if err != nil {
		panic(err)
	}

	return nil
}
