package server

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
	"golang.org/x/crypto/bcrypt"
	"time"
	"errors"
	"github.com/google/uuid"
)

type Builder struct {
	repo *Repo
	dao *db.DAO
}

func (b *Builder) buildUser(tx *sql.Tx, password string, passwordConfirmation string, email string) (*User, error) {
	if password != passwordConfirmation {
		return nil, errors.New("password and confirmation don't match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := b.repo.CreateUser(tx, email, string(hash))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (b *Builder) buildSession(tx *sql.Tx, newSessionUUID uuid.UUID, userID int, sessionToken string, furthestExpiration time.Time) (*Session, error) {
	session, err := b.repo.CreateSession(tx, newSessionUUID, userID, sessionToken, furthestExpiration)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return session, nil
}

func (b *Builder) buildScopeGroupings(tx *sql.Tx, protoScopeGroupings []*proto.ScopeGrouping, sessionID int) ([]*ScopeGrouping, error) {
	scopeGroupings := make([]*ScopeGrouping, len(protoScopeGroupings))
	for i, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		scopeGrouping, err := b.repo.CreateScopeGrouping(tx, sessionID, sg.Scopes, exp)
		if err != nil {
			fmt.Println(err)
			return nil, err
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
			fmt.Println(err)
			return "", "", time.Now(), err
		}
		tf.AddScopeGrouping(sg.Scopes, exp)
	}

	sess, err := tf.GenerateSession()
	if err != nil {
		fmt.Println(err)
		return "", "", time.Now(), err
	}

	return  sess.Token, sess.Json, sess.FurthestExpiration, nil
}
