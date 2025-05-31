package routes

import (
	"fmt"
	"net/http"
	"restaurant_reviews/internal/handlers"
	"restaurant_reviews/internal/jwtAuth"

	"github.com/gin-gonic/gin"
)

func AuthMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("AuthMidleware")
		tokenString, err := c.Request.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		jwtAuth.VerifyToken(tokenString.Value)
		c.Next()
	}
}

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	loggedin := router.Group("/")
	loggedin.Use(AuthMidleware())
	{
		loggedin.GET("/user", handlers.GetUserHandler)
		loggedin.POST("/user/feedback", handlers.FeedBackHandler)
		loggedin.DELETE("/user/:id", handlers.DeleteUserHandler)
	}

	router.POST("/user/register", handlers.RegisterHandler)
	router.POST("/user/login", handlers.LoginHandler)

	return router
}
