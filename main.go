package main

import (
	"os"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	routes "github.com/CodeAndCraft-Online/cortex-api/internal/routes"
	pkg "github.com/CodeAndCraft-Online/cortex-api/pkg"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()
	router := gin.Default()

	// ✅ Apply rate limiter: 100 requests per minute per IP
	limiter := pkg.NewRateLimiter(100, time.Minute)
	router.Use(limiter.Middleware())

	routes.BuildRouteGroups(router)
	router.Run(":" + os.Getenv("PORT"))
}
