package byd50_jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type VpClaims struct {
	// Nonce is used only once and can't be used in second time.
	Nonce string `json:"nonce,omitempty"`

	Vp map[string]interface{} `json:"vp,omitempty"`

	// The aud (audience) claim identifies the audience that the JWT is intended for.
	// exp MUST represent the expirationDate property, encoded as a UNIX timestamp (NumericDate).
	// jti MUST represent the id property of the verifiable credential or verifiable presentation.
	// The iat (issued at) claim identifies the time at which the JWT was issued. This claim can be used to determine the age of the token
	// iss MUST represent the issuer property of a verifiable credential or the holder property of a verifiable presentation.
	// nbf MUST represent issuanceDate, encoded as a UNIX timestamp (NumericDate).
	// sub MUST represent the id property contained in the verifiable credential subject. eg> did:example:ebfeb1f712ebc6f1c276e12ec21
	jwt.StandardClaims
}

func MakeVpSample(issuerDid string, vcJwtArray []string, pvKey *ecdsa.PrivateKey) string {
	typ := []string{"VerifiableCredential", "AlumniCredential"}
	myVp := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":                 typ,
		"verifiableCredential": vcJwtArray,
	}

	nonce := "nonce-142857"
	aud := "did:example:4a57546973436f6f6c4a4a57573"
	exp := time.Now().Add(time.Minute * 5).Unix()
	jti := "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c5"
	iat := time.Now().Unix()
	nbf := iat
	iss := issuerDid
	sub := ""

	// Create the Claims
	claims := VpClaims{
		nonce,
		myVp,
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

	kid := issuerDid
	vpSampleJwt := CreateVp(kid, claims, pvKey)
	return vpSampleJwt
}

func CreateVp(kid string, claims VpClaims, pvKey *ecdsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = kid

	ss, err := token.SignedString(pvKey)
	if err != nil {
		log.Printf(err.Error())
	}
	return ss
}

func VerifyVp(vpJwt string, getPbKey func(string, string) string) (bool, string, error) {
	valid := false
	did := ""
	parseToken, err := jwt.Parse(vpJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		did = token.Header["kid"].(string)
		pbKeyBase58 := getPbKey(did, "")
		pbKey, _ := x509.ParsePKIXPublicKey(base58.Decode(pbKeyBase58))
		return pbKey, nil
	})
	if err != nil {
		return valid, did, err
	}
	valid = parseToken.Valid

	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
		if claims["vp"] != nil {
			vpMapClaims := claims["vp"].(map[string]interface{})
			var vcJwtArray []string
			switch v := vpMapClaims["verifiableCredential"].(type) {
			case string:
				vcJwtArray = append(vcJwtArray, v)
			case []string:
				vcJwtArray = v
			case []interface{}:
				for _, a := range v {
					vs, ok := a.(string)
					if !ok {
						break
					}
					vcJwtArray = append(vcJwtArray, vs)
				}
			}

			log.Printf("Verify each elements.")
			for index, element := range vcJwtArray {
				ok, err = VerifyVc(element, getPbKey)
				if err != nil {
					valid = false
					log.Printf(err.Error())
					break
				}
				if ok {
					log.Printf("index[%v] proof verified correctly", index)
				}
			}
		} else {
			log.Printf("\n\n\n\n \t vpMapClaims not exist ~~~~~~~~~~~~\n\n\n\n")
		}
	}

	return valid, did, err
}

func ParseVp(vpJwt string, getPbKey func(string, string) string) (bool, jwt.MapClaims, error) {
	parseToken, err := jwt.Parse(vpJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		did := token.Header["kid"].(string)
		pbKeyBase58 := getPbKey(did, "")
		pbKey, _ := x509.ParsePKIXPublicKey(base58.Decode(pbKeyBase58))
		return pbKey, nil
	})
	if claims, ok := parseToken.Claims.(jwt.MapClaims); err == nil && ok && parseToken.Valid {
		return parseToken.Valid, claims, err
	}
	return parseToken.Valid, nil, err
}

func ClaimsGetExp(claims jwt.MapClaims) (time.Time, error) {
	var err error
	var expiresAt int64
	t := time.Time{}

	exp, ok := claims["exp"]
	if !ok {
		return t, errors.New("claims hasn't exp field")
	}
	switch expType := exp.(type) {
	case float64:
		expiresAt = int64(expType)
		t = time.Unix(expiresAt, 0)
	case json.Number:
		expiresAt, _ = expType.Int64()
		t = time.Unix(expiresAt, 0)
	default:
		err = errors.New("'exp' type error")
	}

	return t, err
}

func ClaimsGetIat(claims jwt.MapClaims) (time.Time, error) {
	var err error
	var issuedAt int64
	t := time.Time{}

	iat, ok := claims["iat"]
	if !ok {
		return t, errors.New("claims hasn't iat field")
	}
	switch iatType := iat.(type) {
	case float64:
		issuedAt = int64(iatType)
		t = time.Unix(issuedAt, 0)
	case json.Number:
		issuedAt, _ = iatType.Int64()
		t = time.Unix(issuedAt, 0)
	default:
		err = errors.New("'iat' type error")
	}

	return t, err
}
