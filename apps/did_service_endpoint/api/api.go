package api

import (
	"byd50-ssi/pkg/did/pkg/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/web"
	"net/http"
)

type CreateDidRequestBody struct {
	Method          string `form:"method"`
	PublicKeyBase58 string `form:"public_key_base58"`
}

// CreateDid
// @Description create DID
// @Accept  json
// @Produce  json
// @Param   CreateDidRequestBody     body    CreateDidRequestBody     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} web.APIError "We need Method and PublicKey!!"
// @Failure 404 {object} web.APIError "Can not find ID"
// @Router /testapi/create-did/ [post]
func CreateDid(c *gin.Context) {
	err := web.APIError{}
	fmt.Println(err)

	var requestBody CreateDidRequestBody
	if err := c.BindJSON(&requestBody); err != nil {

	}
	pbKeyBase58 := requestBody.PublicKeyBase58
	method := requestBody.Method
	if method != "" && pbKeyBase58 != "" {
		did := controller.CreateDID(pbKeyBase58, method)
		if did != "" {
			c.JSON(http.StatusOK, gin.H{
				"did": did,
			})
		} else {
			c.JSON(http.StatusForbidden, "error occurred")
		}
	} else {
		c.JSON(http.StatusBadRequest, "we need method and public_key_base58")
	}

}

// GetDid
// @Description get DID's Document
// @Accept  json
// @Produce  json
// @Param   some_id     path    string     true        "DID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} web.APIError "We need ID!!"
// @Failure 404 {object} web.APIError "Can not find ID"
// @Router /testapi/get-did/{some_id} [get]
func GetDid(c *gin.Context) {
	did := c.Params.ByName("some_id")
	document := controller.ResolveDID(did)
	if document != "" {
		c.JSON(http.StatusOK, gin.H{
			"document": document,
		})
	} else {
		c.JSON(http.StatusForbidden, "error occurred")
	}
}

// GetDidPublicKey
// @Description get DID's public key
// @Accept  json
// @Produce  json
// @Param   some_id     path    string     true        "DID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} web.APIError "We need ID!!"
// @Failure 404 {object} web.APIError "Can not find ID"
// @Router /testapi/get-did-public-key/{some_id} [get]
func GetDidPublicKey(c *gin.Context) {
	did := c.Params.ByName("some_id")
	publicKeyBase58 := controller.GetPublicKey(did, "")
	if publicKeyBase58 != "" {
		c.JSON(http.StatusOK, gin.H{
			"publicKeyBase58": publicKeyBase58,
		})
	} else {
		c.JSON(http.StatusForbidden, "error occurred")
	}
}
