package transport

import (
	"net"

	userpb "github.com/eotet/project-protos/gen/go/user"
	"github.com/eotet/users-service/internal/user"
	"google.golang.org/grpc"
)

func RunGRPC(svc *user.Service) error {
	listener, err := net.Listen("tcp", ":5501")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, NewHandler(svc))
	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
