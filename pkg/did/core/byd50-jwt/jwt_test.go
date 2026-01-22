package byd50_jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/golang-jwt/jwt"
)

func TestVcVpJwtFlow(t *testing.T) {
	pvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	did := "did:byd50:test"
	pbBytes, err := x509.MarshalPKIXPublicKey(&pvKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pbKeyBase58 := base58.Encode(pbBytes)
	getPbKey := func(_ string, _ string) string {
		return pbKeyBase58
	}

	vcJwt := MakeVcSample(did, pvKey)
	if vcJwt == "" {
		t.Fatal("vc jwt empty")
	}
	if ok, err := VerifyVc(vcJwt, getPbKey); !ok || err != nil {
		t.Fatalf("vc verify failed: %v", err)
	}

	vpJwt := MakeVpSample(did, []string{vcJwt}, pvKey)
	if ok, _, err := VerifyVp(vpJwt, getPbKey); !ok || err != nil {
		t.Fatalf("vp verify failed: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{Issuer: did})
	token.Header["kid"] = did
	noVpJwt, err := token.SignedString(pvKey)
	if err != nil {
		t.Fatal(err)
	}
	if ok, _, err := VerifyVp(noVpJwt, getPbKey); !ok || err != nil {
		t.Fatalf("vp verify without vp failed: %v", err)
	}

	vpString := VpClaims{
		Nonce: "n1",
		Vp: map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
			},
			"type":                 []string{"VerifiablePresentation"},
			"verifiableCredential": vcJwt,
		},
		StandardClaims: jwt.StandardClaims{Issuer: did, ExpiresAt: time.Now().Add(time.Minute).Unix()},
	}
	vpStringJwt := CreateVp(did, vpString, pvKey)
	if ok, _, err := VerifyVp(vpStringJwt, getPbKey); !ok || err != nil {
		t.Fatalf("vp verify string failed: %v", err)
	}

	vpInterface := VpClaims{
		Nonce: "n2",
		Vp: map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
			},
			"type":                 []string{"VerifiablePresentation"},
			"verifiableCredential": []interface{}{vcJwt},
		},
		StandardClaims: jwt.StandardClaims{Issuer: did, ExpiresAt: time.Now().Add(time.Minute).Unix()},
	}
	vpInterfaceJwt := CreateVp(did, vpInterface, pvKey)
	if ok, _, err := VerifyVp(vpInterfaceJwt, getPbKey); !ok || err != nil {
		t.Fatalf("vp verify interface failed: %v", err)
	}

	if ok, claims, err := ParseVp(vpJwt, getPbKey); !ok || err != nil || claims == nil {
		t.Fatalf("parse vp failed: %v", err)
	}

	vpBadCredential := VpClaims{
		Nonce: "n3",
		Vp: map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
			},
			"type":                 []string{"VerifiablePresentation"},
			"verifiableCredential": []interface{}{123},
		},
		StandardClaims: jwt.StandardClaims{Issuer: did, ExpiresAt: time.Now().Add(time.Minute).Unix()},
	}
	vpBadCredentialJwt := CreateVp(did, vpBadCredential, pvKey)
	if ok, _, err := VerifyVp(vpBadCredentialJwt, getPbKey); !ok || err != nil {
		t.Fatalf("vp verify bad credential failed: %v", err)
	}
}

func TestClaimsHelpers(t *testing.T) {
	now := time.Now().Unix()
	claims := MapClaims{
		"aud":      []interface{}{"a", "b"},
		"audArr":   []string{"c", "d"},
		"audStr":   "e",
		"exp":      json.Number("123"),
		"expFloat": float64(321),
		"iat":      json.Number("456"),
		"iss":      "issuer",
		"nbf":      json.Number("789"),
		"vc":       map[string]interface{}{"type": []interface{}{"V1", "V2"}},
		"vp":       map[string]interface{}{"type": []interface{}{"P1", "P2"}},
	}

	if aud, err := claims.GetAudience(); err != nil || len(aud) != 2 {
		t.Fatalf("audience error: %v", err)
	}
	if aud, err := (MapClaims{"aud": claims["audArr"]}).GetAudience(); err != nil || len(aud) != 2 {
		t.Fatalf("audience []string error: %v", err)
	}
	if aud, err := (MapClaims{"aud": claims["audStr"]}).GetAudience(); err != nil || len(aud) != 1 {
		t.Fatalf("audience string error: %v", err)
	}
	if exp, err := claims.GetExpiresAt(); err != nil || exp != 123 {
		t.Fatalf("exp error: %v", err)
	}
	if exp, err := (MapClaims{"exp": claims["expFloat"]}).GetExpiresAt(); err != nil || exp != 321 {
		t.Fatalf("exp float error: %v", err)
	}
	if _, err := (MapClaims{}).GetExpiresAt(); err == nil {
		t.Fatal(errors.New("expected missing exp error"))
	}
	if iat, err := claims.GetIssuedAt(); err != nil || iat != 456 {
		t.Fatalf("iat error: %v", err)
	}
	if iss, err := claims.GetIssuer(); err != nil || iss != "issuer" {
		t.Fatalf("iss error: %v", err)
	}
	if nbf, err := claims.GetNotBefore(); err != nil || nbf != 789 {
		t.Fatalf("nbf error: %v", err)
	}
	if _, err := claims.GetVc(); err != nil {
		t.Fatalf("vc error: %v", err)
	}
	if vtyp, err := claims.GetVcType(); err != nil || len(vtyp) != 2 {
		t.Fatalf("vc type error: %v", err)
	}
	if vtyp, err := (MapClaims{"vc": map[string]interface{}{"type": "V1"}}).GetVcType(); err != nil || len(vtyp) != 1 {
		t.Fatalf("vc type string error: %v", err)
	}
	if vtyp, err := (MapClaims{"vc": map[string]interface{}{"type": []string{"V1", "V2"}}}).GetVcType(); err != nil || len(vtyp) != 2 {
		t.Fatalf("vc type []string error: %v", err)
	}
	if _, err := claims.GetVp(); err != nil {
		t.Fatalf("vp error: %v", err)
	}
	if vtyp, err := claims.GetVpType(); err != nil || len(vtyp) != 2 {
		t.Fatalf("vp type error: %v", err)
	}
	if vtyp, err := (MapClaims{"vp": map[string]interface{}{"type": "P1"}}).GetVpType(); err != nil || len(vtyp) != 1 {
		t.Fatalf("vp type string error: %v", err)
	}
	if vtyp, err := (MapClaims{"vp": map[string]interface{}{"type": []string{"P1", "P2"}}}).GetVpType(); err != nil || len(vtyp) != 2 {
		t.Fatalf("vp type []string error: %v", err)
	}
	if _, err := (MapClaims{"vc": map[string]interface{}{}}).GetVcType(); err == nil {
		t.Fatal(errors.New("expected vc type missing error"))
	}
	if _, err := (MapClaims{"vp": map[string]interface{}{}}).GetVpType(); err == nil {
		t.Fatal(errors.New("expected vp type missing error"))
	}
	if _, err := (MapClaims{"vc": map[string]interface{}{"type": 123}}).GetVcType(); err == nil {
		t.Fatal(errors.New("expected vc type error"))
	}
	if _, err := (MapClaims{"vp": map[string]interface{}{"type": 123}}).GetVpType(); err == nil {
		t.Fatal(errors.New("expected vp type error"))
	}
	if _, err := (MapClaims{"aud": []interface{}{"a", 1}}).GetAudience(); err == nil {
		t.Fatal(errors.New("expected audience type error"))
	}

	if _, err := ClaimsGetExp(jwt.MapClaims{"exp": float64(now)}); err != nil {
		t.Fatalf("claims exp error: %v", err)
	}
	if _, err := ClaimsGetIat(jwt.MapClaims{"iat": float64(now)}); err != nil {
		t.Fatalf("claims iat error: %v", err)
	}
	if _, err := ClaimsGetExp(jwt.MapClaims{"exp": json.Number("123")}); err != nil {
		t.Fatalf("claims exp number error: %v", err)
	}
	if _, err := ClaimsGetIat(jwt.MapClaims{"iat": json.Number("456")}); err != nil {
		t.Fatalf("claims iat number error: %v", err)
	}
	if _, err := ClaimsGetExp(jwt.MapClaims{}); err == nil {
		t.Fatal(errors.New("expected exp error"))
	}
	if _, err := ClaimsGetIat(jwt.MapClaims{}); err == nil {
		t.Fatal(errors.New("expected iat error"))
	}
	if _, err := ClaimsGetExp(jwt.MapClaims{"exp": "bad"}); err == nil {
		t.Fatal(errors.New("expected exp type error"))
	}
	if _, err := ClaimsGetIat(jwt.MapClaims{"iat": "bad"}); err == nil {
		t.Fatal(errors.New("expected iat type error"))
	}
}

func TestClaimErrors(t *testing.T) {
	claims := MapClaims{
		"vc": map[string]interface{}{"type": []interface{}{123}},
		"vp": map[string]interface{}{"type": []interface{}{true}},
	}

	if _, err := claims.GetVcType(); err == nil {
		t.Fatal(errors.New("expected vc type error"))
	}
	if _, err := claims.GetVpType(); err == nil {
		t.Fatal(errors.New("expected vp type error"))
	}

	if _, err := (MapClaims{}).GetVc(); err == nil {
		t.Fatal(errors.New("expected missing vc error"))
	}
	if _, err := (MapClaims{}).GetVp(); err == nil {
		t.Fatal(errors.New("expected missing vp error"))
	}
}

func TestJwtSigningMethodErrors(t *testing.T) {
	pvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	did := "did:byd50:test"
	pbBytes, err := x509.MarshalPKIXPublicKey(&pvKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pbKeyBase58 := base58.Encode(pbBytes)
	getPbKey := func(_ string, _ string) string {
		return pbKeyBase58
	}

	badVcToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{Issuer: did})
	badVcToken.Header["kid"] = did
	badVcJwt, err := badVcToken.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := VerifyVc(badVcJwt, getPbKey); ok || err == nil {
		t.Fatal("expected vc signing method error")
	}

	badVpToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{Issuer: did})
	badVpToken.Header["kid"] = did
	badVpJwt, err := badVpToken.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	if ok, _, err := VerifyVp(badVpJwt, getPbKey); ok || err == nil {
		t.Fatal("expected vp signing method error")
	}
	if ok, claims, err := ParseVp(badVpJwt, getPbKey); ok || err == nil || claims != nil {
		t.Fatal("expected parse vp signing method error")
	}
}

func TestJwtMissingKidAndBadKey(t *testing.T) {
	pvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	getEmptyKey := func(_ string, _ string) string {
		return ""
	}

	noKidToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{Issuer: "did:byd50:test"})
	noKidJwt, err := noKidToken.SignedString(pvKey)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := VerifyVc(noKidJwt, getEmptyKey); ok || err == nil {
		t.Fatal("expected missing kid error for vc")
	}
	if ok, _, err := VerifyVp(noKidJwt, getEmptyKey); ok || err == nil {
		t.Fatal("expected missing kid error for vp")
	}
	if ok, claims, err := ParseVp(noKidJwt, getEmptyKey); ok || err == nil || claims != nil {
		t.Fatal("expected missing kid error for parse vp")
	}

	kidToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{Issuer: "did:byd50:test"})
	kidToken.Header["kid"] = "did:byd50:test"
	kidJwt, err := kidToken.SignedString(pvKey)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := VerifyVc(kidJwt, getEmptyKey); ok || err == nil {
		t.Fatal("expected bad public key error for vc")
	}
	if ok, _, err := VerifyVp(kidJwt, getEmptyKey); ok || err == nil {
		t.Fatal("expected bad public key error for vp")
	}
	if ok, claims, err := ParseVp(kidJwt, getEmptyKey); ok || err == nil || claims != nil {
		t.Fatal("expected bad public key error for parse vp")
	}
}

func TestInternalClaimHelpers(t *testing.T) {
	claims := MapClaims{
		"str":       "one",
		"strArr":    []string{"a"},
		"strArrBad": []string{"a", "b"},
		"ifaceArr":  []interface{}{"x"},
		"badIface":  []interface{}{123},
		"intOk":     float64(12),
		"intBad":    "bad",
	}

	if v, err := claims.getString("str"); err != nil || v != "one" {
		t.Fatal("getString string failed")
	}
	if v, err := claims.getString("strArr"); err != nil || v != "a" {
		t.Fatal("getString []string failed")
	}
	if _, err := claims.getString("strArrBad"); err == nil {
		t.Fatal("expected getString []string error")
	}
	if v, err := claims.getString("ifaceArr"); err != nil || v != "x" {
		t.Fatal("getString []interface failed")
	}
	if _, err := claims.getString("badIface"); err == nil {
		t.Fatal("expected getString type error")
	}
	if _, err := claims.getString("missing"); err == nil {
		t.Fatal("expected missing claim error")
	}

	if arr, err := claims.getStringArray("str"); err != nil || len(arr) != 1 {
		t.Fatal("getStringArray string failed")
	}
	if arr, err := claims.getStringArray("strArr"); err != nil || len(arr) != 1 {
		t.Fatal("getStringArray []string failed")
	}
	if arr, err := claims.getStringArray("ifaceArr"); err != nil || len(arr) != 1 {
		t.Fatal("getStringArray []interface failed")
	}
	if _, err := claims.getStringArray("badIface"); err == nil {
		t.Fatal("expected getStringArray type error")
	}
	if _, err := claims.getStringArray("missing"); err == nil {
		t.Fatal("expected getStringArray missing error")
	}

	if v, err := claims.getInt64("intOk"); err != nil || v != 12 {
		t.Fatal("getInt64 float64 failed")
	}
	if _, err := claims.getInt64("intBad"); err == nil {
		t.Fatal("expected getInt64 type error")
	}
	if _, err := claims.getInt64("missing"); err == nil {
		t.Fatal("expected getInt64 missing error")
	}
}
