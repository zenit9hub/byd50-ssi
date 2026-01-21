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
	"byd50-ssi/did/configs"
	"byd50-ssi/did/core"
	byd50_jwt "byd50-ssi/did/core/byd50-jwt"
	"byd50-ssi/did/pkg/controller"
	"byd50-ssi/did/pkg/logger"
	pb "byd50-ssi/proto-files"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"log"
	"reflect"
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

func UseCase1DefaultAuthentication(dkms core.DKMS) {
	logger.FuncStart()

	// Set up a connection to the server.
	relyingPartyClient := GetRelyingPartyClient(configs.UseConfig.RelyingPartyAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/* Use Case 1. Default Authentication
	1. Request AuthChallenge
	2. Recv 'Auth Challenge String'
	3. decrypt AuthChallenge String
	*/
	log.Printf("myDKMS.DID = " + dkms.Did())
	challengeReply, err := relyingPartyClient.AuthChallenge(ctx, &pb.ChallengeRequest{Did: dkms.Did()})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	challengeString := challengeReply.GetAuthChallenge()
	log.Printf("Received Challenge String(%v)", len(challengeString))

	authResponseString := controller.GetAuthResponseString(challengeString, dkms.PvKeyBase58())
	log.Printf("decrypted string: %s", string(authResponseString))

	/* Use Case 1. Default Authentication
	4. Send 'Auth Response String'
	5. Recv result
	*/
	responseReply, err2 := relyingPartyClient.AuthResponse(ctx, &pb.ResponseRequest{AuthResponse: string(authResponseString)})
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	log.Printf("Auth Response result: %s", responseReply.GetMessage())
	logger.FuncEnd()
}

func UseCase2SimpleAuthentication(dkms core.DKMS) {
	logger.FuncStart()

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
	log.Printf("myDKMS.DID = " + dkms.Did())

	simplePresentString := controller.GetSimplePresent(dkms.Did(), dkms.PvKeyBase58())

	simplePresentReply, err := relyingPartyClient.SimplePresent(ctx, &pb.SimplePresentRequest{SimplePresent: simplePresentString})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Auth Response result: %s", simplePresentReply.GetResult())
	logger.FuncEnd()
}

func UseCase3RequestCredential(dkms core.DKMS) {
	logger.FuncStart()

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
	log.Printf("myDKMS.DID = " + dkms.Did())

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
	pvKey := dkms.PvKey()

	var key *ecdsa.PrivateKey
	if reflect.TypeOf(pvKey) != reflect.TypeOf(key) {
		log.Fatalf(" error: Key type error. %v", reflect.TypeOf(pvKey))
	}
	vcRequestJwt := core.CreateVcWithClaims(kid, claims, pvKey.(*ecdsa.PrivateKey))
	log.Printf(" -- vc request jwt --\n%v", vcRequestJwt)

	// ******************** Request VC ******************** //
	credentialReply, err := issuerClient.RequestCredential(ctx, &pb.CredentialRequest{VcClaimJwt: vcRequestJwt})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("RequestCredential Response result: %s", credentialReply.GetVcJwt())

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
	holderPvKey := dkms.PvKey().(*ecdsa.PrivateKey)
	myVp := core.CreateVpWithClaims(holderDid, vpClaims, holderPvKey)
	log.Printf("myVp: %v", myVp)

	// ******************** Send VP ******************** //
	VpReply, err := relyingPartyClient.VerifyVp(ctxRp, &pb.VerifyVpRequest{Vp: myVp})

	log.Printf("result: %v", VpReply.GetResult())

	logger.FuncEnd()
}

func geth() {
	ks := keystore.NewKeyStore("/path/to/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)

	// Create a new account with the specified encryption passphrase.
	newAcc, _ := ks.NewAccount("Creation password")
	fmt.Println(newAcc)

	// Export the newly created account with a different passphrase. The returned
	// data from this method invocation is a JSON encoded, encrypted key-file.
	jsonAcc, _ := ks.Export(newAcc, "Creation password", "Export password")

	// Update the passphrase on the account created above inside the local keystore.
	_ = ks.Update(newAcc, "Creation password", "Update password")

	// Delete the account updated above from the local keystore.
	_ = ks.Delete(newAcc, "Update password")

	// Import back the account we've exported (and then deleted) above with yet
	// again a fresh passphrase.
	impAcc, _ := ks.Import(jsonAcc, "Export password", "Import password")

	//Signing from Go
	// Create a new account to sign transactions with
	signer, _ := ks.NewAccount("Signer password")
	txHash := common.HexToHash("0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	// Sign a transaction with a single authorization
	signature, _ := ks.SignHashWithPassphrase(signer, "Signer password", txHash.Bytes())

	// Sign a transaction with multiple manually cancelled authorizations
	_ = ks.Unlock(signer, "Signer password")
	signature, _ = ks.SignHash(signer, txHash.Bytes())
	_ = ks.Lock(signer.Address)

	// Sign a transaction with multiple automatically cancelled authorizations
	_ = ks.TimedUnlock(signer, "Signer password", time.Second)
	signature, _ = ks.SignHash(signer, txHash.Bytes())

	log.Printf("signature=%v, impAcc=%v, am=%v", signature, impAcc, am)
}

func main() {
	// Initialize DKMS
	myDKMS, err := core.InitDKMS(core.KeyTypeRSA)
	if err != nil {
		log.Fatalf("could not Init DKMS (%v)", err.Error())
	}

	// Create DID
	method := "byd50"
	did := controller.CreateDID(myDKMS.PbKeyBase58(), method)
	myDKMS.SetDid(did)

	// Use Case 1. Default Authentication
	UseCase1DefaultAuthentication(myDKMS)

	// Use Case 2. Simple Authentication
	UseCase2SimpleAuthentication(myDKMS)

	// Use Case 3. Simple Authentication
	// Initialize DKMS
	myDKMS, err = core.InitDKMS(core.KeyTypeECDSA)
	// Create DID
	did = controller.CreateDID(myDKMS.PbKeyBase58(), method)
	myDKMS.SetDid(did)
	UseCase3RequestCredential(myDKMS)
}
