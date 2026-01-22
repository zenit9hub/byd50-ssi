package driver

import (
	"byd50-ssi/pkg/did/configs"
	pb "byd50-ssi/proto-files"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

var (
	// ErrResolveDidByd50 - Sadly this is missing from method/byd50 registry
	ErrResolveDidByd50 = errors.New("method/byd50: resolve did error")
)

// DidMethodBYD50 - Implements the BYD50 did methods
type DidMethodBYD50 struct {
	Name string
}

// Specific instances for EC256 and company
var (
	didMethodBYD50 *DidMethodBYD50
)

func init() {
	// Register byd50 driver for BYD50 Chain Smart Contract
	didMethodBYD50 = &DidMethodBYD50{"byd50"}
	RegisterDidMethod(didMethodBYD50.Method(), func() DidMethod {
		return didMethodBYD50
	})
}

func (m *DidMethodBYD50) Method() string {
	return m.Name
}

// ResolveDid - Implements the ResolveDid method from DidMethod
// For this resolve did method, did must be an string
func (m *DidMethodBYD50) ResolveDid(did string) (string, string, string, error) {
	// Set up a connection to the server.
	registryClient := GetRegistryClient(configs.UseConfig.DidRegistryAddress)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := registryClient.ScResolveDID(ctx, &pb.ScResolveDIDsRequest{Did: did})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Received: %s", r.GetDidDocument())

	didDocument := r.GetDidDocument()
	ResolutionError := r.GetResolutionError()
	DidDocumentMetadata := r.GetDidDocumentMetadata()

	return didDocument, DidDocumentMetadata, ResolutionError, err
}

// CreateDid - Implements the RegisterDid method from DidMethod
// For this register did method, pbKeyBase58 must be an base58 encoded string
func (m *DidMethodBYD50) CreateDid(pbKeyBase58 string) (string, error) {
	// Set up a connection to the server.
	registryClient := GetRegistryClient(configs.UseConfig.DidRegistryAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := registryClient.ScCreateDID(ctx, &pb.ScCreateDIDsRequest{PublicKey: pbKeyBase58})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	createdDid := r.GetDid()

	return createdDid, nil
}

var (
	once sync.Once
	cli  pb.RegistryClient
)

func GetRegistryClient(serviceHost string) pb.RegistryClient {
	once.Do(func() {
		// Set up a connection to the server.
		conn, err := grpc.Dial(
			serviceHost,
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		cli = pb.NewRegistryClient(conn)
	})
	return cli
}
