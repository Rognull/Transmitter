package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math"
	pb "transmitter/api/gen/proto"
	db "transmitter/internal/db"
)

type Database interface {
	ConnectDB(dsn string) error
	Create(msg db.Message) error
}

func main() {
	var database Database
	database = &db.Database{}
	err := database.ConnectDB("host=localhost user=postgres dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	if err != nil {
		log.Fatal("Error while connect")
	}

	coeff := flag.Float64("f", 1.5, "coefficient")
	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewTransmitterServiceClient(conn)
	stream, err := c.StreamEntries(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	arr := make([]float64, 0, 1000)
	for i := 0; i < 500; i++ {
		response, _ := stream.Recv()
		arr = append(arr, response.Frequency)
	}

	var totalCounter int
	var anomalyCounter int
	targetSTD := SD(arr)
	targetMean := Mean(arr)
	for i := 0; i < 100000; i++ {
		response, err := stream.Recv()
		totalCounter++
		if err != nil {
			log.Fatal(err)
		}
		if response.Frequency > *coeff*targetSTD+targetMean || response.Frequency < targetMean-*coeff*targetSTD {
			r := db.Message{
				Session_id: response.SessionId,
				Frequency:  response.Frequency,
				Timestamp:  response.Timestamp.AsTime(),
			}
			err := database.Create(r)
			if err != nil {
				log.Fatal(err)
			}
			anomalyCounter++
			log.Println(response)
		}
	}
	log.Println(totalCounter, anomalyCounter)
}

func Mean(s []float64) float64 {
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}

func SD(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	mean := Mean(s)
	var sum float64
	for _, v := range s {
		sum += (v - mean) * (v - mean)
	}
	return math.Sqrt(sum / float64(len(s)))
}
