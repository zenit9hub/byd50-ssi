package dids

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/pkg/logger"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	uuid "github.com/satori/go.uuid"
	"log"
)

type DocumentInterface struct {
	// @context
	Context []string `json:"@context"`

	// Decentralized identifiers (DIDs)
	ID string `json:"id"`

	// A string or a set of strings that conform to the rules in § 3.1 DID Syntax.
	Controller string `json:"controller"`

	// A set of strings that conform to the rules of [RFC3986] for URIs.
	// (FEATURE AT RISK) ISSUE 2: Implementation of alsoKnownAs
	// The DID Working Group is seeking implementer feedback regarding the alsoKnownAs feature.
	// If there is not enough implementer interest in implementing this feature,
	// it will be removed from this specification and placed into the DID Specification Registries [DID-SPEC-REGISTRIES] as an extension.
	AlsoKnownAs []string `json:"alsoKnownAs"`

	// The authentication property is OPTIONAL.
	// If present, the associated value MUST be a set of one or more verification methods.
	Authentication []AuthenticationProperty `json:"authentication"`

	// A set of Service Endpoint maps that conform to the rules in § Service properties.
	Service []ServiceProperty `json:"service"`

	// The verificationMethod property is OPTIONAL.
	// If present, the value MUST be a set of verification methods, where each verification method is expressed using a map.
	VerificationMethod []VerificationMethodProperty `json:"verificationMethod"`
}

/**
 * This corresponds to the authentications property of the DIDs specification.
 * www.w3.org/TR/dids-core/#authentication
 */
type AuthenticationProperty struct {
	ID              string `json:"id"`
	Types           string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

/**
 * A network address, such as an HTTP URL, at which services operate on behalf of a DID subject.
 * www.w3.org/TR/dids-core/
 */
type ServiceProperty struct {
	ID              string `json:"id"`
	Types           string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

/**
 * This corresponds to the verification-methods property of the DIDs specification.
 * www.w3.org/TR/dids-core/#verification-methods
 */
type VerificationMethodProperty struct {
	ID              string `json:"id"`
	Types           string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

func CreateDID(method, pbKey string) (string, []byte) {
	var rule = configs.UseConfig.GenerationRule
	did := generateDID(pbKey, method, rule)
	if did == "" {
		return "", nil
	}
	doc, _ := initDocument(did, pbKey)
	return did, doc
}

func initDocument(did, pbKey string) ([]byte, error) {
	var ifDoc DocumentInterface
	ifDoc.Context = []string{"https://www.w3.org/ns/dids/v1"}
	ifDoc.ID = did
	initialAuthKeyId := did + "#keys-1"
	ifDoc.Authentication = []AuthenticationProperty{
		{
			ID:              initialAuthKeyId,
			Controller:      did,
			PublicKeyBase58: pbKey,
		},
	}
	bytes, err := json.MarshalIndent(ifDoc, "", " ")
	log.Printf("[initDocument] - document.ID(%v)\n", ifDoc.ID)
	return bytes, err
}

func UpdateDocument(did string, document []byte) (string, error) {
	logger.FuncStart()

	var ifDoc DocumentInterface
	err := json.Unmarshal(document, &ifDoc)
	if err != nil {
		return err.Error(), err
	}

	logger.FuncEnd()
	return "", nil
}

func generateDID(pbKey, method, rule string) string {
	//generate 'Random DID' or 'Something based on specific identity rule'
	var didMethodSpecificIdentifier string
	switch rule {
	case "base58":
		hash := sha256.New()
		hash.Write([]byte(pbKey))
		digest := hash.Sum(nil)
		didMethodSpecificIdentifier = base58.Encode(digest)
	case "uuid":
		myuuid := uuid.NewV4().String()
		didMethodSpecificIdentifier = myuuid
	case "hexdigit":
		bytes := make([]byte, 20)
		if _, err := rand.Read(bytes); err != nil {
			return ""
		}
		hexDigitStr := hex.EncodeToString(bytes)

		didMethodSpecificIdentifier = hexDigitStr
	default:
		log.Printf("generate rule error - %v", rule)
		return ""
	}

	did := "did:" + method + ":" + didMethodSpecificIdentifier
	log.Printf("[generateDID] (%v) %s", len(did), did)
	return did
}
