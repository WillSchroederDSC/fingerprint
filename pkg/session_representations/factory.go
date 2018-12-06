package session_representations

import (
	"encoding/json"
	"errors"
	"github.com/o1egl/paseto"
	"time"
)

type Factory struct {
	Version        int                          `json:"version"`
	CustomerUUID   string                       `json:"customer_uuid"`
	SessionUUID    string                       `json:"session_uuid"`
	ScopeGroupings []*tokenFactoryScopeGrouping `json:"scope_groupings"`
}

func (tf *Factory) Valid() error {
	if len(tf.ScopeGroupings) < 1 {
		return errors.New("must have at least one scope grouping")
	}

	return nil
}

type tokenFactoryScopeGrouping struct {
	Scopes     []string  `json:"scopes"`
	Expiration time.Time `json:"expiration"`
}

type Representations struct {
	Token              string
	Json               string
}

func NewTokenFactory(userUUID string, sessionUUID string) *Factory {
	return &Factory{Version: 1, CustomerUUID: userUUID, SessionUUID: sessionUUID}
}

func (tf *Factory) AddScopeGrouping(scopes []string, expiration time.Time) {
	tf.ScopeGroupings = append(tf.ScopeGroupings, &tokenFactoryScopeGrouping{Scopes: scopes, Expiration: expiration})
}

func (tf *Factory) GenerateSession() (*Representations, error) {
	err := tf.Valid()
	if err != nil {
		return nil, err
	}

	token, err := tf.generateToken()
	if err != nil {
		return nil, err
	}

	jsonStr, err := tf.generateJSON()
	if err != nil {
		return nil, err
	}

	return &Representations{Token: token, Json: jsonStr}, nil
}

func (tf *Factory) generateToken() (string, error) {
	v2 := paseto.NewV2()
	token, err := v2.Encrypt(secret(), tf, "")
	if err != nil {
		panic(err)
	}

	return token, nil
}

func (tf *Factory) generateJSON() (string, error) {
	bs, err := json.Marshal(tf)
	if err != nil {
		panic(err)
	}
	return string(bs), nil
}

//func (tf *Factory) findFurthestExpiration() time.Time {
//	furthest := tf.ScopeGroupings[0].Expiration
//	for _, sg := range tf.ScopeGroupings {
//		if sg.Expiration.After(furthest) {
//			furthest = sg.Expiration
//		}
//	}
//	return furthest
//}

func secret() []byte {
	// MUST be 32 chars
	return []byte("YELLOW SUBMARINE, BLACK WIZARDRY")
}

func DecodeTokenToJson(sessionToken string) string {
	v2 := paseto.NewV2()
	var token string
	var footer string
	err := v2.Decrypt(sessionToken, secret(), &token, &footer)
	if err != nil {
		panic(err)
	}
	return token
}
