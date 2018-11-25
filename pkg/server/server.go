package server

import "github.com/willschroeder/fingerprint/pkg/db"

const (
	port = ":50051"
)


func NewServer() {
	//lis, err := net.Listen("tcp", port)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//s := grpc.NewServer()
	//
	//proto.RegisterFingerprintServiceServer(s, &Server{})
	//reflection.Register(s)
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}

	db := db.ConnectToDatabase()
	defer db.Conn.Close()

	GetUser(db)
}