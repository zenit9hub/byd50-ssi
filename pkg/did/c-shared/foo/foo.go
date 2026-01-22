package foo

import (
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/golang-jwt/jwt"
	"time"
)

func CreateKeyPairForAndr() {
	_, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		panic(err)
	}
}

func GetPrivateKeyBase58() string {
	dkms := core.GetDKMS()
	return dkms.PvKeyBase58()
}

func GetPublicKeyBase58() string {
	dkms := core.GetDKMS()
	return dkms.PbKeyBase58()
}

func CreateVpForAndr(did, iss, pvKeyBase58, credTyp, vcJwt string) string {
	holderDid := did
	issuer := did
	if len(iss) > 0 {
		issuer = iss
	}
	holderPvKey, _ := x509.ParseECPrivateKey(base58.Decode(pvKeyBase58))

	// ******************** Build VP Claims ******************** //
	nonce := core.RandomString(12)

	typ := credTyp
	typArray := []string{"VerifiableCredential"}
	typArray = append(typArray, typ)

	var vcJwtArray []string
	vcJwtArray = append(vcJwtArray, vcJwt)

	myVp := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":                 typArray,
		"verifiableCredential": vcJwtArray,
	}

	standardClaims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
		IssuedAt:  time.Now().Unix(),
		Issuer:    issuer,
		NotBefore: time.Now().Unix(),
		Subject:   "",
	}

	// Create the Claims
	claims := byd50_jwt.VpClaims{
		Nonce:          nonce,
		Vp:             myVp,
		StandardClaims: standardClaims,
	}

	vpJwt := core.CreateVpWithClaims(holderDid, claims, holderPvKey)

	return vpJwt
}

func ClaimsGetExp(vpJwt string) int64 {
	return ClaimsGetInt64(vpJwt, "exp")
}

func ClaimsGetIat(vpJwt string) int64 {
	return ClaimsGetInt64(vpJwt, "iat")
}

func ClaimsGetInt64(vpJwt, claim string) int64 {
	claimInt64 := int64(0)
	parseToken, _, err := new(jwt.Parser).ParseUnverified(vpJwt, jwt.MapClaims{})
	if err != nil {
		return claimInt64
	}

	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok {
		vClaim, ok := claims[claim]
		if !ok {
			err = errors.New(fmt.Sprintf("claims hasn't [%v] field", claim))
		}

		switch claimType := vClaim.(type) {
		case float64:
			claimInt64 = int64(claimType)
		case json.Number:
			claimInt64, _ = claimType.Int64()
		default:
			err = errors.New(fmt.Sprintf("'%v' type error", claim))
		}
	}

	return claimInt64
}
