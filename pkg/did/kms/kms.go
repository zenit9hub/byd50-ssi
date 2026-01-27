package kms

import (
	didcore "byd50-ssi/pkg/did/core"
	"byd50-ssi/pkg/keys"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"log"
)

type KMS struct {
	privateKey       interface{}
	publicKey        interface{}
	privateKeyPem    string
	publicKeyPem     string
	privateKeyBase58 string
	publicKeyBase58  string
	did              string
}

type Signer interface {
	Sign(message string) (bool, string)
}

type Verifier interface {
	Verify(message, signature string) bool
}

type Encryptor interface {
	Encrypt(plainText string) string
}

type Decryptor interface {
	Decrypt(ciphertextBase64 string) string
}

const (
	KeyTypeRSA   = "rsa"
	KeyTypeECDSA = "ecdsa"
)

// Deprecated: use PvKeyRSA or PvKeyECDSA for type-safe access.
func (p *KMS) PvKey() interface{} {
	return p.privateKey
}

func (p *KMS) SetPvKey(pvKey interface{}) error {
	p.privateKey = pvKey
	return nil
}

// Deprecated: use PbKeyRSA or PbKeyECDSA for type-safe access.
func (p *KMS) PbKey() interface{} {
	return p.publicKey
}

func (p *KMS) SetPbKey(pbKey interface{}) error {
	p.publicKey = pbKey
	return nil
}

func (p *KMS) PvKeyRSA() (*rsa.PrivateKey, error) {
	pvKey, ok := p.privateKey.(*rsa.PrivateKey)
	if !ok || pvKey == nil {
		return nil, errors.New("private key is not RSA")
	}
	return pvKey, nil
}

func (p *KMS) PbKeyRSA() (*rsa.PublicKey, error) {
	pbKey, ok := p.publicKey.(*rsa.PublicKey)
	if !ok || pbKey == nil {
		return nil, errors.New("public key is not RSA")
	}
	return pbKey, nil
}

func (p *KMS) PvKeyECDSA() (*ecdsa.PrivateKey, error) {
	pvKey, ok := p.privateKey.(*ecdsa.PrivateKey)
	if !ok || pvKey == nil {
		return nil, errors.New("private key is not ECDSA")
	}
	return pvKey, nil
}

func (p *KMS) PbKeyECDSA() (*ecdsa.PublicKey, error) {
	pbKey, ok := p.publicKey.(*ecdsa.PublicKey)
	if !ok || pbKey == nil {
		return nil, errors.New("public key is not ECDSA")
	}
	return pbKey, nil
}

func (p *KMS) PvKeyBase58() string {
	return p.privateKeyBase58
}

func (p *KMS) SetPvKeyBase58(pvKeyBase58 string) error {
	p.privateKeyBase58 = pvKeyBase58
	return nil
}

func (p *KMS) PbKeyBase58() string {
	return p.publicKeyBase58
}

func (p *KMS) SetPbKeyBase58(pbKeyBase58 string) error {
	p.publicKeyBase58 = pbKeyBase58
	return nil
}

func (p *KMS) PvKeyPEM() string {
	return p.privateKeyPem
}

func (p *KMS) SetPvKeyPEM(pvKeyPem string) error {
	p.privateKeyPem = pvKeyPem
	return nil
}

func (p *KMS) PbKeyPEM() string {
	return p.publicKeyPem
}

func (p *KMS) SetPbKeyPEM(pbKeyPem string) error {
	p.publicKeyPem = pbKeyPem
	return nil
}

func (p *KMS) Did() string {
	return p.did
}

func (p *KMS) SetDid(did string) error {
	if did == "" {
		return errors.New("invalid did")
	}
	p.did = did
	return nil
}

func (p *KMS) Sign(message string) (bool, string) {
	return didcore.PvKeySign(p.privateKeyBase58, message, "")
}

func (p *KMS) Verify(message, signature string) bool {
	return didcore.PbKeyVerify(p.publicKeyBase58, message, signature)
}

func (p *KMS) Encrypt(plainText string) string {
	return didcore.PbKeyEncrypt(p.publicKeyBase58, plainText)
}

func (p *KMS) Decrypt(ciphertextBase64 string) string {
	return didcore.PvKeyDecrypt(ciphertextBase64, p.privateKeyBase58)
}

var kms KMS

// GenerateKeyPair Generate key pair
func GenerateKeyPair(keyType string) (interface{}, interface{}) {
	var privateKey, publicKey interface{}
	switch keyType {
	case KeyTypeRSA:
		privateKey, publicKey = keys.GenerateKeyPair(2048)
	case KeyTypeECDSA:
		privateKey, publicKey, _ = keys.GenerateECDSAKeyPair()
	default:
		log.Fatal("unknown keyType")
	}
	return privateKey, publicKey
}

// InitKMS Generate and assign a key pair.
func InitKMS(keyType string) (KMS, error) {
	switch keyType {
	case KeyTypeRSA:
		fallthrough
	case KeyTypeECDSA:
		pvKey, pbKey := GenerateKeyPair(keyType)
		kms.SetPvKey(pvKey)
		kms.SetPbKey(pbKey)
		kms.SetPvKeyPEM(ExportPrivateKeyAsPEM(pvKey))
		kms.SetPbKeyPEM(ExportPublicKeyAsPEM(pbKey))
		kms.SetPvKeyBase58(ExportPrivateKeyAsBase58(pvKey))
		kms.SetPbKeyBase58(ExportPublicKeyAsBase58(pbKey))
	default:
		return kms, errors.New("unknown keyType: " + keyType)
	}
	return kms, nil
}

func InitKMSwithKeyPair(vPvKey interface{}, vPbKey interface{}) error {
	switch keyType := vPvKey.(type) {
	case *rsa.PrivateKey:
		log.Println(keyType.Validate(), "init KMS")
		pvKey := vPvKey.(*rsa.PrivateKey)
		pbKey := vPbKey.(*rsa.PublicKey)
		kms.SetPvKey(pvKey)
		kms.SetPbKey(pbKey)
		kms.SetPvKeyPEM(ExportPrivateKeyAsPEM(pvKey))
		kms.SetPbKeyPEM(ExportPublicKeyAsPEM(pbKey))
		kms.SetPvKeyBase58(ExportPrivateKeyAsBase58(pvKey))
		kms.SetPbKeyBase58(ExportPublicKeyAsBase58(pbKey))
	case *ecdsa.PrivateKey:
		pvKey := vPvKey.(*ecdsa.PrivateKey)
		pbKey := vPbKey.(*ecdsa.PublicKey)
		kms.SetPvKey(pvKey)
		kms.SetPbKey(pbKey)
		kms.SetPvKeyPEM(ExportPrivateKeyAsPEM(pvKey))
		kms.SetPbKeyPEM(ExportPublicKeyAsPEM(pbKey))
		kms.SetPvKeyBase58(ExportPrivateKeyAsBase58(pvKey))
		kms.SetPbKeyBase58(ExportPublicKeyAsBase58(pbKey))
	default:
		log.Println(keyType, "unknown keyType:%v", keyType)
	}
	return nil
}

func GetKMS() KMS {
	return kms
}

// ExportPrivateKeyAsPEM : Exports a PrivateKey as PEM.
func ExportPrivateKeyAsPEM(privateKey interface{}) string {
	switch v := privateKey.(type) {
	case *rsa.PrivateKey:
		return keys.ExportRSAPrivateKeyAsPEM(v)
	case *ecdsa.PrivateKey:
		pemStr, err := keys.ExportECDSAPrivateKeyAsPEM(v)
		if err != nil {
			log.Printf("%v", err.Error())
			return ""
		}
		return pemStr
	default:
		log.Printf("unknown key type")
		return ""
	}
}

// ExportPublicKeyAsPEM : Exports a PublicKey as PEM.
func ExportPublicKeyAsPEM(publicKey interface{}) string {
	switch v := publicKey.(type) {
	case *rsa.PublicKey:
		pemStr, err := keys.ExportRSAPublicKeyAsPEM(v)
		if err != nil {
			log.Printf("%v", err.Error())
			return ""
		}
		return pemStr
	case *ecdsa.PublicKey:
		pemStr, err := keys.ExportECDSAPublicKeyAsPEM(v)
		if err != nil {
			log.Printf("%v", err.Error())
			return ""
		}
		return pemStr
	default:
		log.Printf("unknown key type")
		return ""
	}
}

// ExportPrivateKeyAsBase58 : Exports a PrivateKey as Base58.
func ExportPrivateKeyAsBase58(privateKey interface{}) string {
	switch v := privateKey.(type) {
	case *rsa.PrivateKey:
		return keys.ExportRSAPrivateKeyAsBase58(v)
	case *ecdsa.PrivateKey:
		return keys.ExportECDSAPrivateKeyAsBase58(v)
	default:
		log.Printf("unknown key type")
		return ""
	}
}

// ExportPublicKeyAsBase58 : Exports a PublicKey as Base58.
func ExportPublicKeyAsBase58(publicKey interface{}) string {
	switch v := publicKey.(type) {
	case *rsa.PublicKey:
		return keys.ExportRSAPublicKeyAsBase58(v)
	case *ecdsa.PublicKey:
		return keys.ExportECDSAPublicKeyAsBase58(v)
	default:
		log.Printf("unknown key type")
		return ""
	}
}
