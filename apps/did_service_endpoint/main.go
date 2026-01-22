package main

import (
	"byd50-ssi/apps/did_service_endpoint/api"
	_ "byd50-ssi/apps/did_service_endpoint/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title DID Phase 2 Test API
// @version 1.0
// @description This is a sample server for DID ServiceEndpoint.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v2
func main() {
	r := gin.New()

	r.POST("/v2/testapi/create-did/", api.CreateDid)
	r.GET("/v2/testapi/get-did/:some_id", api.GetDid)
	r.GET("/v2/testapi/get-did-public-key/:some_id", api.GetDidPublicKey)

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run()
}
