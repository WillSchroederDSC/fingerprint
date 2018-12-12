package server

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"os"
	"testing"
	"time"

	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
)

var testDAO *db.DAO
var testGRPCServer *GRPCServer

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	testDAO = db.ConnectToDatabase()
	defer testDAO.DB.Close()
	testGRPCServer = NewGRPCServer(testDAO)
	code := m.Run()
	//_, _ = testDAO.DB.Exec("TRUNCATE users CASCADE")
	os.Exit(code)
}

func timestampOneHour() *timestamp.Timestamp {
	oneHour, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * time.Duration(1)))
	return oneHour
}

func timestampTwoHour() *timestamp.Timestamp {
	oneHour, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * time.Duration(2)))
	return oneHour
}

func twoScopeGroupings() []*proto.ScopeGrouping {
	return []*proto.ScopeGrouping{
		{
			Scopes:     []string{"read"},
			Expiration: timestampTwoHour(),
		},
		{
			Scopes:     []string{"write"},
			Expiration: timestampOneHour(),
		},
	}
}

func buildTestUser(password string) *proto.CreateUserResponse {
	usr, _ := testGRPCServer.CreateUser(context.Background(), &proto.CreateUserRequest{
		Email: gofakeit.Email(), Password: password, PasswordConfirmation: password, ScopeGroupings: twoScopeGroupings()})

	return usr
}

func TestGRPCServer_CreateUser(t *testing.T) {
	password := gofakeit.Password(true, true, true, true, true, 10)
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
				Email: gofakeit.Email(), Password: password, PasswordConfirmation: password, ScopeGroupings: twoScopeGroupings()},
			},
			wantErr: false,
		},
		{
			name:   "wont create a new user with mismatching password confirmation",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.CreateUserRequest{
				Email: gofakeit.Email(), Password: password, PasswordConfirmation: "wrong", ScopeGroupings: twoScopeGroupings()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreateUser(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			if got.User.Uuid == "" || got.Session.Uuid == "" || got.User.Email == "" {
				t.Error("User was not created with email")
			}
			if got.Session.Token == "" || got.Session.Json == "" {
				t.Error("didn't generate session representations")
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
		wantErr bool
	}{
		{
			name:    "creates a new user",
			fields:  fields{testDAO},
			args:    args{in0: context.Background(), request: &proto.CreateGuestUserRequest{Email: gofakeit.Email(), ScopeGroupings: twoScopeGroupings()}},
			wantErr: false,
		},
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
			if got.User.Uuid == "" || got.Session.Uuid == "" || got.User.Email == "" {
				t.Error("User was not created with email")
			}
			if got.Session.Token == "" || got.Session.Json == "" {
				t.Error("didn't generate session representations")
			}
		})
	}
}

func TestGRPCServer_GetUser(t *testing.T) {
	testUser := buildTestUser("test")

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
		wantErr bool
	}{
		{
			name:    "gets user using email",
			fields:  fields{testDAO},
			args:    args{in0: context.Background(), request: &proto.GetUserRequest{Identifier: &proto.GetUserRequest_Email{Email: testUser.User.Email}}},
			wantErr: false,
		},
		{
			name:    "gets user using uuid",
			fields:  fields{testDAO},
			args:    args{in0: context.Background(), request: &proto.GetUserRequest{Identifier: &proto.GetUserRequest_Uuid{Uuid: testUser.User.Uuid}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.GetUser(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}

			if got.User.Email != testUser.User.Email {
				t.Error("didn't fetch user")
			}
		})
	}
}

func TestGRPCServer_CreatePasswordResetToken(t *testing.T) {
	testUser := buildTestUser("test")

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
		wantErr bool
	}{
		{
			name:    "creates reset token",
			fields:  fields{testDAO},
			args:    args{in0: context.Background(), request: &proto.CreatePasswordResetTokenRequest{Email: testUser.User.Email, Expiration: timestampOneHour()}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreatePasswordResetToken(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			if got.PasswordResetToken == "" {
				t.Error("didn't build password reset token")
			}
		})
	}
}

func TestGRPCServer_UpdateUserPassword(t *testing.T) {
	testUser := buildTestUser("test")
	resetTokenResp, _ := testGRPCServer.CreatePasswordResetToken(context.Background(), &proto.CreatePasswordResetTokenRequest{Email: testUser.User.Email, Expiration: timestampOneHour()})
	resetToken := resetTokenResp.PasswordResetToken

	type fields struct {
		dao *db.DAO
	}
	type args struct {
		in0     context.Context
		request *proto.UpdateUserPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "updates users password",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.UpdateUserPasswordRequest{
				Email:                testUser.User.Email,
				PasswordResetToken:   resetToken,
				Password:             "test2",
				PasswordConfirmation: "test2",
			},
			},
			wantErr: false,
		},
		// TODO test confiming mismatched confirmation wont update
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.UpdateUserPassword(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			if got.Status != proto.UpdateUserPasswordResponse_SUCCESSFUL {
				t.Error("didn't successfully update password reset token")
			}

			//TODO Verify that the password has changed
		})
	}
}

func TestGRPCServer_CreateSession(t *testing.T) {
	testUser := buildTestUser("test")


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
		wantErr bool
	}{
		{
			name: "creates a session",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.CreateSessionRequest{Email:testUser.User.Email, Password:"test", ScopeGroupings:twoScopeGroupings()},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.CreateSession(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			if got.Session.Uuid == "" || got.Session.Token == "" || got.Session.Json == "" {
				t.Error("didn't create and return a session")
			}
		})
	}
}

func TestGRPCServer_GetSession(t *testing.T) {
	testUser := buildTestUser("test")

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
		wantErr bool
	}{
		{
			name: "gets a session",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.GetSessionRequest{Token:testUser.Session.Token},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			got, err := s.GetSession(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			if got.Session.Uuid == "" || got.Session.Token == "" || got.Session.Json == "" {
				t.Error("didn't create and return a session")
			}
		})
	}
}

func TestGRPCServer_DeleteSession(t *testing.T) {
	testUser := buildTestUser("test")

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
		wantErr bool
	}{
		{
			name: "deletes a session using a uuid",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.DeleteSessionRequest{Representation: &proto.DeleteSessionRequest_Uuid{Uuid:testUser.Session.Uuid}},
			},
			wantErr: false,
		},
		{
			name: "deletes a session using a token",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.DeleteSessionRequest{Representation: &proto.DeleteSessionRequest_Token{Token:testUser.Session.Token}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			_, err := s.DeleteSession(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			// todo validate session actually deleted
		})
	}
}

func TestGRPCServer_DeleteUser(t *testing.T) {
	testUser := buildTestUser("test")

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
		{
			name: "deletes a user",
			fields: fields{testDAO},
			args: args{in0: context.Background(), request: &proto.DeleteUserRequest{Email:testUser.User.Email, Password:"test"},},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GRPCServer{
				dao: tt.fields.dao,
			}
			_, err := s.DeleteUser(tt.args.in0, tt.args.request)
			if err != nil && tt.wantErr {
				return
			}
			// todo validate user actually deleted
		})
	}
}
