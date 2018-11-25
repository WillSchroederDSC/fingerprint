package server

import "github.com/willschroeder/fingerprint/pkg/proto"

type User struct {
	uuid string
	email string
}

func (u *User) ConvertToProtobuff() *proto.User {
	return &proto.User{
		Uuid: u.uuid,
		Email: u.email,
	}
}