package server

import (
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)


func NewServer() {
	dao := db.ConnectToDatabase()
	defer dao.Conn.Close()

	repo := &Repo{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	server := &Server{dao:dao, repo:repo}

	proto.RegisterFingerprintServiceServer(s, server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}