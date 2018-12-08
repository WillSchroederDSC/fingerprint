package db

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/models"
	"github.com/willschroeder/fingerprint/pkg/util"
	"time"
)

type Repo struct {
	db *sql.DB
	tx *sql.Tx
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func NewRepoUsingTransaction(tx *sql.Tx) *Repo {
	return &Repo{tx: tx}
}

func (r *Repo) exec(query string, args ...interface{}) (sql.Result, error) {
	if r.tx != nil {
		return r.tx.Exec(query)
	}

	return r.db.Exec(query)
}

func (r *Repo) queryRow(query string, args ...interface{}) *sql.Row {
	if r.tx != nil {
		return r.tx.QueryRow(query)
	}

	return r.db.QueryRow(query)
}

func (r *Repo) CreateUser(email string, encryptedPassword string, isGuest bool) (*models.User, error) {
	userUUID := uuid.New().String()

	sqlStatement := "INSERT INTO users (Uuid, email, encrypted_password, is_guest, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.exec(sqlStatement, userUUID, email, encryptedPassword, isGuest, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new user")
	}

	user, err := r.GetUserWithUUID(userUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repo) GetUserWithUUID(userUUID string) (*models.User, error) {
	sqlStatement := "SELECT Uuid,email,encrypted_password,is_guest FROM users WHERE Uuid=$1"

	row := r.queryRow(sqlStatement, userUUID)
	var user models.User
	err := row.Scan(&user.Uuid, &user.Email, &user.EncryptedPassword, &user.IsGuest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}

	return &user, nil
}

func (r *Repo) UpdateUserPassword(email string, encryptedPassword string) error {
	sqlStatement := "UPDATE users SET encrypted_password=$1,updated_at=$2 WHERE email=$3"
	_, err := r.exec(sqlStatement, encryptedPassword, time.Now().UTC(), email)
	if err != nil {
		return errors.Wrap(err, "failed to update users password")
	}

	return nil
}

func (r *Repo) GetUserWithEmail(email string) (*models.User, error) {
	sqlStatement := "SELECT Uuid,email,encrypted_password,is_guest FROM users WHERE email=$1"

	row := r.queryRow(sqlStatement, email)
	var user models.User
	err := row.Scan(&user.Uuid, &user.Email, &user.EncryptedPassword, &user.IsGuest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

func (r *Repo) CreateSession(userUUID string) (*models.Session, error) {
	sessionUUID := uuid.New().String()

	sqlStatement := "INSERT INTO sessions (Uuid, user_uuid, token, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.exec(sqlStatement, sessionUUID, userUUID, "TEMP", time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	session, err := r.GetSessionWithUUID(sessionUUID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repo) GetSessionWithUUID(sessionUUID string) (*models.Session, error) {
	sqlStatement := "SELECT Uuid,token FROM sessions WHERE Uuid=$1"

	row := r.queryRow(sqlStatement, sessionUUID)
	var session models.Session
	err := row.Scan(&session.Uuid, &session.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	return &session, nil
}

func (r *Repo) GetSessionWithToken(token string) (*models.Session, error) {
	sqlStatement := "SELECT Uuid,token FROM sessions WHERE token=$1"

	row := r.queryRow(sqlStatement, token)
	var session models.Session
	err := row.Scan(&session.Uuid, &session.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	return &session, nil
}

func (r *Repo) CreateScopeGrouping(sessionUUID string, scopes []string, expiration time.Time) (*models.ScopeGrouping, error) {
	groupingUUID := uuid.New().String()

	sqlStatement := "INSERT INTO scope_groupings (Uuid, session_uuid, scopes, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.exec(sqlStatement, groupingUUID, sessionUUID, pq.Array(scopes), expiration, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scope grouping")
	}

	grouping, err := r.GetScopeGroupingWithUUID(groupingUUID)
	if err != nil {
		return nil, err
	}

	return grouping, nil
}

func (r *Repo) GetScopeGroupingWithUUID(groupingUUID string) (*models.ScopeGrouping, error) {
	sqlStatement := "SELECT Uuid,scopes,expiration FROM scope_groupings WHERE Uuid=$1"
	row := r.queryRow(sqlStatement, groupingUUID)
	var sg models.ScopeGrouping
	err := row.Scan(&sg.Uuid, pq.Array(&sg.Scopes), &sg.Expiration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get scope grouping")
	}

	return &sg, nil
}

func (r *Repo) UpdateSessionToken(sessionUUID string, token string) error {
	sqlStatement := "UPDATE sessions SET token=$1 WHERE uuid=$2"

	_, err := r.exec(sqlStatement, token, sessionUUID)
	if err != nil {
		return errors.Wrap(err, "failed to update session token")
	}

	return nil
}

func (r *Repo) DeleteSessionWithUUID(sessionUUID string) error {
	sqlStatement := "DELETE FROM sessions WHERE Uuid=$1"
	_, err := r.exec(sqlStatement, sessionUUID)
	if err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}

func (r *Repo) DeleteSessionWithToken(token string) error {
	sqlStatement := "DELETE FROM sessions WHERE token=$1"
	_, err := r.exec(sqlStatement, token)
	if err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}

func (r *Repo) CreatePasswordResetToken(userUUID string, expiration time.Time) (*models.PasswordReset, error) {
	resetUUID := uuid.New().String()
	newResetToken := util.String(16)
	sqlStatement := "INSERT INTO password_resets (Uuid, user_uuid, token, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.exec(sqlStatement, resetUUID, userUUID, newResetToken, expiration, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create password reset")
	}

	prt, err := r.GetPasswordResetToken(resetUUID)
	if err != nil {
		return nil, err
	}

	return prt, nil
}

func (r *Repo) GetPasswordResetToken(token string) (*models.PasswordReset, error) {
	sqlStatement := "SELECT Uuid,user_uuid,token,expiration FROM password_resets WHERE token=$1"
	row := r.queryRow(sqlStatement, token)
	var t models.PasswordReset
	err := row.Scan(&t.Uuid, &t.UserUuid, &t.Token, &t.Expiration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get password reset")
	}

	return &t, nil
}

func (r *Repo) DeletePasswordResetToken(resetUUID string) error {
	sqlStatement := "DELETE FROM password_resets WHERE Uuid=$1"
	_, err := r.exec(sqlStatement, resetUUID)
	if err != nil {
		return errors.Wrap(err, "failed to delete password reset")
	}

	return nil
}

func (r *Repo) DeleteAllPasswordResetTokensForUser(userUUID string) error {
	sqlStatement := "DELETE FROM password_reset_tokens WHERE user_uuid=$1"
	_, err := r.exec(sqlStatement, userUUID)
	if err != nil {
		return errors.Wrap(err, "failed to delete password reset")
	}

	return nil
}

func (r *Repo) DeleteUser(email string) error {
	sqlStatement := "DELETE FROM users WHERE email=$1"
	_, err := r.exec(sqlStatement, email)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}
