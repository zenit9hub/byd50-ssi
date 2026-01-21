package core

import (
	byd50_jwt "byd50-ssi/did/core/byd50-jwt"
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt"
)

func CreateVc(kid, typ string, credSub map[string]interface{}, standardClaims jwt.StandardClaims, pvKey *ecdsa.PrivateKey) string {
	claims := buildVcClaims(typ, credSub, standardClaims)
	vcSampleJwt := byd50_jwt.CreateVc(kid, claims, pvKey)
	return vcSampleJwt
}

func CreateVcWithClaims(kid string, claims byd50_jwt.VcClaims, pvKey *ecdsa.PrivateKey) string {
	vcSampleJwt := byd50_jwt.CreateVc(kid, claims, pvKey)
	return vcSampleJwt
}

func VerifyVc(vc string, getPbKey func(string, string) string) (bool, error) {
	ok, err := byd50_jwt.VerifyVc(vc, getPbKey)
	return ok, err
}

func buildVcClaims(typ string, credSub map[string]interface{}, standardClaims jwt.StandardClaims) byd50_jwt.VcClaims {
	typArray := []string{"VerifiableCredential"}
	if typ != "" {
		typArray = append(typArray, typ)
	}

	myVc := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":              typArray,
		"credentialSubject": credSub,
	}

	nonce := RandomString(12)

	// Create the Claims
	claims := byd50_jwt.VcClaims{
		nonce,
		myVc,
		standardClaims,
	}

	return claims
}
