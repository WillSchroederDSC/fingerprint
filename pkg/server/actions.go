package server

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/models"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/random"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
	"golang.org/x/crypto/bcrypt"
)

type Actions struct {
	repo *db.Repo
	dao  *db.DAO
}

func NewActions(dao *db.DAO) *Actions {
	return &Actions{dao:dao, repo:db.NewRepo(dao)}
}

func BuildPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt password")
	}

	return string(hash), nil
}

func confirmPasswordAndHash(password string, passwordConfirmation string) (string, error) {
	if password != passwordConfirmation {
		return "", errors.New("password and confirmation don't match")
	}

	hash, err := BuildPasswordHash(password)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (b *Actions) buildUser(tx *sql.Tx, email string, password string, passwordConfirmation string) (*models.User, error) {
	hash, err := confirmPasswordAndHash(password, passwordConfirmation)
	if err != nil {
		return nil, err
	}

	user, err := b.repo.CreateUser(tx, email, hash, false)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *Actions) updateUserPassword(email string, passwordResetToken string, password string, passwordConfirmation string) error {
	prt, err := b.repo.GetPasswordResetToken(passwordResetToken)
	if err != nil {
		return err
	}

	hash, err := confirmPasswordAndHash(password, passwordConfirmation)
	if err != nil {
		return err
	}

	if passwordResetToken != prt.Token {
		return errors.New("current user reset token does not match given reset token")
	}

	err = b.repo.DeletePasswordResetToken(prt.Uuid)
	if err != nil {
		return err
	}

	err = b.repo.UpdateUserPassword(email, hash)
	if err != nil {
		return err
	}

	return nil
}

func (b *Actions) buildGuestUser(tx *sql.Tx, email string) (*models.User, error) {
	hash, err := BuildPasswordHash(random.String(16))
	if err != nil {
		return nil, err
	}

	email = email + "." + random.String(6) + ".guest"

	user, err := b.repo.CreateUser(tx, email, hash, true)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *Actions) buildSession(tx *sql.Tx, newSessionUUID string, userUUID string, sessionToken string) (*models.Session, error) {
	session, err := b.repo.CreateSession(tx, newSessionUUID, userUUID, sessionToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (b *Actions) buildScopeGroupings(tx *sql.Tx, protoScopeGroupings []*proto.ScopeGrouping, sessionUUID string) ([]*models.ScopeGrouping, error) {
	scopeGroupings := make([]*models.ScopeGrouping, len(protoScopeGroupings))
	for i, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't convert timestamp")
		}

		scopeGrouping, err := b.repo.CreateScopeGrouping(tx, sessionUUID, sg.Scopes, exp)
		if err != nil {
			return nil, err
		}
		scopeGroupings[i] = scopeGrouping
	}

	return scopeGroupings, nil
}

func (b *Actions) buildSessionRepresentation(user *models.User, sessionUUID string, protoScopeGroupings []*proto.ScopeGrouping) (tokenStr string, json string, err error) {
	tf := session_representations.NewTokenFactory(user.Uuid, sessionUUID)
	for _, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			return "", "", errors.Wrap(err, "couldn't convert timestamp")
		}
		tf.AddScopeGrouping(sg.Scopes, exp)
	}

	sess, err := tf.GenerateSession()
	if err != nil {
		return "", "", err
	}

	return sess.Token, sess.Json, nil
}
