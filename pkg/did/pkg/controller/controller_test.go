package controller

import (
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/did/kms"
	pb "byd50-ssi/proto-files"
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"testing"
)

type fakeRegistrarClient struct {
	docs map[string]string
}

func (f *fakeRegistrarClient) CreateDid(_ context.Context, in *pb.CreateDidRequest, _ ...grpc.CallOption) (*pb.CreateDidResponse, error) {
	method := in.GetMethod()
	if method == "" {
		method = "byd50"
	}
	createdDid, doc := dids.CreateDID(method, in.GetPublicKeyBase58())
	f.docs[createdDid] = string(doc)
	return &pb.CreateDidResponse{Did: createdDid}, nil
}

func (f *fakeRegistrarClient) RegisterDid(_ context.Context, _ *pb.RegisterDidRequest, _ ...grpc.CallOption) (*pb.RegisterDidResponse, error) {
	return &pb.RegisterDidResponse{}, nil
}

func (f *fakeRegistrarClient) ResolveDid(_ context.Context, in *pb.ResolveDidRequest, _ ...grpc.CallOption) (*pb.ResolveDidResponse, error) {
	doc, ok := f.docs[in.GetDid()]
	if !ok {
		return &pb.ResolveDidResponse{ResolutionError: "not_found"}, nil
	}
	return &pb.ResolveDidResponse{DidDocument: doc}, nil
}

func (f *fakeRegistrarClient) UpdateDid(_ context.Context, _ *pb.UpdateDidRequest, _ ...grpc.CallOption) (*pb.UpdateDidResponse, error) {
	return &pb.UpdateDidResponse{Result: "ok"}, nil
}

func TestControllerFlow(t *testing.T) {
	oldProvider := registrarClientProvider
	defer func() { registrarClientProvider = oldProvider }()

	fake := &fakeRegistrarClient{docs: map[string]string{}}
	registrarClientProvider = func() pb.RegistrarClient { return fake }

	dkms, err := kms.InitKMS(kms.KeyTypeRSA)
	if err != nil {
		t.Fatal(err)
	}

	did, err := CreateDIDWithErr(dkms.PbKeyBase58(), "byd50")
	if err != nil {
		t.Fatal(err)
	}
	if did == "" {
		t.Fatal(errors.New("did is empty"))
	}

	doc, err := ResolveDIDWithErr(did)
	if err != nil {
		t.Fatal(err)
	}
	if doc == "" {
		t.Fatal(errors.New("did document is empty"))
	}

	pbKey, err := GetPublicKeyWithErr(did, "")
	if err != nil {
		t.Fatal(err)
	}
	if pbKey != dkms.PbKeyBase58() {
		t.Fatal(errors.New("public key mismatch"))
	}

	plainText := "challenge-data"
	challenge := GetAuthChallengeString(did, plainText)
	response := GetAuthResponseString(challenge, dkms.PvKeyBase58())
	if response != plainText {
		t.Fatal(errors.New("challenge response mismatch"))
	}

	simplePresent := GetSimplePresent(did, dkms.PvKeyBase58())
	if VerifySimplePresent(simplePresent) != "success" {
		t.Fatal(errors.New("simple present failed"))
	}
}

func TestGetPublicKeyEmpty(t *testing.T) {
	oldProvider := registrarClientProvider
	defer func() { registrarClientProvider = oldProvider }()

	did := "did:byd50:empty"
	doc := dids.DocumentInterface{
		Context: []string{"https://www.w3.org/ns/dids/v1"},
		ID:      did,
		Authentication: []dids.AuthenticationProperty{
			{
				ID:              did + "#keys-1",
				Controller:      did,
				PublicKeyBase58: "",
			},
		},
	}
	docBytes, err := json.MarshalIndent(doc, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	fake := &fakeRegistrarClient{docs: map[string]string{did: string(docBytes)}}
	registrarClientProvider = func() pb.RegistrarClient { return fake }

	if pbKey := GetPublicKey(did, ""); pbKey != "" {
		t.Fatalf("expected empty public key, got %s", pbKey)
	}
	if _, err := GetPublicKeyWithErr(did, ""); err == nil {
		t.Fatal("expected GetPublicKeyWithErr error")
	}
}

func TestControllerErrorPaths(t *testing.T) {
	oldProvider := registrarClientProvider
	defer func() { registrarClientProvider = oldProvider }()

	fake := &fakeRegistrarClient{docs: map[string]string{}}
	registrarClientProvider = func() pb.RegistrarClient { return fake }

	if _, err := CreateDIDWithErr("", "byd50"); err == nil {
		t.Fatal("expected CreateDIDWithErr error")
	}
	if _, err := ResolveDIDWithErr(""); err == nil {
		t.Fatal("expected ResolveDIDWithErr error")
	}
	if _, err := GetPublicKeyWithErr("did:byd50:missing", ""); err == nil {
		t.Fatal("expected GetPublicKeyWithErr error")
	}
}
