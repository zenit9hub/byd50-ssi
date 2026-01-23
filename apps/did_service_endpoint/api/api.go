package api

import (
	"byd50-ssi/pkg/did/pkg/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateDidRequestBody struct {
	Method          string `json:"method" example:"byd50"`
	PublicKeyBase58 string `json:"public_key_base58" example:"3VZ6oJdR8i1qKX7kH3Yv9d7w7wzgn..."`
}

type CreateDidResponse struct {
	Did string `json:"did" example:"did:byd50:1234567890abcdef"`
}

type GetDidResponse struct {
	Document string `json:"document" example:"{\\\"@context\\\":[\\\"https://www.w3.org/ns/did/v1\\\"],\\\"id\\\":\\\"did:byd50:123\\\",\\\"verificationMethod\\\":[...]}"`
}

type GetDidPublicKeyResponse struct {
	PublicKeyBase58 string `json:"publicKeyBase58" example:"3VZ6oJdR8i1qKX7kH3Yv9d7w7wzgn..."`
}

type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PARAM"`
	Message string `json:"message" example:"method and public_key_base58 are required"`
}

// CreateDid
// @Summary Create DID
// @Description Create a DID using method and public key.
// @ID createDid
// @Accept  json
// @Produce  json
// @Param   CreateDidRequestBody  body    CreateDidRequestBody  true  "Create DID request"
// @Success 200 {object} CreateDidResponse "ok" example({"did":"did:byd50:1234567890abcdef"})
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"method and public_key_base58 are required"})
// @Failure 500 {object} ErrorResponse "internal error" example({"code":"INTERNAL_ERROR","message":"failed to create did"})
// @Security ApiKeyAuth
// @Router /testapi/create-did [post]
func CreateDid(c *gin.Context) {
	var requestBody CreateDidRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid json body",
		})
		return
	}
	pbKeyBase58 := requestBody.PublicKeyBase58
	method := requestBody.Method
	if method != "" && pbKeyBase58 != "" {
		did := controller.CreateDID(pbKeyBase58, method)
		if did != "" {
			c.JSON(http.StatusOK, CreateDidResponse{Did: did})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create did",
		})
		return
	}
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "INVALID_PARAM",
		Message: "method and public_key_base58 are required",
	})

}

// GetDid
// @Summary Get DID Document
// @Description Resolve a DID and return its DID Document.
// @ID getDidDocument
// @Accept  json
// @Produce  json
// @Param   some_id     path    string     true  "DID"
// @Success 200 {object} GetDidResponse "ok" example({"document":"{...did document json...}"})
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"did is required"})
// @Failure 404 {object} ErrorResponse "not found" example({"code":"NOT_FOUND","message":"did document not found"})
// @Security ApiKeyAuth
// @Router /testapi/get-did/{some_id} [get]
func GetDid(c *gin.Context) {
	did := c.Params.ByName("some_id")
	if did == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "did is required",
		})
		return
	}
	document := controller.ResolveDID(did)
	if document != "" {
		c.JSON(http.StatusOK, GetDidResponse{Document: document})
		return
	}
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    "NOT_FOUND",
		Message: "did document not found",
	})
}

// GetDidPublicKey
// @Summary Get DID Public Key
// @Description Resolve a DID and return its public key (Base58).
// @ID getDidPublicKey
// @Accept  json
// @Produce  json
// @Param   some_id     path    string     true  "DID"
// @Success 200 {object} GetDidPublicKeyResponse "ok" example({"publicKeyBase58":"3VZ6oJdR8i1qKX7kH3Yv9d7w7wzgn..."})
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"did is required"})
// @Failure 404 {object} ErrorResponse "not found" example({"code":"NOT_FOUND","message":"public key not found"})
// @Security ApiKeyAuth
// @Router /testapi/get-did-public-key/{some_id} [get]
func GetDidPublicKey(c *gin.Context) {
	did := c.Params.ByName("some_id")
	if did == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "did is required",
		})
		return
	}
	publicKeyBase58 := controller.GetPublicKey(did, "")
	if publicKeyBase58 != "" {
		c.JSON(http.StatusOK, GetDidPublicKeyResponse{PublicKeyBase58: publicKeyBase58})
		return
	}
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    "NOT_FOUND",
		Message: "public key not found",
	})
}
