package core

import (
	byd50_jwt "byd50-ssi/did/core/byd50-jwt"
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt"
	"time"
)

func CreateVp(kid, typ string, vcJwtArray []string, standardClaims jwt.StandardClaims, pvKey *ecdsa.PrivateKey) string {
	claims := buildVpClaims(typ, vcJwtArray, standardClaims)
	vpJwt := byd50_jwt.CreateVp(kid, claims, pvKey)
	return vpJwt
}

func CreateVpWithClaims(kid string, claims byd50_jwt.VpClaims, pvKey *ecdsa.PrivateKey) string {
	vpJwt := byd50_jwt.CreateVp(kid, claims, pvKey)
	return vpJwt
}

func VerifyVp(vp string, getPbKey func(string, string) string) (bool, string, error) {
	ok, did, err := byd50_jwt.VerifyVp(vp, getPbKey)
	return ok, did, err
}

func GetMapClaims(vp string, getPbKey func(string, string) string) (bool, jwt.MapClaims, error) {
	ok, mapClaims, err := byd50_jwt.ParseVp(vp, getPbKey)
	return ok, mapClaims, err
}

func ClaimsGetExp(mapclaims jwt.MapClaims) (time.Time, error) {
	expTime, err := byd50_jwt.ClaimsGetExp(mapclaims)
	return expTime, err
}

func ClaimsGetIat(mapclaims jwt.MapClaims) (time.Time, error) {
	iatTime, err := byd50_jwt.ClaimsGetIat(mapclaims)
	return iatTime, err
}

func buildVpClaims(typ string, vcJwtArray []string, standardClaims jwt.StandardClaims) byd50_jwt.VpClaims {
	typArray := []string{"VerifiablePresentation"}
	if typ != "" {
		typArray = append(typArray, typ)
	}

	myVp := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":                 typArray,
		"verifiableCredential": vcJwtArray,
	}

	nonce := RandomString(12)

	// Create the Claims
	claims := byd50_jwt.VpClaims{
		nonce,
		myVp,
		standardClaims,
	}

	return claims
}
