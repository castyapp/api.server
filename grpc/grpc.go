package grpc

import (
	"fmt"

	"github.com/castyapp/api.server/config"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"google.golang.org/grpc"
)

var (
	UserServiceClient     proto.UserServiceClient
	AuthServiceClient     proto.AuthServiceClient
	TheaterServiceClient  proto.TheaterServiceClient
	MessagesServiceClient proto.MessagesServiceClient
)

func Configure() error {

	var (
		host = config.Map.Grpc.Host
		port = config.Map.Grpc.Port
	)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not dialing grpc.server: %v", err)
	}

	UserServiceClient = proto.NewUserServiceClient(conn)
	AuthServiceClient = proto.NewAuthServiceClient(conn)
	TheaterServiceClient = proto.NewTheaterServiceClient(conn)
	MessagesServiceClient = proto.NewMessagesServiceClient(conn)
	return nil
}
