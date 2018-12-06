package db

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type Repo struct {
	Dao *DAO
}

func NewRepo(dao *DAO) *Repo {
	return &Repo{Dao:dao}
}

func (r *Repo) CreateUser(tx *sql.Tx, email string, encryptedPassword string, isGuest bool) (*User, error) {
	userUUID := uuid.New().String()

	sqlStatement := "INSERT INTO users (Uuid, email, encrypted_password, is_guest, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := tx.Exec(sqlStatement, userUUID, email, encryptedPassword, isGuest, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new user")
	}

	user, err := r.GetUserWithUUIDUsingTx(tx, userUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repo) GetUserWithUUID(userUUID string) (*User, error) {
	sqlStatement := "SELECT Uuid,email,encrypted_password,is_guest FROM users WHERE Uuid=$1"

	row := r.Dao.Conn.QueryRow(sqlStatement, userUUID)
	var user User
	err := row.Scan(&user.Uuid, &user.Email, &user.EncryptedPassword, &user.IsGuest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}

	return &user, nil
}

func (r *Repo) UpdateUserPassword(email string, encryptedPassword string) error {
	sqlStatement := "UPDATE users SET encrypted_password=$1,updated_at=$2 WHERE email=$3"
	_, err := r.Dao.Conn.Exec(sqlStatement, encryptedPassword, time.Now().UTC(), email)
	if err != nil {
		return errors.Wrap(err, "failed to update users password")
	}

	return nil
}

func (r *Repo) GetUserWithEmail(email string) (*User, error) {
	sqlStatement := "SELECT Uuid,email,encrypted_password,is_guest FROM users WHERE email=$1"

	row := r.Dao.Conn.QueryRow(sqlStatement, email)
	var user User
	err := row.Scan(&user.Uuid, &user.Email, &user.EncryptedPassword, &user.IsGuest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

func (r *Repo) GetUserWithUUIDUsingTx(tx *sql.Tx, userUUID string) (*User, error) {
	sqlStatement := "SELECT Uuid,email,encrypted_password,is_guest FROM users WHERE Uuid=$1"

	row := tx.QueryRow(sqlStatement, userUUID)
	var user User
	err := row.Scan(&user.Uuid, &user.Email, &user.EncryptedPassword, &user.IsGuest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

func (r *Repo) CreateSession(tx *sql.Tx, newSessionUUID string, userUUID string, token string) (*Session, error) {
	sqlStatement := "INSERT INTO sessions (Uuid, user_uuid, token, created_at) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(sqlStatement, newSessionUUID, userUUID, token, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	session, err := r.GetSessionWithUUIDUsingTx(tx, newSessionUUID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repo) GetSessionWithUUIDUsingTx(tx *sql.Tx, sessionUUID string) (*Session, error) {
	sqlStatement := "SELECT Uuid,token FROM sessions WHERE Uuid=$1"

	row := tx.QueryRow(sqlStatement, sessionUUID)
	var session Session
	err := row.Scan(&session.Uuid, &session.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	return &session, nil
}

func (r *Repo) GetSessionWithToken(token string) (*Session, error) {
	sqlStatement := "SELECT Uuid,token FROM sessions WHERE token=$1"

	row := r.Dao.Conn.QueryRow(sqlStatement, token)
	var session Session
	err := row.Scan(&session.Uuid, &session.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	return &session, nil
}

func (r *Repo) CreateScopeGrouping(tx *sql.Tx, sessionUUID string, scopes []string, expiration time.Time) (*ScopeGrouping, error) {
	groupingUUID := uuid.New().String()

	sqlStatement := "INSERT INTO scope_groupings (Uuid, session_uuid, scopes, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(sqlStatement, groupingUUID, sessionUUID, pq.Array(scopes), expiration, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scope grouping")
	}

	grouping, err := r.GetScopeGroupingWithUUID(tx, groupingUUID)
	if err != nil {
		return nil, err
	}

	return grouping, nil
}

func (r *Repo) GetScopeGroupingWithUUID(tx *sql.Tx, groupingUUID string) (*ScopeGrouping, error) {
	sqlStatement := "SELECT Uuid,scopes,expiration FROM scope_groupings WHERE Uuid=$1"
	row := tx.QueryRow(sqlStatement, groupingUUID)
	var sg ScopeGrouping
	err := row.Scan(&sg.Uuid, pq.Array(&sg.Scopes), &sg.Expiration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get scope grouping")
	}

	return &sg, nil
}

func (r *Repo) DeleteSessionWithUUID(sessionUUID string) error {
	sqlStatement := "DELETE FROM sessions WHERE Uuid=$1"
	_, err := r.Dao.Conn.Exec(sqlStatement, sessionUUID)
	if err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}

func (r *Repo) DeleteSessionWithToken(token string) interface{} {
	sqlStatement := "DELETE FROM sessions WHERE token=$1"
	_, err := r.Dao.Conn.Exec(sqlStatement, token)
	if err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}

func (r *Repo) CreatePasswordResetToken(userUUID string, expiration time.Time) (*PasswordResets, error) {
	resetUUID := uuid.New().String()
	newResetToken := RandomString(16)
	sqlStatement := "INSERT INTO password_resets (Uuid, user_uuid, token, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.Dao.Conn.Exec(sqlStatement, resetUUID, userUUID, newResetToken, expiration, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create password reset")
	}

	prt, err := r.GetPasswordResetToken(resetUUID)
	if err != nil {
		return nil, err
	}

	return prt, nil
}

func (r *Repo) GetPasswordResetToken(token string) (*PasswordResets, error) {
	sqlStatement := "SELECT Uuid,user_uuid,token,expiration FROM password_resets WHERE token=$1"
	row := r.Dao.Conn.QueryRow(sqlStatement, token)
	var t PasswordResets
	err := row.Scan(&t.Uuid, &t.UserUuid, &t.Token, &t.Expiration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get password reset")
	}

	return &t, nil
}

func (r *Repo) DeletePasswordResetToken(resetUUID string) error {
	sqlStatement := "DELETE FROM password_resets WHERE Uuid=$1"
	_, err := r.Dao.Conn.Exec(sqlStatement, resetUUID)
	if err != nil {
		return errors.Wrap(err, "failed to delete password reset")
	}

	return nil
}

func (r *Repo) DeleteAllPasswordResetTokensForUser(userUUID string) error {
	sqlStatement := "DELETE FROM password_reset_tokens WHERE user_uuid=$1"
	_, err := r.Dao.Conn.Exec(sqlStatement, userUUID)
	if err != nil {
		return errors.Wrap(err, "failed to delete password reset")
	}

	return nil
}
