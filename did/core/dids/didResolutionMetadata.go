package dids

type ResolutionErrorCode int

const (
	InvalidDid ResolutionErrorCode = iota
	InvalidDidUrl
	NotFound
	RepresentationNotSupported
)

func (d ResolutionErrorCode) String() string {
	var errorCode = [...]string{
		"invalidDid",
		"invalidDidurl",
		"notFound",
		"representationNotSupported",
	}
	return errorCode[int(d)%len(errorCode)]
}

type contentType struct {
	// @contentType
	contentType []string `json:"@contentType"`
}

type ResolutionMetadata struct {
	// @ResolutionError
	ResolutionError string `json:"@error"`
}

type DocumentMetadata struct {
	// @Deactivated
	Deactivated string `json:"@deactivated"`
}

type ResolveResponse struct {
	// @ResolutionMetadata
	ResolutionMetadata  ResolutionMetadata `json:"@resolutionMetadata"`
	DidDocument         DocumentInterface  `json:"@didDocument"`
	DidDocumentMetadata DocumentMetadata   `json:"@didDocumentMetadata"`
}
