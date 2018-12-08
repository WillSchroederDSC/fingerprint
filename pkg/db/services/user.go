package services

import (
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/models"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	dao *db.DAO
}

func NewUserService(dao *db.DAO) *UserService {
	return &UserService{dao: dao}
}

func (us *UserService) ValidateEmailAndPassword(email string, password string) (*models.User, error) {
	repo := db.NewRepo(us.dao.DB)
	user, err := repo.GetUserWithEmail(email)
	if err != nil {
		return nil, err
	}

	_, err = validPasswordProvided(user, password)
	if err != nil {
		return nil, err
	}

	return user, nil
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

	sessionService := NewSessionService(repo)
	session, err := sessionService.CreateSession(user.Uuid, request.ScopeGroupings)
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
	hash, err := buildPasswordHash(util.String(16))
	if err != nil {
		return nil, nil, err
	}

	email := request.Email + "." + util.String(6) + ".guest"

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

	sessionService := NewSessionService(repo)
	session, err := sessionService.CreateSession(user.Uuid, request.ScopeGroupings)
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

func (us *UserService) CreatePasswordResetToken(email string, expiration time.Time) (*models.PasswordReset, error) {
	repo := db.NewRepo(us.dao.DB)

	user, err := repo.GetUserWithEmail(email)
	if err != nil {
		return nil, err
	}

	err = repo.DeleteAllPasswordResetTokensForUser(user.Uuid)
	if err != nil {
		return nil, err
	}

	resetToken, err := repo.CreatePasswordResetToken(user.Uuid, expiration)
	if err != nil {
		return nil, err
	}

	return resetToken, nil
}

func (us *UserService) DeleteUser(email string, password string) error {
	repo := db.NewRepo(us.dao.DB)

	user, err := repo.GetUserWithEmail(email)
	if err != nil {
		return err
	}

	_, err = validPasswordProvided(user, password)
	if err != nil {
		return err
	}

	err = repo.DeleteUser(email)
	if err != nil {
		return err
	}

	return nil
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

func validPasswordProvided(user *models.User, password string) (bool, error) {
	hashedPassword, err := buildPasswordHash(password)
	if err != nil {
		return false, err
	}

	if user.EncryptedPassword != hashedPassword {
		return false, errors.New("provided password incorrect")
	}

	return true, nil
}