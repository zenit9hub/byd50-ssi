package byd50_jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type VcClaims struct {

	// The prn (principal) claim identifies the subject of the JWT.
	//Prn string `json:"prn"`

	// The typ (type) claim is used to declare a type for the contents of this JWT Claims Set.
	//Typ string `json:"typ"`

	// Nonce is used only once and can't be used in second time.
	Nonce string `json:"nonce,omitempty"`

	Vc map[string]interface{} `json:"vc,omitempty"`

	// The aud (audience) claim identifies the audience that the JWT is intended for.
	// exp MUST represent the expirationDate property, encoded as a UNIX timestamp (NumericDate).
	// jti MUST represent the id property of the verifiable credential or verifiable presentation.
	// The iat (issued at) claim identifies the time at which the JWT was issued. This claim can be used to determine the age of the token
	// iss MUST represent the issuer property of a verifiable credential or the holder property of a verifiable presentation.
	// nbf MUST represent issuanceDate, encoded as a UNIX timestamp (NumericDate).
	// sub MUST represent the id property contained in the verifiable credential subject. eg> did:example:ebfeb1f712ebc6f1c276e12ec21
	jwt.StandardClaims
}

func MakeVcSample(kid string, pvKey *ecdsa.PrivateKey) string {
	typ := []string{"VerifiableCredential", "AlumniCredential"}
	credSub := map[string]interface{}{
		"degree": "BachelorDegree",
		"name":   "<span lang='fr-CA'>Baccalauréat en musiques numériques</span>",
	}

	myVc := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":              typ,
		"credentialSubject": credSub,
	}

	nonce := "nonce-142857"
	aud := ""
	exp := time.Now().Add(time.Minute * 5).Unix()
	jti := "089a411f-0d88-450f-8cc0-1a3acfebecd3"
	iat := time.Now().Unix()
	nbf := iat
	iss := "http://google.com/issuer"
	sub := ""

	// Create the Claims
	claims := VcClaims{
		nonce,
		myVc,
		jwt.StandardClaims{
			Audience:  aud,
			ExpiresAt: exp,
			Id:        jti,
			IssuedAt:  iat,
			Issuer:    iss,
			NotBefore: nbf,
			Subject:   sub,
		},
	}

	vcSampleJwt := CreateVc(kid, claims, pvKey)
	return vcSampleJwt
}

func CreateVc(kid string, claims VcClaims, pvKey *ecdsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = kid

	ss, err := token.SignedString(pvKey)
	if err != nil {
		log.Printf(err.Error())
	}
	return ss
}

func VerifyVc(vcJwt string, getPbKey func(string, string) string) (bool, error) {
	parseToken, err := jwt.Parse(vcJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		did, ok := token.Header["kid"].(string)
		if !ok || did == "" {
			return nil, errors.New("missing kid header")
		}
		pbKeyBase58 := getPbKey(did, "")
		pbKeyBytes := base58.Decode(pbKeyBase58)
		if len(pbKeyBytes) == 0 {
			return nil, errors.New("invalid public key base58")
		}
		pbKey, err := x509.ParsePKIXPublicKey(pbKeyBytes)
		if err != nil {
			return nil, err
		}
		return pbKey, nil
	})
	return parseToken.Valid, err
}
