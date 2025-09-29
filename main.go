//	@title			Cortex API
//	@version		1.0
//	@description	Cortex API is a Reddit-like social media platform backend API built with Go and PostgreSQL.
//	@contact.name	CodeAndCraft Online
//	@contact.url	https://github.com/CodeAndCraft-Online/cortex-api
//	@contact.email	support@cortex-api.com
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//	@host			localhost:8080
//	@BasePath		/api
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
//	@description "JWT Authorization header using the Bearer scheme. Example: \"Authorization: Bearer {token}\""
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

package main

import (
	"os"
	"time"

	_ "github.com/CodeAndCraft-Online/cortex-api/docs"
	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	routes "github.com/CodeAndCraft-Online/cortex-api/internal/routes"
	pkg "github.com/CodeAndCraft-Online/cortex-api/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	db.InitDB()
	router := gin.Default()

	// ✅ Apply rate limiter: 100 requests per minute per IP
	limiter := pkg.NewRateLimiter(100, time.Minute)
	router.Use(limiter.Middleware())

	routes.BuildRouteGroups(router)

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + os.Getenv("PORT"))
}
