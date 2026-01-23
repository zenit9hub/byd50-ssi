/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Registrar/Resolver service.
package main

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/core"
	"byd50-ssi/pkg/did/core/driver"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "byd50-ssi/proto-files"
	"google.golang.org/grpc"
)

// server is used to implement proto-files.RegistrarServer.
type server struct {
	pb.UnimplementedRegistrarServer
}

// CreateDID implements proto-files.RegistrarServer
func (s *server) CreateDid(ctx context.Context, in *pb.CreateDidRequest) (*pb.CreateDidResponse, error) {
	log.Printf("[CreateDid] req ~> CreateDid")

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Get the method from the request and set default if nil
	methodStr := in.GetMethod()
	if methodStr == "" {
		methodStr = "byd50"
	}

	// Verify that public key length is greater than 1.
	pbKey := in.GetPublicKeyBase58()
	log.Printf("methodStr[{%v}]", methodStr)
	if len(os.Args) > 1 {
		pbKey = os.Args[1]
	}

	// Get a suitable driver with 'did method', call the CreateDid function.
	did, _ := driver.GetDidMethod(methodStr).CreateDid(pbKey)
	log.Printf("[CreateDID] reply <~ %v", did)

	return &pb.CreateDidResponse{Did: did}, nil
}

// ResolveDID implements proto-files.RegistrarServer
func (s *server) ResolveDid(ctx context.Context, in *pb.ResolveDidRequest) (*pb.ResolveDidResponse, error) {
	log.Printf("[ResolveDid] Received DID: %v", in.GetDid())

	// Contact the server and print out its response.
	dID := in.GetDid()

	// validation check
	slice := strings.Split(dID, ":")
	scheme := slice[0]
	didMethod := slice[1]

	if "did" != scheme {
		log.Printf("invalid DID. should be started with 'did'")
		// error handling..
	} else if !core.Contains(configs.UseConfig.AdoptedDriverList, didMethod) {
		log.Printf("invalid Method. driver wasn't adopted")
		log.Printf("Adopted Driver List [%v]", configs.UseConfig.AdoptedDriverList)
		// error handling..
	}

	// Get a suitable driver with 'did method', call the ResolveDid function.
	didDocument, didDocumentMetadata, resolutionError, _ := driver.GetDidMethod(didMethod).ResolveDid(dID)
	if didDocument == "" {
		// error handling..
	}

	return &pb.ResolveDidResponse{DidDocument: didDocument, DidDocumentMetadata: didDocumentMetadata, ResolutionError: resolutionError}, nil
}

// UpdateDID implements proto-files.RegistrarServer
func (s *server) UpdateDid(ctx context.Context, in *pb.UpdateDidRequest) (*pb.UpdateDidResponse, error) {
	log.Printf("[UpdateDid] Received DID: %v", in.GetDid())

	return &pb.UpdateDidResponse{Result: ""}, nil
}

func main() {
	lis, err := net.Listen("tcp", configs.UseConfig.DidRegistrarPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRegistrarServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
