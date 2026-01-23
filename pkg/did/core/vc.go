package core

import (
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	derrors "byd50-ssi/pkg/did/errors"
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt"
	"time"
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
	if err != nil {
		return false, err
	}
	if !ok {
		return false, derrors.New(derrors.CodeInvalidInput, "vc signature invalid")
	}
	return true, nil
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

// ValidateVcClaims ensures required standard claims are present and consistent.
func ValidateVcClaims(claims byd50_jwt.VcClaims) error {
	if claims.Issuer == "" {
		return derrors.New(derrors.CodeInvalidInput, "vc issuer is empty")
	}
	if claims.ExpiresAt == 0 {
		return derrors.New(derrors.CodeInvalidInput, "vc exp is empty")
	}
	if claims.IssuedAt == 0 {
		return derrors.New(derrors.CodeInvalidInput, "vc iat is empty")
	}
	if claims.ExpiresAt <= claims.IssuedAt {
		return derrors.New(derrors.CodeInvalidInput, "vc exp must be after iat")
	}
	if claims.NotBefore != 0 && claims.NotBefore > claims.IssuedAt {
		return derrors.New(derrors.CodeInvalidInput, "vc nbf must be <= iat")
	}
	if claims.ExpiresAt <= time.Now().Unix() {
		return derrors.New(derrors.CodeInvalidInput, "vc expired")
	}
	return nil
}
