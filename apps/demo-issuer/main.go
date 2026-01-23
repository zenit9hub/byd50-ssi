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
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/kms"
	"byd50-ssi/pkg/did/pkg/controller"
	pb "byd50-ssi/proto-files"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var sourceData = "randomStr;2021-06-08T14:04:43UTC"
var issuerDid string
var myDkms kms.KMS

// server is used to implement proto-files.GreeterServer.
type server struct {
	pb.UnimplementedIssuerServer
}

func mustPvKeyECDSA(dkms kms.KMS) *ecdsa.PrivateKey {
	pvKey, err := dkms.PvKeyECDSA()
	if err != nil {
		log.Fatalf("invalid ECDSA private key: %v", err)
	}
	return pvKey
}

// RequestCredential implements proto-files.GreeterServer
func (s *server) RequestCredential(_ context.Context, in *pb.CredentialRequest) (*pb.CredentialReply, error) {
	log.Printf("[RequestCredential][Request]")
	subjectDid := ""
	parseToken, err := jwt.Parse(in.VcClaimJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		subjectDid = token.Header["kid"].(string)
		pbKeyBase58 := controller.GetPublicKey(subjectDid, "")
		pbKey, _ := x509.ParsePKIXPublicKey(base58.Decode(pbKeyBase58))
		return pbKey, nil
	})
	if err != nil {
		log.Fatalf("error: Request Credential has an error: %v", err)
	}

	vcJwt := ""
	if parseToken.Valid {
		kid := issuerDid
		typ := "AlumniCredential"
		pvKey := mustPvKeyECDSA(myDkms)
		log.Printf("pvkey : %v", pvKey)
		standardClaims := jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://demo-issuer.com/issuer142857",
			NotBefore: time.Now().Unix(),
			Subject:   subjectDid,
		}

		if claims, ok := parseToken.Claims.(jwt.MapClaims); ok {
			vc := claims["vc"].(map[string]interface{})
			credSub := vc["credentialSubject"].(map[string]interface{})
			vcJwt = core.CreateVc(kid, typ, credSub, standardClaims, pvKey)
		}

	}
	log.Printf("[RequestCredential][Reply] vcJwt: %v", vcJwt)

	return &pb.CredentialReply{VcJwt: vcJwt}, nil
}

// ReqCredIdCard implements proto-files.GreeterServer
func (s *server) ReqCredIdCard(_ context.Context, in *pb.IdCardRequest) (*pb.IdCardReply, error) {
	log.Printf("[ReqCredIdCard][Request]")

	// ******************** Build VC Claims ******************** //
	subjectDid := in.GetDid()
	nonce := core.RandomString(12)

	typ := "eIdCardCredential"
	typArray := []string{"VerifiableCredential"}
	typArray = append(typArray, typ)

	credSub := map[string]interface{}{
		"country": "S.Korea",
		"name":    "Hong Gil-Dong",
		"birth":   "2000-11-08",
	}
	myVc := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":              typArray,
		"credentialSubject": credSub,
	}

	standardClaims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Minute * 3).Unix(),
		Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
		IssuedAt:  time.Now().Unix(),
		Issuer:    "http://www.gov.kr/residentregistration",
		NotBefore: time.Now().Unix(),
		Subject:   subjectDid,
	}

	// Create the Claims
	claims := byd50_jwt.VcClaims{
		nonce,
		myVc,
		standardClaims,
	}
	kid := issuerDid
	pvKey := mustPvKeyECDSA(myDkms)
	eIdVcJwt := core.CreateVcWithClaims(kid, claims, pvKey)

	log.Printf("[ReqCredIdCard][Reply] vcJwt: %v", eIdVcJwt)
	return &pb.IdCardReply{EidVcJwt: eIdVcJwt}, nil
}

// ReqCredDlCard implements proto-files.GreeterServer
func (s *server) ReqCredDlCard(_ context.Context, in *pb.DlCardRequest) (*pb.DlCardReply, error) {
	log.Printf("[ReqCredDlCard][Request]")

	log.Printf("EIdVcJwt >>\n%v", in.GetEidVcJwt())
	valid, did, err := core.VerifyVp(in.GetEidVcJwt(), controller.GetPublicKey)
	log.Printf("ReqCredDlCard ~~~   %v, %v, err:%v", valid, did, err)
	result := ""
	eDlVcJwt := ""
	if err != nil {
		result = err.Error()
	}
	if valid {
		// ******************** Build VC Claims ******************** //
		subjectDid := did
		nonce := core.RandomString(12)

		typ := "eDriver'sLicenceCardCredential"
		typArray := []string{"VerifiableCredential"}
		typArray = append(typArray, typ)

		credSub := map[string]interface{}{
			"identityinfo": map[string]interface{}{
				"country": "S.Korea",
				"name":    "Hong Gil-Dong",
				"birth":   "2000-11-08",
			},
			"driverLicense": map[string]interface{}{
				"documentnumber":        "15-03-142857-74",
				"certificateprivileges": "1-normal",
				"aptitude_test":         "",
			},
		}
		myVc := map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			"type":              typArray,
			"credentialSubject": credSub,
		}

		standardClaims := jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://www.gov.kr/residentregistration",
			NotBefore: time.Now().Unix(),
			Subject:   subjectDid,
		}

		// Create the Claims
		claims := byd50_jwt.VcClaims{
			nonce,
			myVc,
			standardClaims,
		}
		kid := issuerDid
		pvKey := mustPvKeyECDSA(myDkms)
		eDlVcJwt = core.CreateVcWithClaims(kid, claims, pvKey)
	}

	log.Printf("[ReqCredDlCard][Reply] vcJwt: %v", eDlVcJwt)
	return &pb.DlCardReply{Valid: valid, Result: result, EdlVcJwt: eDlVcJwt}, nil
}

// ReqCredRentalCarAgreement implements proto-files.GreeterServer
func (s *server) ReqCredRentalCarAgreement(_ context.Context, in *pb.RentalCarAgreementRequest) (*pb.RentalCarAgreementReply, error) {
	valid, did, err := core.VerifyVp(in.GetEdlVcJwt(), controller.GetPublicKey)
	result := ""
	rentalCarAgreementVcJwt := ""
	if err != nil {
		result = err.Error()
	}
	if valid {
		// ******************** Build VC Claims ******************** //
		subjectDid := did
		nonce := core.RandomString(12)

		typ := "RentalCarAgreementCredential"
		typArray := []string{"VerifiableCredential"}
		typArray = append(typArray, typ)

		credSub := map[string]interface{}{
			"identityinfo": map[string]interface{}{
				"country": "S.Korea",
				"name":    "Hong Gil-Dong",
				"birth":   "2000-11-08",
			},
			"driverLicense": map[string]interface{}{
				"documentnumber":        "15-03-142857-74",
				"certificateprivileges": "1-normal",
				"aptitude_test":         "",
			},
			"rentalCarInfo": map[string]interface{}{
				"numberPlate": "49ho 2832",
			},
		}
		myVc := map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			"type":              typArray,
			"credentialSubject": credSub,
		}

		standardClaims := jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: time.Now().Add(time.Second * 15).Unix(),
			Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://www.gov.kr/residentregistration",
			NotBefore: time.Now().Unix(),
			Subject:   subjectDid,
		}

		// Create the Claims
		claims := byd50_jwt.VcClaims{
			nonce,
			myVc,
			standardClaims,
		}
		kid := issuerDid
		pvKey := mustPvKeyECDSA(myDkms)
		rentalCarAgreementVcJwt = core.CreateVcWithClaims(kid, claims, pvKey)
	}
	log.Printf("[ReqCredRentalCarAgreement][Reply] rentalCarAgreementVcJwt: %v", rentalCarAgreementVcJwt)

	return &pb.RentalCarAgreementReply{Valid: valid, Result: result, RentalCarAgreementVcJwt: rentalCarAgreementVcJwt}, nil
}

// RentalCarControl implements proto-files.GreeterServer
func (s *server) RentalCarControl(_ context.Context, in *pb.RentalCarControlRequest) (*pb.RentalCarControlReply, error) {
	log.Printf("[RentalCarControl][Request]")
	log.Printf("GetRentalCarAgreementVpJwt >>\n%v", in.GetRentalCarAgreementVcJwt())
	result := ""
	valid, did, err := core.VerifyVp(in.GetRentalCarAgreementVcJwt(), controller.GetPublicKey)
	if valid {
		result = "Welcome to our rental car system. " + did
	}
	if err != nil {
		result = err.Error()
	}

	return &pb.RentalCarControlReply{Valid: valid, Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", configs.UseConfig.IssuerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Initialize KMS
	myDkms, err = kms.InitKMS(kms.KeyTypeECDSA)
	log.Printf("pvkey : %v", mustPvKeyECDSA(myDkms))
	if err != nil {
		log.Fatalf("could not Init KMS (%v)", err.Error())
	}

	// Create DID
	method := "byd50"
	did := controller.CreateDID(myDkms.PbKeyBase58(), method)
	myDkms.SetDid(did)
	issuerDid = did

	s := grpc.NewServer()
	pb.RegisterIssuerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
