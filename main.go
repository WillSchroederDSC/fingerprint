package main

import (
	"context"
	"log"
	"net"

	pb "github.com/willschroeder/fingerprint/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type Server struct{}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterFingerprintServiceServer(s, &Server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) CreateUser(_ context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &pb.User{Uuid: "1111", Email: request.Email}
	return &pb.CreateUserResponse{User: user}, nil
}
