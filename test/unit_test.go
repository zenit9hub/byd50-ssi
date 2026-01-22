//go:build integration
// +build integration

package main_test

import (
	"byd50-ssi/pkg/did/core"
	"byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/pkg/controller"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	"testing"
	"time"
)

func TestInitDKMS(t *testing.T) {
	// ******************** KeyTypeRSA Test ******************** //
	dkms, err := core.InitDKMS(core.KeyTypeRSA)
	if err != nil {
		t.Fatal(err.Error())
	}
	pvKey, err := dkms.PvKeyRSA()
	if err != nil {
		t.Fatal(err)
	}
	err = pvKey.Validate()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestEncryptDecrypt(t *testing.T) {
	myDkms := core.GetDKMS()
	plainText := "TestEncryptDecrypt"
	encryptedText := myDkms.Encrypt(plainText)
	decryptedText := myDkms.Decrypt(encryptedText)
	if plainText != decryptedText {
		t.Fatal(errors.New("plainText and decryptedText are not same"))
	}
}

func TestSignVerify(t *testing.T) {
	myDkms := core.GetDKMS()
	plainText := "TestEncryptDecrypt"
	ret, result := core.PvKeySign(myDkms.PvKeyBase58(), plainText, "")
	if !ret {
		t.Fatal(result)
	}
	ret = core.PbKeyVerify(myDkms.PbKeyBase58(), plainText, result)
	if !ret {
		t.Fatal(errors.New("PbKeyVerify error"))
	}
}

func TestVerifyAuthChallengeAndResponse(t *testing.T) {
	myDkms := core.GetDKMS()
	method := "byd50"
	did := controller.CreateDID(myDkms.PbKeyBase58(), method)
	plainText := "TestVerifyAuthChallengeAndResponse"
	authChallengeString := controller.GetAuthChallengeString(did, plainText)
	authResponseString := controller.GetAuthResponseString(authChallengeString, myDkms.PvKeyBase58())
	if plainText != authResponseString {
		t.Fatal(errors.New("TestVerifyAuthChallengeAndResponse error"))
	}
}

func TestVerifySimplePresent(t *testing.T) {
	myDkms := core.GetDKMS()
	method := "byd50"
	did := controller.CreateDID(myDkms.PbKeyBase58(), method)
	myDkms.SetDid(did)
	simplePresentString := controller.GetSimplePresent(myDkms.Did(), myDkms.PvKeyBase58())
	result := controller.VerifySimplePresent(simplePresentString)
	if result != "success" {
		t.Fatal(errors.New("TestVerifySimplePresent error"))
	}
}

func TestKeyExport(t *testing.T) {
	myDkms := core.GetDKMS()
	pvKeyPem := core.ExportPrivateKeyAsPEM(myDkms.PvKey())
	if pvKeyPem != myDkms.PvKeyPEM() {
		t.Fatal(errors.New("export result error"))
	}

	block, _ := pem.Decode([]byte(pvKeyPem))
	if block == nil {
		t.Fatal(errors.New("failed to parse PEM block containing the private key"))
	}

	t.Logf("block.Type=[%v], PbKeyPEM=[%v]", block.Type, myDkms.PbKeyPEM())
}

func TestCreateAndResolveDID(t *testing.T) {
	_, err := core.InitDKMS(core.KeyTypeRSA)
	if err != nil {
		t.Fatal(err.Error())
	}
	myDkms := core.GetDKMS()
	method := "byd50"
	did := controller.CreateDID(myDkms.PbKeyBase58(), method)
	myDkms.SetDid(did)
	if did == "" {
		t.Fatal(errors.New("did is null"))
	}
	didDocument := controller.ResolveDID(did)
	if didDocument == "" {
		t.Fatal(errors.New("document is null"))
	}

}

func TestGetPublicKey(t *testing.T) {
	myDkms := core.GetDKMS()
	method := "byd50"
	did := controller.CreateDID(myDkms.PbKeyBase58(), method)
	pbKey := controller.GetPublicKey(did, "")
	if pbKey == "" {
		t.Fatal(errors.New("public key is null"))
	}
}

func TestMakeVcAndVerify(t *testing.T) {
	// ******************** KeyTypeECDSA ******************** //
	issuerDkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}

	// ******************** Create DID ******************** //
	method := "byd50"
	did := controller.CreateDID(issuerDkmsEcdsa.PbKeyBase58(), method)
	issuerDkmsEcdsa.SetDid(did)
	pvKey := issuerDkmsEcdsa.PvKey().(*ecdsa.PrivateKey)

	// ******************** Build VC Claims ******************** //
	nonce := core.RandomString(12)

	typ := "AlumniCredential"
	typArray := []string{"VerifiableCredential"}
	typArray = append(typArray, typ)

	credSub := map[string]interface{}{
		"degree": "BachelorDegree",
		"name":   "<span lang='fr-CA'>Baccalauréat en musiques numériques</span>",
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
		myVc,
		standardClaims,
	}

	// ******************** Create VC ******************** //
	vcJwt := core.CreateVc(did, typ, credSub, standardClaims, pvKey)
	vcJwt2 := core.CreateVcWithClaims(did, claims, pvKey)
	t.Logf(" -- vc sample jwt --\n%v\n%v", vcJwt, vcJwt2)

	// ******************** Verify VC ******************** //
	valid, err := core.VerifyVc(vcJwt, controller.GetPublicKey)
	t.Logf("token.Valid (%v), err (%v)", valid, err)
}

func TestMakeVpAndVerify(t *testing.T) {
	// ******************** Sequence 1 preparing for vc ******************** //
	// ******************** KeyTypeECDSA ******************** //
	issuerDkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}
	// ******************** Create DID ******************** //
	method := "byd50"
	issuerDid := controller.CreateDID(issuerDkmsEcdsa.PbKeyBase58(), method)
	issuerDkmsEcdsa.SetDid(issuerDid)
	issuerPvKey := issuerDkmsEcdsa.PvKey().(*ecdsa.PrivateKey)

	// ******************** Create VC ******************** //
	kid := issuerDid
	vcJwt := byd50_jwt.MakeVcSample(kid, issuerPvKey)
	t.Logf(" -- vc sample jwt --\n%v", vcJwt)

	// ******************** Verify VC ******************** //
	valid, err := core.VerifyVc(vcJwt, controller.GetPublicKey)
	t.Logf("token.Valid (%v), err (%v)", valid, err)

	holderDkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}

	// ******************** Sequence 2 preparing for vp ******************** //
	// ******************** Create DID ******************** //
	holderDid := controller.CreateDID(holderDkmsEcdsa.PbKeyBase58(), method)
	holderDkmsEcdsa.SetDid(holderDid)
	holderPvKey := holderDkmsEcdsa.PvKey().(*ecdsa.PrivateKey)

	// ******************** Build VP Claims ******************** //
	typ := "CredentialManagerPresentation"

	var vcJwtArray []string
	vcJwtArray = append(vcJwtArray, vcJwt)

	standardClaims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Id:        "089a411f-0d88-450f-8cc0-1a3acfebecd3",
		IssuedAt:  time.Now().Unix(),
		Issuer:    "http://google.com/issuer",
		NotBefore: time.Now().Unix(),
		Subject:   "",
	}

	// ******************** Create VP ******************** //
	vpJwt := core.CreateVp(holderDid, typ, vcJwtArray, standardClaims, holderPvKey)
	t.Logf(" -- vp sample jwt --\n%v", vpJwt)

	// ******************** Verify VP ******************** //
	valid, _, err = core.VerifyVp(vpJwt, controller.GetPublicKey)
	t.Logf("token.Valid (%v), err (%v)", valid, err)

	if ok, mapclaims, err := core.GetMapClaims(vpJwt, controller.GetPublicKey); ok && err == nil {
		expTime, err := core.ClaimsGetExp(mapclaims)
		t.Logf("expTime (%v), err (%v)", expTime, err)

		mc := byd50_jwt.MapClaims(mapclaims)
		aud, _ := mc.GetAudience()
		exp, _ := mc.GetExpiresAt()
		iat, _ := mc.GetIssuedAt()
		iss, _ := mc.GetIssuer()
		nbf, _ := mc.GetNotBefore()
		vpTyp, _ := mc.GetVpType()

		t.Logf("aud(%v), exp(%v), iat(%v), iss(%v), nbf(%v), vpTyp(%v)", aud, exp, iat, iss, nbf, vpTyp)
	}
}
