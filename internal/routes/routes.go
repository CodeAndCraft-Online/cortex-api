package routes

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/routes/auth"
	"github.com/CodeAndCraft-Online/cortex-api/internal/routes/comments"
	"github.com/CodeAndCraft-Online/cortex-api/internal/routes/posts"
	"github.com/CodeAndCraft-Online/cortex-api/internal/routes/subs"
	"github.com/CodeAndCraft-Online/cortex-api/internal/routes/users"
	"github.com/CodeAndCraft-Online/cortex-api/internal/routes/votes"
	"github.com/gin-gonic/gin"
)

func BuildRouteGroups(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Cortex API"})
	})

	api := router.Group("/api")

	auth.RegisterAuthRoutes(api)
	posts.RegisterPostRoutes(api)
	subs.RegisterSubRoutes(api)
	users.RegisterUserRoutes(api)
	comments.RegisterCommentsRoutes(api)
	votes.RegisterVotesRoutes(api)
}
