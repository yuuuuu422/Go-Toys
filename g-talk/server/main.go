package main

import (
	"context"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	pb "g-talk/proto"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	password string
	center   *Center
}

var (
	password string
	addr     string
)

func init() {
	flag.StringVar(&password, "password", "", "clent password")
	flag.StringVar(&addr, "addr", "", "server address")
	flag.Parse()
}

func newServer(password string) *Server {
	return &Server{
		password: password,
		center:   NewCenter(),
	}
}

func creatToken() token {
	b := make([]byte, 4)
	rand.Read(b)
	return token(fmt.Sprintf("%x", b))
}

func (s *Server) getToken(stream pb.Chat_StreamServer) (token, error) {
	md, _ := metadata.FromIncomingContext(stream.Context())
	tkn, ok := md["go-token"]
	if !ok {
		return "", errors.New("token not found")
	}
	return token(tkn[0]), nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Password != s.password {
		return nil, errors.New("invalid password")
	}
	log.Infof("[%s] is logged in", req.Username)

	tkn := creatToken()
	s.center.Register(tkn, req.Username)

	return &pb.LoginResponse{
		Token: string(tkn),
	}, nil
}

func (s *Server) Stream(stream pb.Chat_StreamServer) error {
	tkn, err := s.getToken(stream)
	if err != nil {
		return err
	}
	user, err := s.center.Get(tkn)
	if err != nil {
		return err
	}
	user.Stream = stream
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				s.center.LoginOut(tkn)
				return
			}
			resp := &pb.StreamResponse{
				Username: req.Username,
				Message:  req.Message,
			}
			s.center.Broadcast(resp)
		}
	}()
	select {
	case <-user.Done:
	case <-stream.Context().Done():
	}
	log.Infof("[%s] is logged out", user.Name)
	return nil
}
func main() {
	if flag.NFlag() != 2 {
		flag.Usage()
		return
	}

	addr := fmt.Sprintf(":%s", addr)
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Started listening on %s\n", addr)

	server := grpc.NewServer()

	pb.RegisterChatServer(server, newServer(password))
	if err := server.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
