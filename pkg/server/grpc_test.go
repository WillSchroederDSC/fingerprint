package server

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
)

var testDAO *db.DAO

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	testDAO = db.ConnectToDatabase()
	defer testDAO.DB.Close()
	code := m.Run()
	_, _ = testDAO.DB.Exec("TRUNCATE users CASCADE")
	os.Exit(code)
}

func TestGRPCServer_CreateUser(t *testing.T) {
	oneHour, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * time.Duration(1)))
	twoHour, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * time.Duration(2)))
	scopeGroupings := []*proto.ScopeGrouping{
		{
			Scopes:     []string{"read"},
			Expiration: oneHour,
		},
		{
			Scopes:     []string{"write"},
			Expiration: twoHour,
		},
	}

	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.CreateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "creates a new user",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.CreateUserRequest{
				Email: gofakeit.Email(), Password: "test", PasswordConfirmation: "test", ScopeGroupings: scopeGroupings },
			},
			wantErr: false,
		},
	}
	for _, tt := range tests{
		t.Run(tt.name, func (t *testing.T){
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreateUser(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr{
				t.Errorf("GRPCServer.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.User.Uuid == "" || got.Session.Uuid == "" || got.User.Email == "" {
				t.Error("User was not created with email")
			}
		})
	}
}


func TestGRPCServer_CreateGuestUser(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.CreateGuestUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.CreateGuestUserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreateGuestUser(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateGuestUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateGuestUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_GetUser(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.GetUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.GetUserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.GetUser(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_CreatePasswordResetToken(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.CreatePasswordResetTokenRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.CreatePasswordResetTokenResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreatePasswordResetToken(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreatePasswordResetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreatePasswordResetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_UpdateUserPassword(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.ResetUserPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.ResetUserPasswordResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.UpdateUserPassword(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.UpdateUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.UpdateUserPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_CreateSession(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.CreateSessionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.CreateSessionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreateSession(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_GetSession(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.GetSessionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.GetSessionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.GetSession(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_DeleteSession(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.DeleteSessionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.DeleteSessionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.DeleteSession(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.DeleteSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_DeleteUser(t *testing.T) {
	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.DeleteUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.DeleteUserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.DeleteUser(tt.args.in0, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}