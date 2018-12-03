package server

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"time"
)

type User struct {
	id int
	uuid string
	email string
	encryptedPassword string
}

func (u *User) ConvertToProtobuff() *proto.User {
	return &proto.User{
		Uuid: u.uuid,
		Email: u.email,
	}
}

type Session struct {
	id int
	uuid string
	token string
	customerId int
	expiration time.Time
}

func (s *Session) ConvertToProtobuff(json string) *proto.Session {
	return &proto.Session{
		Uuid: s.uuid,
		Token: s.token,
		Json: json,
	}
}

type ScopeGrouping struct {
	id int
	uuid string
	sessionId int
	scopes []string
	expiration time.Time
}

func (sg *ScopeGrouping) ConvertToProtobuff() (*proto.ScopeGrouping, error) {
	timestamp, err := ptypes.TimestampProto(sg.expiration)
	if err != nil {
		return nil, err
	}

	return &proto.ScopeGrouping{
		Scopes: sg.scopes,
		Expiration: timestamp,
	}, nil
}