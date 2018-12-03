package server

import (
	"database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Builder struct {
	repo *Repo
	dao *db.DAO
}

func BuildPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hash), nil
}

func (b *Builder) buildUser(tx *sql.Tx, password string, passwordConfirmation string, email string) (*User, error) {
	if password != passwordConfirmation {
		return nil, errors.New("password and confirmation don't match")
	}

	hash, err := BuildPasswordHash(password)
	if err != nil {
		panic(err)
	}

	user, err := b.repo.CreateUser(tx, email, hash)
	if err != nil {
		panic(err)
	}

	return user, nil
}

func (b *Builder) buildSession(tx *sql.Tx, newSessionUUID uuid.UUID, userID int, sessionToken string, furthestExpiration time.Time) (*Session, error) {
	session, err := b.repo.CreateSession(tx, newSessionUUID, userID, sessionToken, furthestExpiration)
	if err != nil {
		panic(err)
	}

	return session, nil
}

func (b *Builder) buildScopeGroupings(tx *sql.Tx, protoScopeGroupings []*proto.ScopeGrouping, sessionID int) ([]*ScopeGrouping, error) {
	scopeGroupings := make([]*ScopeGrouping, len(protoScopeGroupings))
	for i, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			panic(err)
		}

		scopeGrouping, err := b.repo.CreateScopeGrouping(tx, sessionID, sg.Scopes, exp)
		if err != nil {
			panic(err)
		}
		scopeGroupings[i] = scopeGrouping
	}

	return scopeGroupings, nil
}

func (b *Builder) buildToken(user *User, sessionUUID uuid.UUID, protoScopeGroupings []*proto.ScopeGrouping) (tokenStr string, json string, furthestExpiration time.Time, err error) {
	tf := session_representations.NewTokenFactory(user.uuid, sessionUUID.String())
	for _, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			panic(err)
		}
		tf.AddScopeGrouping(sg.Scopes, exp)
	}

	sess, err := tf.GenerateSession()
	if err != nil {
		panic(err)
	}

	return  sess.Token, sess.Json, sess.FurthestExpiration, nil
}
