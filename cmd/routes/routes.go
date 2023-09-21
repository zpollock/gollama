package routes

import (
	"gollama/cmd/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/chat/completions", controllers.ChatCompletionsHandler)
	r.POST("/completions", controllers.CompletionsHandler)
}
