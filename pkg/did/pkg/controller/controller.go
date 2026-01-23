package controller

import (
	"byd50-ssi/pkg/did/core"
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/did/core/rc"
	derrors "byd50-ssi/pkg/did/errors"
	pb "byd50-ssi/proto-files"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"
)

/**
 * Create a DID Document.
 *
 * @param publicKey the json string returned by calling
 * @return the Document object
 */
func CreateDID(pbKeyBase58, method string) string {
	did, err := CreateDIDWithErr(pbKeyBase58, method)
	if err != nil {
		log.Printf("CreateDID error: %v", err)
		return ""
	}
	return did
}

func CreateDIDWithErr(pbKeyBase58, method string) (string, error) {
	if pbKeyBase58 == "" {
		return "", derrors.New(derrors.CodeEmptyKey, "public key is empty")
	}
	registrarClient := getRegistrarClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := registrarClient.CreateDid(ctx, &pb.CreateDidRequest{PublicKeyBase58: pbKeyBase58, Method: method})
	if err != nil {
		return "", derrors.Wrap(derrors.CodeUpstream, "registrar create did failed", err)
	}
	if r.GetDid() == "" {
		return "", derrors.New(derrors.CodeInternal, "registrar returned empty did")
	}

	log.Printf("Created DID: %s", r.GetDid())
	return r.GetDid(), nil
}

/**
 * Add a publicKey to DID Document.
 *
 * @param signedJwt the string that signed the object returned by calling
 * @return the Document object
 */
/**
func addPublicKey(pbKey, signedJwt string) string {
	// Add PublicKey in to the Document
	// ToDo..
	document := "document"

	return document
}
*/

/**
 * Revoke a publicKey in the DID Document.
 *
 * @param signedJwt the string that signed the object returned by calling
 * @return the Document object
 */
/**
func revokePublicKey(pbKey, signedJwt string) string {
	// Add PublicKey in to the Document
	document := "document"

	return document
}
*/

/**
 * Get a DID Document.
 *
 * @param did the id of a DID Document
 * @return the Document object
 */
func ResolveDID(dID string) string {
	doc, err := ResolveDIDWithErr(dID)
	if err != nil {
		log.Printf("ResolveDID error: %v", err)
		return ""
	}
	return doc
}

func ResolveDIDWithErr(dID string) (string, error) {
	if dID == "" {
		return "", derrors.New(derrors.CodeInvalidInput, "did is empty")
	}
	// Set up a connection to the server.
	registrarClient := getRegistrarClient()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := registrarClient.ResolveDid(ctx, &pb.ResolveDidRequest{Did: dID})
	if err != nil {
		return "", derrors.Wrap(derrors.CodeUpstream, "registrar resolve did failed", err)
	}
	log.Printf("ResolveDID(%v)", dID)

	var resolveResponse dids.ResolveResponse
	resolveResponse.ResolutionMetadata.ResolutionError = r.GetResolutionError()

	documents := r.GetDidDocument()
	if documents == "" {
		return "", derrors.New(derrors.CodeNotFound, "registrar returned empty document")
	}
	return documents, nil
}

/**
 * Get a publicKey that matches the id of DID document and the id of publicKey.
 *
 * @param did   the id of DID document
 * @param keyId the id of publicKey
 * @return the publicKey object
 */
func GetPublicKey(did, keyId string) string {
	pbKeyBase58, err := GetPublicKeyWithErr(did, keyId)
	if err != nil {
		log.Printf("GetPublicKey error: %v", err)
		return ""
	}
	return pbKeyBase58
}

func GetPublicKeyWithErr(did, keyId string) (string, error) {
	// Add PublicKey in to the Document
	var ifDoc dids.DocumentInterface
	document, err := ResolveDIDWithErr(did)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal([]byte(document), &ifDoc); err != nil {
		return "", derrors.Wrap(derrors.CodeInternal, "failed to parse did document", err)
	}
	if len(ifDoc.Authentication) == 0 {
		return "", derrors.New(derrors.CodeNotFound, "no authentication keys in document")
	}
	pbKeyBase58 := ifDoc.Authentication[0].PublicKeyBase58

	if len(pbKeyBase58) == 0 {
		return "", derrors.New(derrors.CodeEmptyKey, "public key is empty")
	}

	return pbKeyBase58, nil
}

var registrarClientProvider = rc.GetRegistrarClient

func getRegistrarClient() pb.RegistrarClient {
	return registrarClientProvider()
}

func GetAuthChallengeString(did, plainText string) string {
	pbKeyBase58, err := GetPublicKeyWithErr(did, "")
	if err != nil {
		log.Printf("GetAuthChallengeString error: %v", err)
		return ""
	}

	authChallengeString := core.PbKeyEncrypt(pbKeyBase58, plainText)
	log.Printf("  -- plainText: %v", plainText)
	log.Printf("  == authChallengeString(%v)", len(authChallengeString))

	return authChallengeString
}

func GetAuthResponseString(challengeString, pvKeyBase58 string) string {
	authResponseString := core.PvKeyDecrypt(challengeString, pvKeyBase58)
	return authResponseString
}

func GetSimplePresent(did, pvKeyBase58 string) string {
	didAndTime := did + ";" + time.Now().UTC().Format(time.RFC3339)

	log.Printf("didAndTime= %v", didAndTime)

	_, signedStr := core.PvKeySign(pvKeyBase58, didAndTime, "")
	simplePresentString := didAndTime + ";" + signedStr

	return simplePresentString
}

func VerifySimplePresent(simplePresentString string) string {
	var result = "fail"

	slice := strings.Split(simplePresentString, ";")
	aDid := slice[0]
	aTime := slice[1]
	aSign := slice[2]
	presentTime, _ := time.Parse(time.RFC3339, aTime)

	log.Printf("presentTime: %v", presentTime)
	duration := time.Now().UTC().Sub(presentTime)
	log.Printf("duration = %v", duration)

	if duration < time.Second*10 {
		pbKeyBase58, err := GetPublicKeyWithErr(aDid, "")
		if err != nil {
			log.Printf("VerifySimplePresent error: %v", err)
			return result
		}
		if core.PbKeyVerify(pbKeyBase58, aDid+";"+aTime, aSign) {
			result = "success"
		}
	} else {
		result = "time out"
	}

	return result
}
