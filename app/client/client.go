package client

import (
	itempb "github.com/rnd/kudu/golang/protogen/item"
	"google.golang.org/grpc"
)

// kuduServiceClient holds interface of all kudu services.
type kuduServiceClient interface {
	itempb.ItemServiceClient
}

// New creates new kudu service client.
func New() (kuduServiceClient, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return struct {
		itempb.ItemServiceClient
	}{
		ItemServiceClient: itempb.NewItemServiceClient(conn),
	}, nil
}
