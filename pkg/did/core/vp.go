package core

import (
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	derrors "byd50-ssi/pkg/did/errors"
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
	if err != nil {
		return false, did, err
	}
	if !ok {
		return false, did, derrors.New(derrors.CodeInvalidInput, "vp signature invalid")
	}
	return true, did, nil
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

// ValidateVpClaims ensures required standard claims are present and consistent.
func ValidateVpClaims(claims byd50_jwt.VpClaims) error {
	if claims.Issuer == "" {
		return derrors.New(derrors.CodeInvalidInput, "vp issuer is empty")
	}
	if claims.ExpiresAt == 0 {
		return derrors.New(derrors.CodeInvalidInput, "vp exp is empty")
	}
	if claims.IssuedAt == 0 {
		return derrors.New(derrors.CodeInvalidInput, "vp iat is empty")
	}
	if claims.ExpiresAt <= claims.IssuedAt {
		return derrors.New(derrors.CodeInvalidInput, "vp exp must be after iat")
	}
	if claims.NotBefore != 0 && claims.NotBefore > claims.IssuedAt {
		return derrors.New(derrors.CodeInvalidInput, "vp nbf must be <= iat")
	}
	if claims.ExpiresAt <= time.Now().Unix() {
		return derrors.New(derrors.CodeInvalidInput, "vp expired")
	}
	return nil
}
