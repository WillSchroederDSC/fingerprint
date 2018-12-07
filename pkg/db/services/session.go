package services

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/models"
	"github.com/willschroeder/fingerprint/pkg/proto"
)

type SessionService struct {
	repo *db.Repo
}

func NewSessionService(repo *db.Repo) *SessionService {
	return &SessionService{repo: repo}
}

func (ss *SessionService) CreateSession(userUUID string) (*models.Session, error) {
	session, err := ss.repo.CreateSession(userUUID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (ss *SessionService) AddTokenToSession(sessionUUID string, token string) error {
	err := ss.repo.UpdateSessionToken(sessionUUID, token)
	if err != nil {
		return err
	}

	return nil
}


func (ss *SessionService) BuildScopeGroupings(sessionUUID string, protoScopeGroupings []*proto.ScopeGrouping) ([]*models.ScopeGrouping, error) {
	scopeGroupings := make([]*models.ScopeGrouping, len(protoScopeGroupings))
	for i, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't convert timestamp")
		}

		scopeGrouping, err := ss.repo.CreateScopeGrouping(sessionUUID, sg.Scopes, exp)
		if err != nil {
			return nil, err
		}
		scopeGroupings[i] = scopeGrouping
	}

	return scopeGroupings, nil
}

