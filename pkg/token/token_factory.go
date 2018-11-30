package token

import (
	"encoding/json"
	"github.com/o1egl/paseto"
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
	v2 := paseto.NewV2()
	token, err := v2.Encrypt(secret(), tf, "")
	if err != nil {
		panic(err)
	}

	return token, nil
}

func (tf *Encoder) GenerateJSON() (string, error) {
	bs, err := json.Marshal(tf)
	if err != nil {
		panic(err)
	}
	return string(bs), nil
}

func BuildTokenFactory(userUUID string, sessionUUID string, isGuest bool) *Encoder {
	return &Encoder{Version: 1, CustomerUUID: userUUID, SessionUUID:sessionUUID, IsGuest:isGuest}
}

func secret() []byte {
	// MUST be 32 chars
	return []byte("YELLOW SUBMARINE, BLACK WIZARDRY")
}

func DecodeToken(sessionToken string) string {
	v2 := paseto.NewV2()
	var token string
	var footer string
	err := v2.Decrypt(sessionToken, secret(), &token, &footer)
	if err != nil {
		panic(err)
	}
	return token
}