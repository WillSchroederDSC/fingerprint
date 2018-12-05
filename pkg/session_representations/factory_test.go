package session_representations

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	factory := NewTokenFactory("111", "222")
	factory.AddScopeGrouping([]string{"read", "write"}, time.Now())
	factory.AddScopeGrouping([]string{"test", "another"}, time.Now())
	session, _ := factory.GenerateSession()
	println(session.Json)
	println(session.Token)
	DecodeTokenToJson(session.Token)
}
