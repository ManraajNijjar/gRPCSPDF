package main

import (
	"grpcCourse/pdf/pdfpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) SendFile(stream pdfpb.PDFService_SendFileServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading: %v", err)
			return err
		}

		bytes := req.GetContent()

		sendErr := stream.Send(&pdfpb.PDFSendResponse{
			Message: "Its good",
			Code:    0,
		})

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pdfpb.RegisterPDFServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
