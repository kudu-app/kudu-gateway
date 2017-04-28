package client

import (
	"log"

	taskpb "github.com/rnd/kudu/golang/protogen/task"
	userpb "github.com/rnd/kudu/golang/protogen/user"
	"google.golang.org/grpc"
)

// KuduServiceClient holds interface of all kudu services.
type KuduServiceClient interface {
	userpb.UserServiceClient
	taskpb.TaskServiceClient
}

// New creates new kudu service client.
func New() (KuduServiceClient, error) {
	return struct {
		userpb.UserServiceClient
		taskpb.TaskServiceClient
	}{
		UserServiceClient: userpb.NewUserServiceClient(mustDial("localhost:9111")),
		TaskServiceClient: taskpb.NewTaskServiceClient(mustDial("localhost:9112")),
	}, nil
}

// mustDial ensures a tcp connection to specified address.
func mustDial(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return conn
}
