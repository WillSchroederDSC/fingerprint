package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"time"
)

type User struct {
	Uuid              string
	Email             string
	EncryptedPassword string
	IsGuest           bool
}

func (u *User) ConvertToProtobuff() *proto.User {
	return &proto.User{
		Uuid:  u.Uuid,
		Email: u.Email,
	}
}

type Session struct {
	Uuid     string
	Token    string
	UserUuid string
}

func (s *Session) ConvertToProtobuff(json string) *proto.Session {
	return &proto.Session{
		Uuid:  s.Uuid,
		Token: s.Token,
		Json:  json,
	}
}

type ScopeGrouping struct {
	Uuid        string
	SessionUuid string
	Scopes      []string
	Expiration  time.Time
}

func (sg *ScopeGrouping) ConvertToProtobuff() (*proto.ScopeGrouping, error) {
	timestamp, err := ptypes.TimestampProto(sg.Expiration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode time")
	}

	return &proto.ScopeGrouping{
		Scopes:     sg.Scopes,
		Expiration: timestamp,
	}, nil
}

type PasswordReset struct {
	Uuid       string
	UserUuid   string
	Token      string
	Expiration time.Time
}

type SessionRepresentation struct {
	Token string
	Json  string
}
