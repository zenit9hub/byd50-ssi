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

// Package main implements a client for Greeter service.
package main

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/kms"
	"byd50-ssi/pkg/did/pkg/controller"
	pb "byd50-ssi/proto-files"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

var (
	onceRpRPC          sync.Once
	relyingPartyClient pb.RelyingPartyClient
)

var (
	onceIssRPC   sync.Once
	issuerClient pb.IssuerClient
)

func logSectionStart(title string) {
	log.Printf("\n[%s] START", title)
}

func logSectionEnd(title string) {
	log.Printf("[%s] END\n", title)
}

func logDidDocument(label, doc string) {
	if doc == "" {
		log.Printf("\n[%s]\n<empty>", label)
		return
	}
	pretty := doc
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(doc), &payload); err == nil {
		if buf, err := json.MarshalIndent(payload, "", "  "); err == nil {
			pretty = string(buf)
		}
	}
	log.Printf("\n[%s]\n%s", label, pretty)
}

func mustPvKeyECDSA(dkms kms.KMS) *ecdsa.PrivateKey {
	pvKey, err := dkms.PvKeyECDSA()
	if err != nil {
		log.Fatalf("invalid ECDSA private key: %v", err)
	}
	return pvKey
}

func GetRelyingPartyClient(serviceHost string) pb.RelyingPartyClient {
	onceRpRPC.Do(func() {
		// Set up a connection to the server.
		conn, err := grpc.Dial(
			serviceHost,
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		relyingPartyClient = pb.NewRelyingPartyClient(conn)
	})
	return relyingPartyClient
}

func GetIssuerClient(serviceHost string) pb.IssuerClient {
	onceIssRPC.Do(func() {
		// Set up a connection to the server.
		conn, err := grpc.Dial(
			serviceHost,
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		issuerClient = pb.NewIssuerClient(conn)
	})
	return issuerClient
}

func UseCase1DefaultAuthentication(dkms kms.KMS) {
	// Set up a connection to the server.
	relyingPartyClient := GetRelyingPartyClient(configs.UseConfig.RelyingPartyAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/* Use Case 1. Default Authentication
	1. Request AuthChallenge
	2. Recv 'Auth Challenge String'
	3. decrypt AuthChallenge String
	*/
	challengeReply, err := relyingPartyClient.AuthChallenge(ctx, &pb.ChallengeRequest{Did: dkms.Did()})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	challengeString := challengeReply.GetAuthChallenge()
	log.Printf("Challenge received (len=%v)", len(challengeString))

	authResponseString := dkms.Decrypt(challengeString)
	log.Printf("Challenge decrypted")

	/* Use Case 1. Default Authentication
	4. Send 'Auth Response String'
	5. Recv result
	*/
	responseReply, err2 := relyingPartyClient.AuthResponse(ctx, &pb.ResponseRequest{AuthResponse: string(authResponseString)})
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	log.Printf("Auth response: %s", responseReply.GetMessage())
}

func UseCase2SimpleAuthentication(dkms kms.KMS) {
	// Set up a connection to the server.
	relyingPartyClient := GetRelyingPartyClient(configs.UseConfig.RelyingPartyAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/* Use Case 2. Simple Authentication
	1. SimplePresentDID
	   didAndTime := "DID" + ";" + time.Now().UTC().String()
	   signedStr := RsaSign(didAndTime)
	   challengeString := didAndTime + ";" + signedStr
	*/
	simplePresentString := controller.GetSimplePresent(dkms.Did(), dkms.PvKeyBase58())

	simplePresentReply, err := relyingPartyClient.SimplePresent(ctx, &pb.SimplePresentRequest{SimplePresent: simplePresentString})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Simple present result: %s", simplePresentReply.GetResult())
}

func UseCase3RequestCredential(dkms kms.KMS) {
	// Set up a connection to the server.
	issuerClient := GetIssuerClient(configs.UseConfig.IssuerAddress)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Set up a connection to the server.
	relyingPartyClient := GetRelyingPartyClient(configs.UseConfig.RelyingPartyAddress)
	ctxRp, cancelRp := context.WithTimeout(context.Background(), time.Second)
	defer cancelRp()

	/* Use Case 3. Request Credential
	1. JWT
		credentialsubject
	*/
	log.Printf("Holder DID: %s", dkms.Did())

	// ******************** Build VC Claims ******************** //
	nonce := core.RandomString(12)

	typ := "AlumniCredential"
	typArray := []string{"VerifiableCredential"}
	typArray = append(typArray, typ)

	credSub := map[string]interface{}{
		"degree": "BachelorDegree",
		"name":   "<span lang='fr-CA'>Baccalauréat en musiques numériques</span>",
	}
	myVcClaims := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":              typArray,
		"credentialSubject": credSub,
	}

	standardClaims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
		IssuedAt:  time.Now().Unix(),
		Issuer:    "http://google.com/issuer",
		NotBefore: time.Now().Unix(),
		Subject:   "",
	}

	// Create the Claims
	claims := byd50_jwt.VcClaims{
		nonce,
		myVcClaims,
		standardClaims,
	}

	// ******************** Create VC with claims ******************** //
	kid := dkms.Did()
	pvKey := mustPvKeyECDSA(dkms)
	vcRequestJwt := core.CreateVcWithClaims(kid, claims, pvKey)
	log.Printf("\n[VC Request JWT]\n%v", vcRequestJwt)

	// ******************** Request VC ******************** //
	credentialReply, err := issuerClient.RequestCredential(ctx, &pb.CredentialRequest{VcClaimJwt: vcRequestJwt})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("\n[VC JWT]\n%v", credentialReply.GetVcJwt())

	myVc := credentialReply.GetVcJwt()

	// ******************** Make VP ******************** //
	// ******************** Build VP Claims ******************** //
	nonce = core.RandomString(12)

	typ = "CredentialManagerPresentation"
	typArray = []string{"VerifiablePresentation"}
	typArray = append(typArray, typ)

	var vcJwtArray []string
	vcJwtArray = append(vcJwtArray, myVc)

	myVpClaims := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":                 typArray,
		"verifiableCredential": vcJwtArray,
	}

	standardClaims = jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
		IssuedAt:  time.Now().Unix(),
		Issuer:    "client make this vp",
		NotBefore: time.Now().Unix(),
		Subject:   "",
	}

	// Create the Claims
	vpClaims := byd50_jwt.VpClaims{
		nonce,
		myVpClaims,
		standardClaims,
	}

	holderDid := dkms.Did()
	holderPvKey := mustPvKeyECDSA(dkms)
	myVp := core.CreateVpWithClaims(holderDid, vpClaims, holderPvKey)
	log.Printf("\n[VP JWT]\n%v", myVp)

	// ******************** Send VP ******************** //
	VpReply, err := relyingPartyClient.VerifyVp(ctxRp, &pb.VerifyVpRequest{Vp: myVp})

	log.Printf("VP verify result: %v", VpReply.GetResult())
	log.Printf("\n[VP Verify PublicKey PEM]\n%v", dkms.PbKeyPEM())
}

func main() {
	// Auth flows use RSA; VC/VP flows use ECDSA.
	authKMS, err := kms.InitKMS(kms.KeyTypeRSA)
	if err != nil {
		log.Fatalf("could not Init KMS (%v)", err.Error())
	}

	method := "byd50"
	authDid := controller.CreateDID(authKMS.PbKeyBase58(), method)
	authKMS.SetDid(authDid)
	log.Printf("\n[Auth DID]\n%s", authDid)
	authDidDoc := controller.ResolveDID(authDid)
	logDidDocument("DID Document (Auth)", authDidDoc)

	logSectionStart("DID Auth Challenge & Response")
	UseCase1DefaultAuthentication(authKMS)
	logSectionEnd("DID Auth Challenge & Response")

	logSectionStart("DID Simple Presentation")
	UseCase2SimpleAuthentication(authKMS)
	logSectionEnd("DID Simple Presentation")

	credKMS, err := kms.InitKMS(kms.KeyTypeECDSA)
	if err != nil {
		log.Fatalf("could not Init KMS (%v)", err.Error())
	}
	credDid := controller.CreateDID(credKMS.PbKeyBase58(), method)
	credKMS.SetDid(credDid)
	log.Printf("\n[Credential DID]\n%s", credDid)

	logSectionStart("VC Issue & VP Submit")
	UseCase3RequestCredential(credKMS)
	logSectionEnd("VC Issue & VP Submit")
}
