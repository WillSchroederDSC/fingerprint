package server

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/willschroeder/fingerprint/pkg/db"
	"time"
)

type Repo struct {
	dao *db.DAO
}

func (r *Repo) CreateUser(email string, encryptedPassword string) (*User, error) {
	userUUID := uuid.New().String()

	sqlStatement := "INSERT INTO users (uuid, email, encrypted_password, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.dao.Conn.Exec(sqlStatement, userUUID, email, encryptedPassword, time.Now().UTC())
	if err != nil {
		panic(err)
		return nil, err
	}

	user, err := r.GetUserWithUUID(userUUID)
	if err != nil {
		panic(err)
		return nil, err
	}

	return user, nil
}

func (r *Repo) GetUserWithUUID(userUUID string) (*User, error) {
	sqlStatement := "SELECT id,uuid,email FROM users WHERE uuid=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, userUUID)
	var user User
	err := row.Scan(&user.id, &user.uuid, &user.email)
	if err != nil {
		panic(err)
		return nil, err
	}

	return &user, nil
}

func (r *Repo) CreateSession(userId int, token string, expiration time.Time) (*Session, error) {
	sessionUUID := uuid.New().String()

	sqlStatement := "INSERT INTO sessions (uuid, user_id, token, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.dao.Conn.Exec(sqlStatement, sessionUUID, userId, token, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		panic(err)
		return nil, err
	}

	session, err := r.GetSessionWithUUID(sessionUUID)
	if err != nil {
		panic(err)
		return nil, err
	}

	return session, nil
}

func (r *Repo) GetSessionWithUUID(sessionUUID string) (*Session, error) {
	sqlStatement := "SELECT id,uuid,token,expiration FROM sessions WHERE uuid=$1"

	row := r.dao.Conn.QueryRow(sqlStatement, sessionUUID)
	var session Session
	err := row.Scan(&session.id, &session.uuid, &session.token, &session.expiration)
	if err != nil {
		panic(err)
		return nil, err
	}

	return &session, nil
}

func (r *Repo) CreateScopeGrouping(sessionId int, scopes []string, expiration time.Time) (*ScopeGrouping, error) {
	groupingUUID := uuid.New().String()

	sqlStatement := "INSERT INTO scope_groupings (uuid, session_id, scopes, expiration, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.dao.Conn.Exec(sqlStatement, groupingUUID, sessionId, pq.Array(scopes), expiration, time.Now().UTC())
	if err != nil {
		panic(err)
		return nil, err
	}
	grouping, err := r.GetScopeGroupingWithUUID(groupingUUID)
	if err != nil {
		panic(err)
		return nil, err
	}

	return grouping, nil
}

func (r *Repo) GetScopeGroupingWithUUID(groupingUUID string) (*ScopeGrouping, error) {
	sqlStatement := "SELECT id,uuid,scopes,expiration FROM scope_groupings WHERE uuid=$1"
	row := r.dao.Conn.QueryRow(sqlStatement, groupingUUID)
	var sg ScopeGrouping
	err := row.Scan(&sg.id,&sg.uuid,pq.Array(&sg.scopes),&sg.expiration)
	if err != nil {
		panic(err)
		return nil, err
	}

	return &sg, nil
}