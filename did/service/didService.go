package service

// ToDo.. TBD
// Service API - 상위 서비스 레이어로 REST API 지원하도록 수정 예정

/**
 * Create a DID Document.
 *
 * @param publicKey the json string returned by calling
 * @return the Document object
 */
func createDocument(pbKey string) string {
	// Create Document
	document := "document"

	return document
}

/**
 * Add a publicKey to DID Document.
 *
 * @param signedJwt the string that signed the object returned by calling
 * @return the Document object
 */
func addPublicKey(pbKey, signedJwt string) string {
	// Add PublicKey in to the Document
	document := "document"

	return document
}

/**
 * Revoke a publicKey in the DID Document.
 *
 * @param signedJwt the string that signed the object returned by calling
 * @return the Document object
 */
func revokePublicKey(pbKey, signedJwt string) string {
	// Add PublicKey in to the Document
	document := "document"

	return document
}

/**
 * Get a DID Document.
 *
 * @param did the id of a DID Document
 * @return the Document object
 */
func readDocument(did string) string {
	// Read the Document
	document := "document"

	return document
}

/**
 * Get a publicKey that matches the id of DID document and the id of publicKey.
 *
 * @param did   the id of DID document
 * @param keyId the id of publicKey
 * @return the publicKey object
 */
func getPublicKey(did, keyId string) string {
	// Add PublicKey in to the Document
	document := "document"
	publicKey := document

	return publicKey
}
