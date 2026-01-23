//go:build integration
// +build integration

package main_test

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/kms"
	pb "byd50-ssi/proto-files"
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestRegistrarSmoke(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, configs.UseConfig.DidRegistrarAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("failed to connect registrar: %v", err)
	}
	defer conn.Close()

	client := pb.NewRegistrarClient(conn)
	sessionKMS, err := kms.InitKMS(kms.KeyTypeRSA)
	if err != nil {
		t.Fatalf("failed to init KMS: %v", err)
	}

	createResp, err := client.CreateDid(ctx, &pb.CreateDidRequest{
		PublicKeyBase58: sessionKMS.PbKeyBase58(),
		Method:          "byd50",
	})
	if err != nil {
		t.Fatalf("create did failed: %v", err)
	}
	if createResp.GetDid() == "" {
		t.Fatalf("create did returned empty DID")
	}

	resolveResp, err := client.ResolveDid(ctx, &pb.ResolveDidRequest{Did: createResp.GetDid()})
	if err != nil {
		t.Fatalf("resolve did failed: %v", err)
	}
	if resolveResp.GetResolutionError() != "" {
		t.Fatalf("resolve did error: %v", resolveResp.GetResolutionError())
	}
	if resolveResp.GetDidDocument() == "" {
		t.Fatalf("resolve did returned empty document")
	}
}
