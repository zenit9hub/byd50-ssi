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

// Package main implements a server for Greeter service.
package main

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/did/pkg/database"
	"byd50-ssi/pkg/did/registry"
	pb "byd50-ssi/proto-files"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	schemeMethod = "byd50"
)

var registryStore registry.Store

// server is used to implement proto-files.GreeterServer.
type server struct {
	pb.UnimplementedRegistryServer
}

// ScCreateDID implements proto-files.GreeterServer
func (s *server) ScCreateDID(ctx context.Context, in *pb.ScCreateDIDsRequest) (*pb.ScCreateDIDsReply, error) {
	pbKey := in.GetPublicKey()
	method := "byd50"
	createdDID, doc := dids.CreateDID(method, pbKey)

	if err := registryStore.Put(ctx, createdDID, doc); err != nil {
		log.Printf("[ScCreateDID] - [%v] store error: %v", createdDID, err)
	}

	log.Printf("[ScCreateDID] - [%v] %s", createdDID, doc)
	return &pb.ScCreateDIDsReply{Did: createdDID}, nil
}

// ScResolveDID implements proto-files.GreeterServer
func (s *server) ScResolveDID(ctx context.Context, in *pb.ScResolveDIDsRequest) (*pb.ScResolveDIDsReply, error) {
	// resolve DID's Document
	var resolutionError string
	var didDocument string
	var didDocumentMetadata string
	docuByteArray, err := registryStore.Get(ctx, in.GetDid())

	if err != nil {
		resolutionError = dids.NotFound.String()
		log.Printf("ScResolveDID error:%v", resolutionError)
	} else {
		didDocument = string(docuByteArray)
		didDocumentMetadata = ""
	}
	log.Printf("ScResolveDID - [%v] %v", in.GetDid(), string(docuByteArray))
	return &pb.ScResolveDIDsReply{ResolutionError: resolutionError, DidDocument: didDocument, DidDocumentMetadata: didDocumentMetadata}, nil
}

// ScUpdateDID implements proto-files.GreeterServer
func (s *server) ScUpdateDID(ctx context.Context, in *pb.ScUpdateDIDsRequest) (*pb.ScUpdateDIDsReply, error) {
	// update DID's Document
	result := "success"

	//validation check
	//if an error -> result = "Invalid document"

	ret, err := registryStore.Has(ctx, in.GetDid())
	if ret {
		if err := registryStore.Put(ctx, in.GetDid(), []byte(in.GetDocument())); err != nil {
			log.Printf("error caused by.. err[%v], ret[%v]", err, ret)
		}
		log.Printf("ScUpdateDID(%v) - [%v] %v", result, in.GetDid(), in.GetDocument())
	} else {
		log.Printf("error caused by.. err[%v], ret[%v]", err, ret)
	}

	return &pb.ScUpdateDIDsReply{Result: result}, nil
}

func initRegistry() {
	db, _ := database.Initialize()
	store, err := registry.NewLevelDBStore(db)
	if err != nil {
		log.Fatalf("failed to init registry store: %v", err)
	}
	registryStore = store
}

func main() {
	initRegistry()
	if store, ok := registryStore.(*registry.LevelDBStore); ok {
		defer store.Close()
	}

	lis, err := net.Listen("tcp", configs.UseConfig.DidRegistryPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterRegistryServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
