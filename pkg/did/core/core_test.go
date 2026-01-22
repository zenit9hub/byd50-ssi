package core_test

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/core/dids"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"testing"
	"time"
)

type Test struct {
	in  int
	out string
}

var tests = []Test{
	{-1, "negative"},
	{5, "small"},
}

func TestRandomHex(t *testing.T) {
	hexDigitStr, err := core.RandomHex(20)
	log.Printf("%v", hexDigitStr)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestContains(t *testing.T) {
	if core.Contains(configs.UseConfig.AdoptedDriverList, "intended") {
		t.Fail()
	}
	if !core.Contains(configs.UseConfig.AdoptedDriverList, "byd50") {
		t.Fail()
	}
}

func TestDKMS(t *testing.T) {
	dkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}
	pvKeyEcdsa, err := dkmsEcdsa.PvKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}
	pbKeyEcdsa, err := dkmsEcdsa.PbKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}
	core.InitDKMSwithKeyPair(pvKeyEcdsa, pbKeyEcdsa)
	core.InitDKMSwithKeyPair(nil, nil)
	core.InitDKMS("")

	dkms, err := core.InitDKMS(core.KeyTypeRSA)
	if err != nil {
		t.Fatal(err.Error())
	}

	log.Printf("%v", dkmsEcdsa)

	pvKeyRsa, err := dkms.PvKeyRSA()
	if err != nil {
		t.Fatal(err)
	}
	pbKeyRsa, err := dkms.PbKeyRSA()
	if err != nil {
		t.Fatal(err)
	}

	dkms.SetPvKey(pvKeyRsa)
	dkms.SetPvKeyPEM(dkms.PvKeyPEM())
	dkms.SetPvKeyBase58(dkms.PvKeyBase58())
	dkms.SetPbKey(pbKeyRsa)
	dkms.SetPbKeyPEM(dkms.PbKeyPEM())
	dkms.SetPbKeyBase58(dkms.PbKeyBase58())
	dkms.SetDid(dkms.Did())

	core.InitDKMSwithKeyPair(pvKeyRsa, pbKeyRsa)

	err = pvKeyRsa.Validate()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestKeyExport(t *testing.T) {
	myDkms, err := core.InitDKMS(core.KeyTypeRSA)
	if err != nil {
		t.Fatal(err)
	}
	pvKey, err := myDkms.PvKeyRSA()
	if err != nil {
		t.Fatal(err)
	}
	pvKeyPem := core.ExportPrivateKeyAsPEM(pvKey)
	if pvKeyPem != myDkms.PvKeyPEM() {
		t.Fatal(errors.New("export result error"))
	}

	block, _ := pem.Decode([]byte(pvKeyPem))
	if block == nil {
		t.Fatal(errors.New("failed to parse PEM block containing the private key"))
	}

	if block.Type != "RSA PRIVATE KEY" {
		t.Fatal(errors.New("key type mismatched"))
	}
}

func TestEncryptDecrypt(t *testing.T) {
	myDkms := core.GetDKMS()
	plainText := "TestEncryptDecrypt"
	encryptedText := core.PbKeyEncrypt(myDkms.PbKeyBase58(), plainText)
	decryptedText := core.PvKeyDecrypt(encryptedText, myDkms.PvKeyBase58())
	if plainText != decryptedText {
		t.Fatal(errors.New("plainText and decryptedText are not same"))
	}

	encryptedText = myDkms.Encrypt(plainText)
	decryptedText = myDkms.Decrypt(encryptedText)
	if plainText != decryptedText {
		t.Fatal(errors.New("dkms encrypt/decrypt mismatch"))
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

	ret, result = myDkms.Sign(plainText)
	if !ret {
		t.Fatal("dkms sign failed")
	}
	if !myDkms.Verify(plainText, result) {
		t.Fatal(errors.New("dkms verify error"))
	}
}

func TestRandomString(t *testing.T) {
	rndStr1 := core.RandomString(12)
	rndStr2 := core.RandomString(12)
	if rndStr1 == rndStr2 {
		t.Fatal(errors.New("random strings are same"))
	}
}

//func TestCreateResolveDID(t *testing.T) {
//	myDkms := core.GetDKMS()
//	createdDid := core.CreateDID(myDkms.PbKeyBase58(), "byd50")
//	createdDoc := core.ResolveDID(createdDid)
//	t.Logf("createdDid: %v\n createdDoc: %v\n", createdDid, string(createdDoc))
//}

func TestVC(t *testing.T) {
	// ******************** KeyTypeECDSA ******************** //
	issuerDkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}

	// ******************** Create DID ******************** //
	method := "byd50"
	pbKey := issuerDkmsEcdsa.PbKeyBase58()
	issuerDid, _ := dids.CreateDID(method, pbKey)

	issuerDkmsEcdsa.SetDid(issuerDid)
	pvKey, err := issuerDkmsEcdsa.PvKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}

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
	vcJwt := core.CreateVc(issuerDid, typ, credSub, standardClaims, pvKey)
	vcJwt2 := core.CreateVcWithClaims(issuerDid, claims, pvKey)
	t.Logf(" -- vc sample jwt --\n%v\n%v", vcJwt, vcJwt2)

	// ******************** Verify VC ******************** //
}

func TestVP(t *testing.T) {
	// ******************** Sequence 1 preparing for vc ******************** //
	// ******************** KeyTypeECDSA ******************** //
	issuerDkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}
	// ******************** Create DID ******************** //
	method := "byd50"
	pbKey := issuerDkmsEcdsa.PbKeyBase58()
	issuerDid, _ := dids.CreateDID(method, pbKey)
	issuerDkmsEcdsa.SetDid(issuerDid)
	issuerPvKey, err := issuerDkmsEcdsa.PvKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}

	// ******************** Create VC ******************** //
	kid := issuerDid
	vcJwt := byd50_jwt.MakeVcSample(kid, issuerPvKey)
	t.Logf(" -- vc sample jwt --\n%v", vcJwt)

	holderDkmsEcdsa, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err.Error())
	}

	// ******************** Sequence 2 preparing for vp ******************** //
	// ******************** Create DID ******************** //

	holderDid, _ := dids.CreateDID(method, holderDkmsEcdsa.PbKeyBase58())
	holderDkmsEcdsa.SetDid(holderDid)
	holderPvKey, err := holderDkmsEcdsa.PvKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}

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
	nonce := core.RandomString(12)

	// Create the Claims
	claims := byd50_jwt.VpClaims{
		nonce,
		myVp,
		standardClaims,
	}

	vpJwt2 := core.CreateVpWithClaims(holderDid, claims, holderPvKey)
	t.Logf(" -- vp sample jwt --\n%v", vpJwt)
	t.Logf(" -- vp sample jwt2 --\n%v", vpJwt2)

}

func TestDKMSExportsECDSA(t *testing.T) {
	dkms, err := core.InitDKMS(core.KeyTypeECDSA)
	if err != nil {
		t.Fatal(err)
	}
	pvKey, err := dkms.PvKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}
	pbKey, err := dkms.PbKeyECDSA()
	if err != nil {
		t.Fatal(err)
	}
	if core.ExportPrivateKeyAsPEM(pvKey) == "" {
		t.Fatal(errors.New("ecdsa private key pem empty"))
	}
	if core.ExportPublicKeyAsPEM(pbKey) == "" {
		t.Fatal(errors.New("ecdsa public key pem empty"))
	}
	if core.ExportPrivateKeyAsBase58(pvKey) == "" {
		t.Fatal(errors.New("ecdsa private key base58 empty"))
	}
	if core.ExportPublicKeyAsBase58(pbKey) == "" {
		t.Fatal(errors.New("ecdsa public key base58 empty"))
	}

	if core.ExportPrivateKeyAsPEM("bad") != "" {
		t.Fatal(errors.New("unexpected pem for invalid private key"))
	}
	if core.ExportPublicKeyAsPEM("bad") != "" {
		t.Fatal(errors.New("unexpected pem for invalid public key"))
	}
	if core.ExportPrivateKeyAsBase58("bad") != "" {
		t.Fatal(errors.New("unexpected base58 for invalid private key"))
	}
	if core.ExportPublicKeyAsBase58("bad") != "" {
		t.Fatal(errors.New("unexpected base58 for invalid public key"))
	}
}
