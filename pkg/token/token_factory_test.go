package token

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	factory := NewTokenFactory("111", "222", true)
	factory.AddScopeGrouping([]string{"read", "write"}, time.Now())
	factory.AddScopeGrouping([]string{"test", "another"}, time.Now())
	json, _ := factory.GenerateJSON()
	println(json)
	token, _ := factory.GenerateToken()
	println(token)
	DecodeTokenToJson(token)
}