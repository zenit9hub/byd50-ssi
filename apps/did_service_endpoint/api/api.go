package api

import (
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/pkg/controller"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"time"
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

type CreateVcRequestBody struct {
	Kid               string                 `json:"kid" example:"did:byd50:1234567890abcdef"`
	PvKeyBase58       string                 `json:"pv_key_base58" example:"3VZ6oJdR8i1qKX7kH3Yv9d7w7wzgn..."`
	Type              string                 `json:"type" example:"AlumniCredential"`
	CredentialSubject map[string]interface{} `json:"credential_subject"`
	Issuer            string                 `json:"issuer" example:"http://demo-issuer.example"`
	Subject           string                 `json:"subject" example:"did:byd50:holder123"`
	ExpiresInMinutes  int                    `json:"expires_in_minutes" example:"5"`
}

type CreateVcResponse struct {
	VcJwt string `json:"vc_jwt" example:"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type VerifyVcRequestBody struct {
	VcJwt string `json:"vc_jwt" example:"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type CreateVpRequestBody struct {
	HolderDid        string   `json:"holder_did" example:"did:byd50:holder123"`
	PvKeyBase58      string   `json:"pv_key_base58" example:"3VZ6oJdR8i1qKX7kH3Yv9d7w7wzgn..."`
	Type             string   `json:"type" example:"CredentialManagerPresentation"`
	VcJwts           []string `json:"vc_jwts"`
	Issuer           string   `json:"issuer" example:"client make this vp"`
	Subject          string   `json:"subject" example:"did:byd50:holder123"`
	ExpiresInMinutes int      `json:"expires_in_minutes" example:"5"`
	Audience         string   `json:"aud" example:"did:byd50:rental456"`
	Nonce            string   `json:"nonce" example:"n-123456"`
	SimplePresentation bool   `json:"simple_presentation" example:"false"`
}

type CreateVpResponse struct {
	VpJwt string `json:"vp_jwt" example:"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type VerifyVpRequestBody struct {
	VpJwt string `json:"vp_jwt" example:"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpectedAud   string `json:"expected_aud,omitempty" example:"did:byd50:rental456"`
	ExpectedNonce string `json:"expected_nonce,omitempty" example:"n-123456"`
}

type VerifyResponse struct {
	Valid bool   `json:"valid" example:"true"`
	Error string `json:"error,omitempty" example:"signature invalid"`
}

func standardClaims(issuer, subject string, expiresInMinutes int) jwt.StandardClaims {
	if expiresInMinutes <= 0 {
		expiresInMinutes = 5
	}
	now := time.Now()
	return jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: now.Add(time.Duration(expiresInMinutes) * time.Minute).Unix(),
		Id:        core.RandomString(12),
		IssuedAt:  now.Unix(),
		Issuer:    issuer,
		NotBefore: now.Unix(),
		Subject:   subject,
	}
}

func parseEcPrivateKeyBase58(pvKeyBase58 string) (*ecdsa.PrivateKey, error) {
	keyBytes := base58.Decode(pvKeyBase58)
	return x509.ParseECPrivateKey(keyBytes)
}

func logReq(c *gin.Context, action string, fields map[string]string) {
	log.Printf("[did_service_endpoint][%s] %s %s ip=%s fields=%v",
		action, c.Request.Method, c.Request.URL.Path, c.ClientIP(), fields)
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
		logReq(c, "CreateDid.BadRequest", map[string]string{"error": "invalid json"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid json body",
		})
		return
	}
	pbKeyBase58 := requestBody.PublicKeyBase58
	method := requestBody.Method
	logReq(c, "CreateDid.Request", map[string]string{
		"method": method,
		"hasKey": strconv.FormatBool(pbKeyBase58 != ""),
	})
	if method != "" && pbKeyBase58 != "" {
		did := controller.CreateDID(pbKeyBase58, method)
		if did != "" {
			logReq(c, "CreateDid.Success", map[string]string{"did": did})
			c.JSON(http.StatusOK, CreateDidResponse{Did: did})
			return
		}
		logReq(c, "CreateDid.Error", map[string]string{"error": "create did failed"})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create did",
		})
		return
	}
	logReq(c, "CreateDid.BadRequest", map[string]string{"error": "missing method or public_key_base58"})
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
	logReq(c, "GetDid.Request", map[string]string{"did": did})
	if did == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "did is required",
		})
		return
	}
	document := controller.ResolveDID(did)
	if document != "" {
		logReq(c, "GetDid.Success", map[string]string{"did": did})
		c.JSON(http.StatusOK, GetDidResponse{Document: document})
		return
	}
	logReq(c, "GetDid.NotFound", map[string]string{"did": did})
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
	logReq(c, "GetDidPublicKey.Request", map[string]string{"did": did})
	if did == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "did is required",
		})
		return
	}
	publicKeyBase58 := controller.GetPublicKey(did, "")
	if publicKeyBase58 != "" {
		logReq(c, "GetDidPublicKey.Success", map[string]string{"did": did})
		c.JSON(http.StatusOK, GetDidPublicKeyResponse{PublicKeyBase58: publicKeyBase58})
		return
	}
	logReq(c, "GetDidPublicKey.NotFound", map[string]string{"did": did})
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    "NOT_FOUND",
		Message: "public key not found",
	})
}

// CreateVc
// @Summary Create VC
// @Description Create a Verifiable Credential (JWT) signed by the issuer's private key.
// @ID createVc
// @Accept  json
// @Produce  json
// @Param   CreateVcRequestBody  body    CreateVcRequestBody  true  "Create VC request"
// @Success 200 {object} CreateVcResponse "ok" example({"vc_jwt":"eyJhbGciOi..."} )
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"kid, pv_key_base58, and credential_subject are required"})
// @Failure 500 {object} ErrorResponse "internal error" example({"code":"INTERNAL_ERROR","message":"failed to create vc"})
// @Security ApiKeyAuth
// @Router /testapi/vc/create [post]
func CreateVc(c *gin.Context) {
	var requestBody CreateVcRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logReq(c, "CreateVc.BadRequest", map[string]string{"error": "invalid json"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid json body",
		})
		return
	}
	logReq(c, "CreateVc.Request", map[string]string{
		"kid":     requestBody.Kid,
		"type":    requestBody.Type,
		"subject": requestBody.Subject,
	})
	if requestBody.Kid == "" || requestBody.PvKeyBase58 == "" || requestBody.CredentialSubject == nil {
		logReq(c, "CreateVc.BadRequest", map[string]string{"error": "missing required fields"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "kid, pv_key_base58, and credential_subject are required",
		})
		return
	}
	pvKey, err := parseEcPrivateKeyBase58(requestBody.PvKeyBase58)
	if err != nil {
		logReq(c, "CreateVc.BadRequest", map[string]string{"error": "invalid pv_key_base58"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid pv_key_base58",
		})
		return
	}
	typArray := []string{"VerifiableCredential"}
	if requestBody.Type != "" {
		typArray = append(typArray, requestBody.Type)
	}
	claims := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":              typArray,
		"credentialSubject": requestBody.CredentialSubject,
	}
	issuer := requestBody.Issuer
	if issuer == "" {
		issuer = requestBody.Kid
	}
	stdClaims := standardClaims(issuer, requestBody.Subject, requestBody.ExpiresInMinutes)
	vcClaims := byd50_jwt.VcClaims{
		core.RandomString(12),
		claims,
		stdClaims,
	}
	vcJwt := core.CreateVcWithClaims(requestBody.Kid, vcClaims, pvKey)
	if vcJwt == "" {
		logReq(c, "CreateVc.Error", map[string]string{"error": "create vc failed"})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create vc",
		})
		return
	}
	logReq(c, "CreateVc.Success", map[string]string{"vc_len": strconv.Itoa(len(vcJwt))})
	c.JSON(http.StatusOK, CreateVcResponse{VcJwt: vcJwt})
}

// VerifyVc
// @Summary Verify VC
// @Description Verify a VC (JWT) using DID resolver for public key lookup.
// @ID verifyVc
// @Accept  json
// @Produce  json
// @Param   VerifyVcRequestBody  body    VerifyVcRequestBody  true  "Verify VC request"
// @Success 200 {object} VerifyResponse "ok" example({"valid":true})
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"vc_jwt is required"})
// @Security ApiKeyAuth
// @Router /testapi/vc/verify [post]
func VerifyVc(c *gin.Context) {
	var requestBody VerifyVcRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logReq(c, "VerifyVc.BadRequest", map[string]string{"error": "invalid json"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid json body",
		})
		return
	}
	if requestBody.VcJwt == "" {
		logReq(c, "VerifyVc.BadRequest", map[string]string{"error": "vc_jwt is required"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "vc_jwt is required",
		})
		return
	}
	logReq(c, "VerifyVc.Request", map[string]string{"vc_len": strconv.Itoa(len(requestBody.VcJwt))})
	ok, err := core.VerifyVc(requestBody.VcJwt, controller.GetPublicKey)
	if err != nil {
		logReq(c, "VerifyVc.Result", map[string]string{"valid": "false", "error": err.Error()})
		c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: err.Error()})
		return
	}
	logReq(c, "VerifyVc.Result", map[string]string{"valid": strconv.FormatBool(ok)})
	c.JSON(http.StatusOK, VerifyResponse{Valid: ok})
}

// CreateVp
// @Summary Create VP
// @Description Create a Verifiable Presentation (JWT). Supports aud/nonce and simple presentation (no VC).
// @ID createVp
// @Accept  json
// @Produce  json
// @Param   CreateVpRequestBody  body    CreateVpRequestBody  true  "Create VP request"
// @Success 200 {object} CreateVpResponse "ok" example({"vp_jwt":"eyJhbGciOi..."} )
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"holder_did, pv_key_base58, and vc_jwts are required"})
// @Failure 500 {object} ErrorResponse "internal error" example({"code":"INTERNAL_ERROR","message":"failed to create vp"})
// @Security ApiKeyAuth
// @Router /testapi/vp/create [post]
func CreateVp(c *gin.Context) {
	var requestBody CreateVpRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logReq(c, "CreateVp.BadRequest", map[string]string{"error": "invalid json"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid json body",
		})
		return
	}
	if requestBody.HolderDid == "" || requestBody.PvKeyBase58 == "" {
		logReq(c, "CreateVp.BadRequest", map[string]string{"error": "missing required fields"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "holder_did and pv_key_base58 are required",
		})
		return
	}
	if len(requestBody.VcJwts) == 0 && !requestBody.SimplePresentation {
		logReq(c, "CreateVp.BadRequest", map[string]string{"error": "vc_jwts required"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "vc_jwts are required unless simple_presentation is true",
		})
		return
	}
	pvKey, err := parseEcPrivateKeyBase58(requestBody.PvKeyBase58)
	if err != nil {
		logReq(c, "CreateVp.BadRequest", map[string]string{"error": "invalid pv_key_base58"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid pv_key_base58",
		})
		return
	}
	issuer := requestBody.Issuer
	if issuer == "" {
		issuer = requestBody.HolderDid
	}
	subject := requestBody.Subject
	if subject == "" {
		subject = requestBody.HolderDid
	}
	stdClaims := standardClaims(issuer, subject, requestBody.ExpiresInMinutes)
	if requestBody.Audience != "" {
		stdClaims.Audience = requestBody.Audience
	}
	types := []string{"VerifiablePresentation"}
	if requestBody.Type != "" {
		types = append(types, requestBody.Type)
	}
	vpClaims := byd50_jwt.VpClaims{
		requestBody.Nonce,
		map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			"type":                 types,
			"verifiableCredential": requestBody.VcJwts,
		},
		stdClaims,
	}
	vpJwt := core.CreateVpWithClaims(requestBody.HolderDid, vpClaims, pvKey)
	if vpJwt == "" {
		logReq(c, "CreateVp.Error", map[string]string{"error": "create vp failed"})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create vp",
		})
		return
	}
	logReq(c, "CreateVp.Success", map[string]string{"vp_len": strconv.Itoa(len(vpJwt))})
	c.JSON(http.StatusOK, CreateVpResponse{VpJwt: vpJwt})
}

// VerifyVp
// @Summary Verify VP
// @Description Verify a VP (JWT) using DID resolver for public key lookup. Optionally checks aud/nonce.
// @ID verifyVp
// @Accept  json
// @Produce  json
// @Param   VerifyVpRequestBody  body    VerifyVpRequestBody  true  "Verify VP request"
// @Success 200 {object} VerifyResponse "ok" example({"valid":true})
// @Failure 400 {object} ErrorResponse "bad request" example({"code":"INVALID_PARAM","message":"vp_jwt is required"})
// @Security ApiKeyAuth
// @Router /testapi/vp/verify [post]
func VerifyVp(c *gin.Context) {
	var requestBody VerifyVpRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logReq(c, "VerifyVp.BadRequest", map[string]string{"error": "invalid json"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "invalid json body",
		})
		return
	}
	if requestBody.VpJwt == "" {
		logReq(c, "VerifyVp.BadRequest", map[string]string{"error": "vp_jwt is required"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PARAM",
			Message: "vp_jwt is required",
		})
		return
	}
	logReq(c, "VerifyVp.Request", map[string]string{"vp_len": strconv.Itoa(len(requestBody.VpJwt))})
	ok, did, err := core.VerifyVp(requestBody.VpJwt, controller.GetPublicKey)
	if err != nil || !ok {
		if err == nil {
			err = errors.New("vp signature invalid")
		}
		logReq(c, "VerifyVp.Result", map[string]string{"valid": "false", "error": err.Error()})
		c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: err.Error()})
		return
	}
	_, claims, err := core.GetMapClaims(requestBody.VpJwt, controller.GetPublicKey)
	if err == nil && claims != nil {
		if requestBody.ExpectedAud != "" && !matchAudience(claims["aud"], requestBody.ExpectedAud) {
			c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: "audience mismatch"})
			return
		}
		if requestBody.ExpectedNonce != "" {
			nonce, _ := claims["nonce"].(string)
			if nonce != requestBody.ExpectedNonce {
				c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: "nonce mismatch"})
				return
			}
		}
		vcJwts, _ := extractVcJwtsFromVp(claims)
		if len(vcJwts) > 0 {
			vcValid, _ := core.VerifyVc(vcJwts[0], controller.GetPublicKey)
			if !vcValid {
				c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: "vc signature invalid"})
				return
			}
			if vcExpired(vcJwts[0]) {
				c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: "vc expired"})
				return
			}
			holderDid := vcHolderDid(vcJwts[0])
			if holderDid != "" && did != "" && holderDid != did {
				c.JSON(http.StatusOK, VerifyResponse{Valid: false, Error: "vc holder mismatch"})
				return
			}
		}
	}
	logReq(c, "VerifyVp.Result", map[string]string{"valid": "true"})
	c.JSON(http.StatusOK, VerifyResponse{Valid: true})
}
