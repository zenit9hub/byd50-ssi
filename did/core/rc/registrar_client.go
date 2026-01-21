package rc

import (
	"byd50-ssi/did/configs"
	pb "byd50-ssi/proto-files"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var (
	onceRC          sync.Once
	registrarClient pb.RegistrarClient
)

func GetRegistrarClient() pb.RegistrarClient {
	serviceHost := configs.UseConfig.DidRegistrarAddress
	onceRC.Do(func() {
		// Set up a connection to the server.
		conn, err := grpc.Dial(
			serviceHost,
			grpc.WithInsecure(),
			//grpc_middleware.ChainUnaryClient(), //ToDo.. reconnection logic, logging, etc..
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		registrarClient = pb.NewRegistrarClient(conn)
	})
	return registrarClient
}
