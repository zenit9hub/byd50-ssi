package controller

import (
	"byd50-ssi/did/core"
	"byd50-ssi/did/core/dids"
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

func (f *fakeRegistrarClient) CreateDID(_ context.Context, in *pb.CreateDIDsRequest, _ ...grpc.CallOption) (*pb.CreateDIDsReply, error) {
	method := in.GetMethod()
	if method == "" {
		method = "byd50"
	}
	createdDid, doc := dids.CreateDID(method, in.GetPublicKeyBase58())
	f.docs[createdDid] = string(doc)
	return &pb.CreateDIDsReply{Did: createdDid}, nil
}

func (f *fakeRegistrarClient) RegisterDID(_ context.Context, _ *pb.RegisterDIDsRequest, _ ...grpc.CallOption) (*pb.RegisterDIDsReply, error) {
	return &pb.RegisterDIDsReply{}, nil
}

func (f *fakeRegistrarClient) ResolveDID(_ context.Context, in *pb.ResolveDIDsRequest, _ ...grpc.CallOption) (*pb.ResolveDIDsReply, error) {
	doc, ok := f.docs[in.GetDid()]
	if !ok {
		return &pb.ResolveDIDsReply{ResolutionError: "not_found"}, nil
	}
	return &pb.ResolveDIDsReply{DidDocument: doc}, nil
}

func (f *fakeRegistrarClient) UpdateDID(_ context.Context, _ *pb.UpdateDIDsRequest, _ ...grpc.CallOption) (*pb.UpdateDIDsReply, error) {
	return &pb.UpdateDIDsReply{Result: "ok"}, nil
}

func TestControllerFlow(t *testing.T) {
	oldProvider := registrarClientProvider
	defer func() { registrarClientProvider = oldProvider }()

	fake := &fakeRegistrarClient{docs: map[string]string{}}
	registrarClientProvider = func() pb.RegistrarClient { return fake }

	dkms, err := core.InitDKMS(core.KeyTypeRSA)
	if err != nil {
		t.Fatal(err)
	}

	did := CreateDID(dkms.PbKeyBase58(), "byd50")
	if did == "" {
		t.Fatal(errors.New("did is empty"))
	}

	doc := ResolveDID(did)
	if doc == "" {
		t.Fatal(errors.New("did document is empty"))
	}

	pbKey := GetPublicKey(did, "")
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
}
