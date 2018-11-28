package server

import (
	"github.com/willschroeder/fingerprint/pkg/db"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

const (
	port = ":50051"
)


func NewServer() {
	dao := db.ConnectToDatabase()
	defer dao.Conn.Close()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	repo := &Repo{dao:dao}
	server := &GRPCSerer{
		dao: dao,
		repo:repo,
		logger:logger,
	}

	// GRPC Setup, taken from google's Hello World example
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterFingerprintServiceServer(s, server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}