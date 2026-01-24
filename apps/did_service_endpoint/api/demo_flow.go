package api

import (
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/kms"
	"byd50-ssi/pkg/did/pkg/controller"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"sync"
	"time"
)

type DemoActorsResponse struct {
	LicenseIssuerDid string `json:"license_issuer_did" example:"did:byd50:issuer123"`
	RentalCompanyDid string `json:"rental_company_did" example:"did:byd50:rental456"`
}

type ChallengeResponse struct {
	Aud   string `json:"aud" example:"did:byd50:issuer123"`
	Nonce string `json:"nonce" example:"n-123456"`
}

type IssueLicenseRequestBody struct {
	HolderDid        string `json:"holder_did" example:"did:byd50:holder123"`
	SimpleVpJwt      string `json:"simple_vp_jwt" example:"eyJhbGciOi..."`
	ExpectedAud      string `json:"expected_aud" example:"did:byd50:issuer123"`
	ExpectedNonce    string `json:"expected_nonce" example:"n-123456"`
	ExpiresInMinutes int    `json:"expires_in_minutes" example:"5"`
	ExpiresInSeconds int    `json:"expires_in_seconds" example:"300"`
}

type IssueLicenseResponse struct {
	SimplePresentationValid bool   `json:"simple_presentation_valid" example:"true"`
	VcJwt                   string `json:"vc_jwt,omitempty" example:"eyJhbGciOi..."`
	Error                   string `json:"error,omitempty" example:"invalid presentation"`
}

type IssueRentalRequestBody struct {
	VpJwt            string `json:"vp_jwt" example:"eyJhbGciOi..."`
	ExpectedAud      string `json:"expected_aud" example:"did:byd50:rental456"`
	ExpectedNonce    string `json:"expected_nonce" example:"n-789012"`
	ExpiresInMinutes int    `json:"expires_in_minutes" example:"5"`
	ExpiresInSeconds int    `json:"expires_in_seconds" example:"60"`
}

type IssueRentalResponse struct {
	VpSignatureValid bool   `json:"vp_signature_valid" example:"true"`
	AudNonceValid    bool   `json:"aud_nonce_valid" example:"true"`
	VcValid          bool   `json:"vc_valid" example:"true"`
	VcNotExpired     bool   `json:"vc_not_expired" example:"true"`
	HolderDidMatch   bool   `json:"holder_did_match" example:"true"`
	VcJwt            string `json:"vc_jwt,omitempty" example:"eyJhbGciOi..."`
	Error            string `json:"error,omitempty" example:"vp invalid"`
}

type demoActor struct {
	Did         string
	PvKey       *ecdsa.PrivateKey
	PvKeyBase58 string
	PbKeyBase58 string
}

var demoActors struct {
	once    sync.Once
	license demoActor
	rental  demoActor
}

func ensureDemoActors() {
	demoActors.once.Do(func() {
		demoActors.license = createDemoActor("license-issuer")
		demoActors.rental = createDemoActor("rental-company")
		log.Printf("[did_service_endpoint][demo] license_issuer=%s rental_company=%s",
			demoActors.license.Did, demoActors.rental.Did)
	})
}

func createDemoActor(tag string) demoActor {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey := &pvKey.PublicKey
	pvKeyBase58 := kms.ExportPrivateKeyAsBase58(pvKey)
	pbKeyBase58 := kms.ExportPublicKeyAsBase58(pbKey)
	did := controller.CreateDID(pbKeyBase58, "byd50")
	if did == "" {
		log.Printf("[did_service_endpoint][demo] failed to create did for %s", tag)
	}
	return demoActor{
		Did:         did,
		PvKey:       pvKey,
		PvKeyBase58: pvKeyBase58,
		PbKeyBase58: pbKeyBase58,
	}
}

// GetDemoActors
// @Summary Get demo actor DIDs
// @Description Return demo issuer and rental company DID values.
// @ID getDemoActors
// @Accept  json
// @Produce  json
// @Success 200 {object} DemoActorsResponse "ok"
// @Security ApiKeyAuth
// @Router /testapi/demo/actors [get]
func GetDemoActors(c *gin.Context) {
	ensureDemoActors()
	c.JSON(http.StatusOK, DemoActorsResponse{
		LicenseIssuerDid: demoActors.license.Did,
		RentalCompanyDid: demoActors.rental.Did,
	})
}

// LicenseChallenge
// @Summary Get license issuer challenge
// @Description Return aud/nonce used for DID simple presentation (VP without VC).
// @ID licenseChallenge
// @Accept  json
// @Produce  json
// @Success 200 {object} ChallengeResponse "ok"
// @Security ApiKeyAuth
// @Router /testapi/license/challenge [post]
func LicenseChallenge(c *gin.Context) {
	ensureDemoActors()
	c.JSON(http.StatusOK, ChallengeResponse{
		Aud:   demoActors.license.Did,
		Nonce: core.RandomString(12),
	})
}

// RentalChallenge
// @Summary Get rental company challenge
// @Description Return aud/nonce required for rental contract VP submission.
// @ID rentalChallenge
// @Accept  json
// @Produce  json
// @Success 200 {object} ChallengeResponse "ok"
// @Security ApiKeyAuth
// @Router /testapi/rental/challenge [post]
func RentalChallenge(c *gin.Context) {
	ensureDemoActors()
	c.JSON(http.StatusOK, ChallengeResponse{
		Aud:   demoActors.rental.Did,
		Nonce: core.RandomString(12),
	})
}

// IssueLicense
// @Summary Issue license VC
// @Description Verify simple presentation (VP without VC) and issue license VC.
// @ID issueLicense
// @Accept  json
// @Produce  json
// @Param   IssueLicenseRequestBody  body    IssueLicenseRequestBody  true  "Issue license request"
// @Success 200 {object} IssueLicenseResponse "ok"
// @Failure 400 {object} ErrorResponse "bad request"
// @Security ApiKeyAuth
// @Router /testapi/license/issue [post]
func IssueLicense(c *gin.Context) {
	ensureDemoActors()
	var requestBody IssueLicenseRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_PARAM", Message: "invalid json body"})
		return
	}
	if requestBody.HolderDid == "" || requestBody.SimpleVpJwt == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_PARAM", Message: "holder_did and simple_vp_jwt are required"})
		return
	}

	sigValid, vpDid, audOk, nonceOk, _, err := verifyVpExpectations(
		requestBody.SimpleVpJwt,
		requestBody.ExpectedAud,
		requestBody.ExpectedNonce,
	)
	if err != nil || !sigValid || !audOk || !nonceOk {
		c.JSON(http.StatusOK, IssueLicenseResponse{
			SimplePresentationValid: sigValid && audOk && nonceOk,
			Error:                   "simple presentation invalid",
		})
		return
	}
	if vpDid != "" && vpDid != requestBody.HolderDid {
		c.JSON(http.StatusOK, IssueLicenseResponse{
			SimplePresentationValid: false,
			Error:                   "holder did mismatch",
		})
		return
	}

	subject := map[string]interface{}{
		"holderDid":  requestBody.HolderDid,
		"licenseType": "Type-1",
		"country":    "KR",
	}
	stdClaims := standardClaimsWithSeconds(
		demoActors.license.Did,
		requestBody.HolderDid,
		requestBody.ExpiresInSeconds,
		requestBody.ExpiresInMinutes,
	)
	claims := byd50_jwt.VcClaims{
		core.RandomString(12),
		map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			"type":              []string{"VerifiableCredential", "DriverLicenseCredential"},
			"credentialSubject": subject,
		},
		stdClaims,
	}
	vcJwt := core.CreateVcWithClaims(demoActors.license.Did, claims, demoActors.license.PvKey)
	c.JSON(http.StatusOK, IssueLicenseResponse{
		SimplePresentationValid: true,
		VcJwt:                   vcJwt,
	})
}

// IssueRental
// @Summary Issue rental contract VC
// @Description Verify VP (aud/nonce + license VC) and issue rental contract VC.
// @ID issueRental
// @Accept  json
// @Produce  json
// @Param   IssueRentalRequestBody  body    IssueRentalRequestBody  true  "Issue rental request"
// @Success 200 {object} IssueRentalResponse "ok"
// @Failure 400 {object} ErrorResponse "bad request"
// @Security ApiKeyAuth
// @Router /testapi/rental/issue [post]
func IssueRental(c *gin.Context) {
	ensureDemoActors()
	var requestBody IssueRentalRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_PARAM", Message: "invalid json body"})
		return
	}
	if requestBody.VpJwt == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_PARAM", Message: "vp_jwt is required"})
		return
	}

	sigValid, vpDid, audOk, nonceOk, vcJwts, err := verifyVpExpectations(
		requestBody.VpJwt,
		requestBody.ExpectedAud,
		requestBody.ExpectedNonce,
	)
	if err != nil {
		c.JSON(http.StatusOK, IssueRentalResponse{VpSignatureValid: false, Error: err.Error()})
		return
	}
	vcValid := false
	vcNotExpired := false
	holderMatch := false
	if sigValid && len(vcJwts) > 0 {
		vcValid, _ = core.VerifyVc(vcJwts[0], controller.GetPublicKey)
		vcNotExpired = !vcExpired(vcJwts[0])
		holderDid := vcHolderDid(vcJwts[0])
		holderMatch = holderDid != "" && holderDid == vpDid
	}
	if !sigValid || !audOk || !nonceOk || !vcValid || !vcNotExpired || !holderMatch {
		c.JSON(http.StatusOK, IssueRentalResponse{
			VpSignatureValid: sigValid,
			AudNonceValid:    audOk && nonceOk,
			VcValid:          vcValid,
			VcNotExpired:     vcNotExpired,
			HolderDidMatch:   holderMatch,
			Error:            "vp or vc invalid",
		})
		return
	}
	if vpDid == "" {
		c.JSON(http.StatusOK, IssueRentalResponse{
			VpSignatureValid: sigValid,
			AudNonceValid:    audOk && nonceOk,
			VcValid:          vcValid,
			VcNotExpired:     vcNotExpired,
			Error:            "holder did missing",
		})
		return
	}

	subject := map[string]interface{}{
		"holderDid":   vpDid,
		"agreementId": "rent-" + core.RandomString(8),
		"validDays":   1,
	}
	stdClaims := standardClaimsWithSeconds(
		demoActors.rental.Did,
		vpDid,
		requestBody.ExpiresInSeconds,
		requestBody.ExpiresInMinutes,
	)
	claims := byd50_jwt.VcClaims{
		core.RandomString(12),
		map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			"type":              []string{"VerifiableCredential", "RentalCarAgreementCredential"},
			"credentialSubject": subject,
		},
		stdClaims,
	}
	vcJwt := core.CreateVcWithClaims(demoActors.rental.Did, claims, demoActors.rental.PvKey)
	c.JSON(http.StatusOK, IssueRentalResponse{
		VpSignatureValid: true,
		AudNonceValid:    true,
		VcValid:          true,
		VcNotExpired:     true,
		HolderDidMatch:   true,
		VcJwt:            vcJwt,
	})
}

func verifyVpExpectations(vpJwt, expectedAud, expectedNonce string) (bool, string, bool, bool, []string, error) {
	ok, did, err := core.VerifyVp(vpJwt, controller.GetPublicKey)
	if err != nil || !ok {
		if err == nil {
			err = errors.New("vp signature invalid")
		}
		return false, did, false, false, nil, err
	}
	_, claims, err := core.GetMapClaims(vpJwt, controller.GetPublicKey)
	if err != nil || claims == nil {
		return false, did, false, false, nil, err
	}
	audOk := true
	if expectedAud != "" {
		audOk = matchAudience(claims["aud"], expectedAud)
	}
	nonceOk := true
	if expectedNonce != "" {
		nonceVal, _ := claims["nonce"].(string)
		nonceOk = nonceVal == expectedNonce
	}
	vcJwts, _ := extractVcJwtsFromVp(claims)
	return true, did, audOk, nonceOk, vcJwts, nil
}

func matchAudience(audClaim interface{}, expected string) bool {
	switch v := audClaim.(type) {
	case string:
		return v == expected
	case []interface{}:
		for _, a := range v {
			if s, ok := a.(string); ok && s == expected {
				return true
			}
		}
	}
	return false
}

func extractVcJwtsFromVp(claims jwt.MapClaims) ([]string, error) {
	vp, ok := claims["vp"].(map[string]interface{})
	if !ok {
		return nil, errors.New("vp claim missing")
	}
	raw := vp["verifiableCredential"]
	switch v := raw.(type) {
	case string:
		return []string{v}, nil
	case []interface{}:
		var arr []string
		for _, a := range v {
			if s, ok := a.(string); ok {
				arr = append(arr, s)
			}
		}
		return arr, nil
	case []string:
		return v, nil
	default:
		return nil, errors.New("unsupported vc list")
	}
}

func vcExpired(vcJwt string) bool {
	claims, err := parseVcClaims(vcJwt)
	if err != nil {
		return true
	}
	exp, err := core.ClaimsGetExp(claims)
	if err != nil {
		return true
	}
	return time.Now().After(exp)
}

func parseVcClaims(vcJwt string) (jwt.MapClaims, error) {
	parseToken, err := jwt.Parse(vcJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		did, ok := token.Header["kid"].(string)
		if !ok || did == "" {
			return nil, errors.New("missing kid header")
		}
		pbKeyBase58 := controller.GetPublicKey(did, "")
		pbKeyBytes := base58.Decode(pbKeyBase58)
		if len(pbKeyBytes) == 0 {
			return nil, errors.New("invalid public key base58")
		}
		pbKey, err := x509.ParsePKIXPublicKey(pbKeyBytes)
		if err != nil {
			return nil, err
		}
		return pbKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid vc claims")
}

func standardClaimsWithSeconds(issuer, subject string, expiresInSeconds, expiresInMinutes int) jwt.StandardClaims {
	if expiresInSeconds > 0 {
		now := time.Now()
		return jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: now.Add(time.Duration(expiresInSeconds) * time.Second).Unix(),
			Id:        core.RandomString(12),
			IssuedAt:  now.Unix(),
			Issuer:    issuer,
			NotBefore: now.Unix(),
			Subject:   subject,
		}
	}
	return standardClaims(issuer, subject, expiresInMinutes)
}

func vcHolderDid(vcJwt string) string {
	claims, err := parseVcClaims(vcJwt)
	if err != nil || claims == nil {
		return ""
	}
	vc, ok := claims["vc"].(map[string]interface{})
	if !ok || vc == nil {
		return ""
	}
	subject, ok := vc["credentialSubject"].(map[string]interface{})
	if !ok || subject == nil {
		return ""
	}
	if v, ok := subject["holderDid"].(string); ok {
		return v
	}
	if v, ok := subject["id"].(string); ok {
		return v
	}
	if v, ok := subject["did"].(string); ok {
		return v
	}
	return ""
}
