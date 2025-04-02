package router

import (
	"bug-tracker/controller"
	"bug-tracker/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authController *controller.AuthController
}

func NewRouter(authController *controller.AuthController) *Router {
	return &Router{
		authController: authController,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vue.js dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", r.authController.Register)
		auth.POST("/login", r.authController.Login)
	}

	return router
}

// AuthMiddleware validates the JWT token
func AuthMiddleware(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		user, err := authUseCase.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}
