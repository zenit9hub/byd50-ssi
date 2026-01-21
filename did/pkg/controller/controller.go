package controller

import (
	"byd50-ssi/did/core"
	"byd50-ssi/did/core/dids"
	"byd50-ssi/did/core/rc"
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

	registrarClient := getRegistrarClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := registrarClient.CreateDID(ctx, &pb.CreateDIDsRequest{PublicKeyBase58: pbKeyBase58, Method: method})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Created DID: %s", r.GetDid())
	return r.GetDid()
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
	// Set up a connection to the server.
	registrarClient := getRegistrarClient()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := registrarClient.ResolveDID(ctx, &pb.ResolveDIDsRequest{Did: dID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ResolveDID(%v)", dID)

	var resolveResponse dids.ResolveResponse
	resolveResponse.ResolutionMetadata.ResolutionError = r.GetResolutionError()

	documents := r.GetDidDocument()
	return documents
}

/**
 * Get a publicKey that matches the id of DID document and the id of publicKey.
 *
 * @param did   the id of DID document
 * @param keyId the id of publicKey
 * @return the publicKey object
 */
func GetPublicKey(did, keyId string) string {
	// Add PublicKey in to the Document
	var ifDoc dids.DocumentInterface
	document := ResolveDID(did)
	json.Unmarshal([]byte(document), &ifDoc)
	pbKeyBase58 := ifDoc.Authentication[0].PublicKeyBase58

	if len(pbKeyBase58) == 0 {
		log.Printf("\n\t !!warning!! pbKeyBase58 length is zero!! \n")
	}

	return pbKeyBase58
}

var registrarClientProvider = rc.GetRegistrarClient

func getRegistrarClient() pb.RegistrarClient {
	return registrarClientProvider()
}

func GetAuthChallengeString(did, plainText string) string {
	pbKeyBase58 := GetPublicKey(did, "")

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
		pbKeyBase58 := GetPublicKey(aDid, "")
		if core.PbKeyVerify(pbKeyBase58, aDid+";"+aTime, aSign) {
			result = "success"
		}
	} else {
		result = "time out"
	}

	return result
}
