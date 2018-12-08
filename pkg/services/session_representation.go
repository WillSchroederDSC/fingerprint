package services

import (
	"encoding/json"
	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/models"
	"time"
)

type SessionRepresentationService struct {
	Version        int                          `json:"version"`
	CustomerUUID   string                       `json:"customer_uuid"`
	SessionUUID    string                       `json:"session_uuid"`
	ScopeGroupings []*tokenFactoryScopeGrouping `json:"scope_groupings"`
}

func (srs *SessionRepresentationService) Valid() error {
	if len(srs.ScopeGroupings) < 1 {
		return errors.New("must have at least one scope grouping")
	}

	return nil
}

type tokenFactoryScopeGrouping struct {
	Scopes     []string  `json:"scopes"`
	Expiration time.Time `json:"expiration"`
}

func NewSessionRepresentationService(userUUID string, sessionUUID string) *SessionRepresentationService {
	return &SessionRepresentationService{Version: 1, CustomerUUID: userUUID, SessionUUID: sessionUUID}
}

func (srs *SessionRepresentationService) AddScopeGrouping(scopes []string, expiration time.Time) {
	srs.ScopeGroupings = append(srs.ScopeGroupings, &tokenFactoryScopeGrouping{Scopes: scopes, Expiration: expiration})
}

func (srs *SessionRepresentationService) GenerateSession() (*models.SessionRepresentation, error) {
	err := srs.Valid()
	if err != nil {
		return nil, err
	}

	token, err := srs.generateToken()
	if err != nil {
		return nil, err
	}

	jsonStr, err := srs.generateJSON()
	if err != nil {
		return nil, err
	}

	return &models.SessionRepresentation{Token: token, Json: jsonStr}, nil
}

func (srs *SessionRepresentationService) generateToken() (string, error) {
	v2 := paseto.NewV2()
	token, err := v2.Encrypt(secret(), srs, "")
	if err != nil {
		return "", errors.Wrap(err, "failed to generate paseto token")
	}

	return token, nil
}

func (srs *SessionRepresentationService) generateJSON() (string, error) {
	bs, err := json.Marshal(srs)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate json")
	}
	return string(bs), nil
}

//func (srs *SessionRepresentationService) findFurthestExpiration() time.Time {
//	furthest := srs.ScopeGroupings[0].Expiration
//	for _, sg := range srs.ScopeGroupings {
//		if sg.Expiration.After(furthest) {
//			furthest = sg.Expiration
//		}
//	}
//	return furthest
//}

func secret() []byte {
	// MUST be 32 chars
	// TODO: Make this configurable
	return []byte("YELLOW SUBMARINE, BLACK WIZARDRY")
}

func DecodeTokenToJson(sessionToken string) (string, error) {
	v2 := paseto.NewV2()
	var token string
	var footer string
	err := v2.Decrypt(sessionToken, secret(), &token, &footer)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt paseto token")
	}
	return token, nil
}
