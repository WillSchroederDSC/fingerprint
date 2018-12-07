package services

import (
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/models"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/random"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	dao *db.DAO
}

func NewUserService(dao *db.DAO) *UserService {
	return &UserService{dao: dao}
}

func (us *UserService) CreateUser(request *proto.CreateUserRequest) (*models.User, *models.Session, error) {
	hashedPassword, err := confirmPasswordAndHash(request.Password, request.PasswordConfirmation)
	if err != nil {
		return nil, nil, err
	}

	tx, err := us.dao.NewTransaction()
	if err != nil {
		return nil, nil, err
	}
	repo := db.NewRepoUsingTransaction(tx)

	user, err := repo.CreateUser(request.Email, hashedPassword, false)
	if err != nil {
		db.HandleRollback(tx)
		return nil, nil, err
	}

	session, err := buildSession(repo, user.Uuid, request.ScopeGroupings)
	if err != nil {
		db.HandleRollback(tx)
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (us *UserService) CreateGuestUser(request *proto.CreateGuestUserRequest) (*models.User, *models.Session, error) {
	hash, err := buildPasswordHash(random.String(16))
	if err != nil {
		return nil, nil, err
	}

	email := request.Email + "." + random.String(6) + ".guest"

	tx, err := us.dao.NewTransaction()
	if err != nil {
		return nil, nil, err
	}
	repo := db.NewRepoUsingTransaction(tx)

	user, err := repo.CreateUser(email, hash, true)
	if err != nil {
		db.HandleRollback(tx)
		return nil, nil, err
	}

	session, err := buildSession(repo, user.Uuid, request.ScopeGroupings)
	if err != nil {
		db.HandleRollback(tx)
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (us *UserService) UpdateUserPassword(email string, passwordResetToken string, password string, passwordConfirmation string) error {
	repo := db.NewRepo(us.dao.DB)

	prt, err := repo.GetPasswordResetToken(passwordResetToken)
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

	err = repo.DeletePasswordResetToken(prt.Uuid)
	if err != nil {
		return err
	}

	err = repo.UpdateUserPassword(email, hash)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) DeleteUser() error {
	panic("build me")
}

func confirmPasswordAndHash(password string, passwordConfirmation string) (string, error) {
	if password != passwordConfirmation {
		return "", errors.New("password and confirmation don't match")
	}

	hash, err := buildPasswordHash(password)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func buildPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt password")
	}

	return string(hash), nil
}

func buildSession(repo * db.Repo, userUUID string, groupings []*proto.ScopeGrouping) (*models.Session, error) {
	sessionService := NewSessionService(repo)
	session, err := sessionService.CreateSession(userUUID)
	if err != nil {
		return nil, err
	}

	scopeGroupings, err := sessionService.BuildScopeGroupings(session.Uuid, groupings)
	if err != nil {
		return nil, err
	}

	representationsService := NewSessionRepresentationService(userUUID, session.Uuid)
	for _, sg := range scopeGroupings {
		representationsService.AddScopeGrouping(sg.Scopes, sg.Expiration)
	}

	representation, err := representationsService.GenerateSession()
	if err != nil {
		return nil, err
	}

	err = sessionService.AddTokenToSession(session.Uuid, representation.Token)
	if err != nil {
		return nil, err
	}

	return session, nil
}