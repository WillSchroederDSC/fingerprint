package server

import (
	"github.com/willschroeder/fingerprint/pkg/proto"
	"time"
)

type User struct {
	id int
	uuid string
	email string
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
	customerId int
	expiration time.Time
}

func (s *Session) ConvertToProtobuff(token string, json string) *proto.Session {
	return &proto.Session{
		Uuid: s.uuid,
		Token: token,
		Json: json,
	}
}

type ScopeGrouping struct {

}