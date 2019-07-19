package grpcclient

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/digitalbitbox/bitbox-wallet-app/backend/bitboxbase/grpc/messages"
)

//go:generate protoc -I messages/ messages/bbb.proto --go_out=plugins=grpc:messages

type Client struct {
	address    string
	conn       *grpc.ClientConn
	baseClient pb.BitBoxBaseClient
}

func NewClient(address string) (*Client, error) {
	client := &Client{
		address: address,
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could no dial: %v", err)
		return client, err
	}
	client.conn = conn
	//defer client.conn.Close()
	baseClient := pb.NewBitBoxBaseClient(client.conn)
	client.baseClient = baseClient
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return client, err
}

func (client *Client) Connect() {
	connectRequest := "clientHello"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.baseClient.SayHello(ctx, &pb.HelloRequest{Name: connectRequest})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func (client *Client) GetEnv() *pb.BaseSystemEnvResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.baseClient.GetEnv(ctx, &pb.BaseSystemEnvRequest{})
	if err != nil {
		log.Fatalf("could not get base system environment %v", err)
	}
	return r
}
