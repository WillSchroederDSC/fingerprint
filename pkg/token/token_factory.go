package token

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Encoder struct {
	Version        int    `json:"version"`
	CustomerUUID   string `json:"customer_id"`
	SessionUUID    string `json:"session_id"`
	IsGuest        bool   `json:"is_guest"`
	ScopeGroupings []*tokenFactoryScopeGrouping `json:"scope_groupings"`
}

// Required to let object be encoded by JWT lib
func (_ Encoder) Valid() error {
	return nil
}

type tokenFactoryScopeGrouping struct {
	Scopes     []string `json:"scopes"`
	Expiration time.Time `json:"expiration"`
}

func (tf *Encoder) AddScopeGrouping(scopes []string, expiration time.Time) {
	tf.ScopeGroupings = append(tf.ScopeGroupings, &tokenFactoryScopeGrouping{Scopes: scopes, Expiration:expiration})
}

func (tf *Encoder) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = tf
	tokenString, err := token.SignedString(secret())
	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

func (tf *Encoder) GenerateJSON() (string, error) {
	bs, err := json.Marshal(tf)
	if err != nil {
		panic(err)
	}
	return string(bs), nil
}

func BuildTokenFactory(customerUUID string, sessionUUID string, isGuest bool) *Encoder {
	return &Encoder{Version: 1, CustomerUUID: customerUUID, SessionUUID:sessionUUID, IsGuest:isGuest}
}

func secret() []byte {
	return []byte("foobar")
}