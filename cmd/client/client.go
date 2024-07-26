package main

import (
	"context"
	"io"
	"log"
	"os"

	grpcservice "github.com/adetunjii/streaming-server/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := grpcservice.NewVideoServiceClient(conn)

	stream, err := client.GetVideoStream(context.Background(), &grpcservice.VideoRequest{Id: "hero-bg.png"})
	if err != nil {
		log.Fatalf("could not stream video: %v", err)
	}

	outFile, err := os.Create("received.png")
	if err != nil {
		log.Fatalf("could not create output file: %v", err)
	}
	defer outFile.Close()

	total := 0

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving chunk: %v", err)
		}
		outFile.Write(chunk.Buffer)

		total += len(chunk.Buffer)
	}

	log.Println("Video received and saved as received.png")
	log.Printf("Received %d bytes of data", total)
}
