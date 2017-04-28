package client

import (
	"log"
	"os"

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
		UserServiceClient: userpb.NewUserServiceClient(mustDial(os.Getenv("USER_SVC_SERVICE_HOST") + ":" + os.Getenv("USER_SVC_SERVICE_PORT"))),
		TaskServiceClient: taskpb.NewTaskServiceClient(mustDial(os.Getenv("TASK_SVC_SERVICE_HOST") + ":" + os.Getenv("TASK_SVC_SERVICE_PORT"))),
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
