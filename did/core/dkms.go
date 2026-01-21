package core

import (
	"byd50-ssi/did/utility"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"log"
)

type DKMS struct {
	privateKey       interface{}
	publicKey        interface{}
	privateKeyPem    string
	publicKeyPem     string
	privateKeyBase58 string
	publicKeyBase58  string
	did              string
}

const (
	KeyTypeRSA   = "rsa"
	KeyTypeECDSA = "ecdsa"
)

func (p *DKMS) PvKey() interface{} {
	return p.privateKey
}

func (p *DKMS) SetPvKey(pvKey interface{}) error {
	p.privateKey = pvKey
	return nil
}

func (p *DKMS) PbKey() interface{} {
	return p.publicKey
}

func (p *DKMS) SetPbKey(pbKey interface{}) error {
	p.publicKey = pbKey
	return nil
}

func (p *DKMS) PvKeyBase58() string {
	return p.privateKeyBase58
}

func (p *DKMS) SetPvKeyBase58(pvKeyBase58 string) error {
	p.privateKeyBase58 = pvKeyBase58
	return nil
}

func (p *DKMS) PbKeyBase58() string {
	return p.publicKeyBase58
}

func (p *DKMS) SetPbKeyBase58(pbKeyBase58 string) error {
	p.publicKeyBase58 = pbKeyBase58
	return nil
}

func (p *DKMS) PvKeyPEM() string {
	return p.privateKeyPem
}

func (p *DKMS) SetPvKeyPEM(pvKeyPem string) error {
	p.privateKeyPem = pvKeyPem
	return nil
}

func (p *DKMS) PbKeyPEM() string {
	return p.publicKeyPem
}

func (p *DKMS) SetPbKeyPEM(pbKeyPem string) error {
	p.publicKeyPem = pbKeyPem
	return nil
}

func (p *DKMS) Did() string {
	return p.did
}

func (p *DKMS) SetDid(did string) error {
	if did == "" {
		return errors.New("invalid did")
	}
	p.did = did
	return nil
}

var dkms DKMS

// GenerateKeyPair Generate key pair
func GenerateKeyPair(keyType string) (interface{}, interface{}) {
	var privateKey, publicKey interface{}
	switch keyType {
	case KeyTypeRSA:
		privateKey, publicKey = utility.GenerateKeyPair(2048)
	case KeyTypeECDSA:
		privateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		publicKey = &privateKey.(*ecdsa.PrivateKey).PublicKey
	default:
		log.Fatal("unknown keyType")
	}
	return privateKey, publicKey
}

// InitDKMS Generate and assign a key pair.
func InitDKMS(keyType string) (DKMS, error) {
	switch keyType {
	case KeyTypeRSA:
		fallthrough
	case KeyTypeECDSA:
		pvKey, pbKey := GenerateKeyPair(keyType)
		dkms.SetPvKey(pvKey)
		dkms.SetPbKey(pbKey)
		dkms.SetPvKeyPEM(ExportPrivateKeyAsPEM(pvKey))
		dkms.SetPbKeyPEM(ExportPublicKeyAsPEM(pbKey))
		dkms.SetPvKeyBase58(ExportPrivateKeyAsBase58(pvKey))
		dkms.SetPbKeyBase58(ExportPublicKeyAsBase58(pbKey))
	default:
		return dkms, errors.New("unknown keyType: " + keyType)
	}
	return dkms, nil
}

func InitDKMSwithKeyPair(vPvKey interface{}, vPbKey interface{}) error {
	switch keyType := vPvKey.(type) {
	case *rsa.PrivateKey:
		log.Println(keyType.Validate(), "init DKMS")
		pvKey := vPvKey.(*rsa.PrivateKey)
		pbKey := vPbKey.(*rsa.PublicKey)
		dkms.SetPvKey(pvKey)
		dkms.SetPbKey(pbKey)
		dkms.SetPvKeyPEM(ExportPrivateKeyAsPEM(pvKey))
		dkms.SetPbKeyPEM(ExportPublicKeyAsPEM(pbKey))
		dkms.SetPvKeyBase58(ExportPrivateKeyAsBase58(pvKey))
		dkms.SetPbKeyBase58(ExportPublicKeyAsBase58(pbKey))
	case *ecdsa.PrivateKey:
		pvKey := vPvKey.(*ecdsa.PrivateKey)
		pbKey := vPbKey.(*ecdsa.PublicKey)
		dkms.SetPvKey(pvKey)
		dkms.SetPbKey(pbKey)
		dkms.SetPvKeyPEM(ExportPrivateKeyAsPEM(pvKey))
		dkms.SetPbKeyPEM(ExportPublicKeyAsPEM(pbKey))
		dkms.SetPvKeyBase58(ExportPrivateKeyAsBase58(pvKey))
		dkms.SetPbKeyBase58(ExportPublicKeyAsBase58(pbKey))
	default:
		log.Println(keyType, "unknown keyType:%v", keyType)
	}
	return nil
}

func GetDKMS() DKMS {
	return dkms
}

// ExportPrivateKeyAsPEM : Exports a PrivateKey as PEM.
func ExportPrivateKeyAsPEM(privateKey interface{}) string {
	var privateKeyPEM []byte

	switch v := privateKey.(type) {
	case *rsa.PrivateKey:
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(v)
		privateKeyPEM = pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: privateKeyBytes,
			},
		)
	case *ecdsa.PrivateKey:
		privateKeyBytes, err := x509.MarshalECPrivateKey(v)
		if err != nil {
			log.Printf("%v", err.Error())
		} else {
			privateKeyPEM = pem.EncodeToMemory(
				&pem.Block{
					Type:  "ECDSA PRIVATE KEY",
					Bytes: privateKeyBytes,
				},
			)
		}
	default:
		log.Printf("unknown key type")
	}

	return string(privateKeyPEM)
}

// ExportPublicKeyAsPEM : Exports a PublicKey as PEM.
func ExportPublicKeyAsPEM(publicKey interface{}) string {
	var publicKeyPEM []byte
	switch v := publicKey.(type) {
	case *rsa.PublicKey:
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(v)
		if err != nil {
			log.Printf("%v", err.Error())
			return ""
		}
		publicKeyPEM = pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: publicKeyBytes,
			},
		)
	case *ecdsa.PublicKey:
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(v)
		if err != nil {
			log.Printf("%v", err.Error())
		} else {
			publicKeyPEM = pem.EncodeToMemory(
				&pem.Block{
					Type:  "ECDSA PUBLIC KEY",
					Bytes: publicKeyBytes,
				},
			)
		}
	default:
		log.Printf("unknown key type")
	}

	return string(publicKeyPEM)
}

// ExportPrivateKeyAsBase58 : Exports a PrivateKey as Base58.
func ExportPrivateKeyAsBase58(privateKey interface{}) string {
	var privateKeyBase58 string
	switch privateKey.(type) {
	case *rsa.PrivateKey:
		privateKeyBase58 = base58.Encode(x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey)))
	case *ecdsa.PrivateKey:
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey.(*ecdsa.PrivateKey))
		if err != nil {
			log.Printf("error occured: %v", err.Error())
		} else {
			privateKeyBase58 = base58.Encode(privateKeyBytes)
		}
	default:
		log.Printf("unknown key type")
	}

	return privateKeyBase58
}

// ExportPublicKeyAsBase58 : Exports a PublicKey as Base58.
func ExportPublicKeyAsBase58(publicKey interface{}) string {
	var publicKeyBase58 string
	switch publicKey.(type) {
	case *rsa.PublicKey:
		publicKeyBase58 = base58.Encode(x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey)))
	case *ecdsa.PublicKey:
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			log.Printf("error occured: %v", err.Error())
		} else {
			publicKeyBase58 = base58.Encode(publicKeyBytes)
		}
	default:
		log.Printf("unknown key type")
	}

	return publicKeyBase58
}
