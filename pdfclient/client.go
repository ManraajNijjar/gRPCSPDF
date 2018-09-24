package main

import (
	"context"
	"fmt"
	"grpcCourse/pdf/pdfpb"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := pdfpb.NewPDFServiceClient(conn)

	sendFileStream(c)
}

func sendFileStream(c pdfpb.PDFServiceClient) {

	f, err := os.Open("/testpdf/test")

	stream, err := c.SendFile(context.Background())
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		dat, err := ioutil.ReadFile("/testpdf/test")
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}
		f, err := os.Open("/testpdf/test")
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		requests := []*pdfpb.PDFSendRequest{
			&pdfpb.PDFSendRequest{
				Content: make([]byte, 2),
			},
		}

		for _, req := range requests {
			fmt.Println("Sending data")
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error closing stream: %v", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
			}
			if err != nil {
				log.Fatalf("Error recieving: %v", err)
				close(waitc)
			}
			fmt.Printf("Received: %v", res.GetMessage())
		}
	}()

	<-waitc
}
