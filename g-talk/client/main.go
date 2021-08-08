package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	pb "g-talk/proto"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	username string
	password string
	addr     string
)

type Client struct {
	grpcClient pb.ChatClient
}

func init() {
	flag.StringVar(&username, "name", "", "client username")
	flag.StringVar(&password, "password", "", "clent password")
	flag.StringVar(&addr, "addr", "", "server address")
	flag.Parse()
}

func (c *Client) Login(addr string, username string, password string) (string, error) {
	req := &pb.LoginRequest{
		Username: username,
		Password: password,
	}
	resp, err := c.grpcClient.Login(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil

}

func (c *Client) Stream(username string, token string) error {
	md := metadata.Pairs("go-token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.grpcClient.Stream(ctx)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for {
			if !sc.Scan() {
				log.Fatalln(sc.Err())
			}
			req := &pb.StreamRequest{
				Token:    token,
				Username: username,
				Message:  sc.Text(),
			}
			if err := stream.Send(req); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("[%s]: %s\n", resp.Username, resp.Message)
	}
}

func main() {
	if flag.NFlag() != 3 {
		flag.Usage()
		return
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := &Client{
		grpcClient: pb.NewChatClient(conn),
	}
	token, err := c.Login(addr, username, password)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Connect to server: %s",addr)
	if err := c.Stream(username, token); err != nil {
		log.Fatalln(err)
	}
}
