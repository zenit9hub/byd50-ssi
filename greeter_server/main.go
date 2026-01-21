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
	"byd50-ssi/did/configs"
	"byd50-ssi/did/core"
	"byd50-ssi/did/pkg/controller"
	pb "byd50-ssi/proto-files"
	"context"
	"github.com/btcsuite/btcutil/base58"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var sourceData = "randomStr;2021-06-08T14:04:43UTC"

// server is used to implement proto-files.GreeterServer.
type server struct {
	pb.UnimplementedRelyingPartyServer
}

// AuthChallenge implements proto-files.GreeterServer
func (s *server) AuthChallenge(_ context.Context, in *pb.ChallengeRequest) (*pb.ChallengeReply, error) {
	log.Printf("[AuthChallenge][Request] DID: %v", in.GetDid())
	sourceData = time.Now().UTC().String()
	sourceData = base58.Encode([]byte(sourceData)) + ";" + sourceData
	authChallengeString := controller.GetAuthChallengeString(in.GetDid(), sourceData)
	log.Printf("[AuthChallenge][Reply] authChallengeString(%v)", len(authChallengeString))

	return &pb.ChallengeReply{AuthChallenge: authChallengeString}, nil
}

// AuthResponse implements proto-files.GreeterServer
func (s *server) AuthResponse(_ context.Context, in *pb.ResponseRequest) (*pb.ResponseReply, error) {
	log.Printf("[AuthResponse][Request] %v", in.AuthResponse)
	message := "error"

	// compare Challenge and Response string
	if in.GetAuthResponse() == sourceData {
		log.Printf("[AuthResponse] Compare success")
		message = "success"
	}

	log.Printf("[AuthResponse][Reply] : " + message)
	return &pb.ResponseReply{Message: message}, nil
}

// SimplePresent implements proto-files.GreeterServer
func (s *server) SimplePresent(_ context.Context, in *pb.SimplePresentRequest) (*pb.SimplePresentReply, error) {
	log.Printf("[SimplePresent][Request]")
	result := controller.VerifySimplePresent(in.GetSimplePresent())
	log.Printf("[SimplePresent][Reply] result: %v", result)

	return &pb.SimplePresentReply{Result: result}, nil
}

// SimplePresent implements proto-files.GreeterServer
func (s *server) VerifyVp(_ context.Context, in *pb.VerifyVpRequest) (*pb.VerifyVpReply, error) {
	log.Printf("[VerifyVp][Request]")
	valid, _, err := core.VerifyVp(in.GetVp(), controller.GetPublicKey)
	log.Printf("[VerifyVp][Reply] valid: %v err: %v", valid, err)

	result := ""
	if valid {
		result = "verified"
	}
	return &pb.VerifyVpReply{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", configs.UseConfig.RelyingPartyPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRelyingPartyServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
