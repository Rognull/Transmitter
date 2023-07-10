package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
	pb "transmitter/api/gen/proto"
)

type Server struct {
	pb.UnimplementedTransmitterServiceServer
}

func (s *Server) StreamEntries(empty *emptypb.Empty, stream pb.TransmitterService_StreamEntriesServer) error {
	// Generate random values for mean and standard deviation
	mean := rand.Float64()*20 - 10
	std := rand.Float64()*1.2 + 0.3

	// Generate a random UUID for session_id
	sessionID := uuid.New().String()

	log.Printf("Session ID: %s, Mean: %f, STD: %f", sessionID, mean, std)

	for {
		// Create and send entry with random values
		//time.Sleep(time.Millisecond * 100)
		entry := &pb.Entry{
			SessionId: sessionID,
			Frequency: rand.NormFloat64()*std + mean,
			Timestamp: timestamppb.Now(),
		}

		if err := stream.Send(entry); err != nil {
			return status.Errorf(codes.Unknown, "error sending entry: %v", err)
		}

		// Log the generated values

	}
}

func StartServer(ctx context.Context, s *grpc.Server) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pb.RegisterTransmitterServiceServer(s, &Server{})
	log.Println("Server started")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
