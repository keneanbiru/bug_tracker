package router

import (
	"bug-tracker/controller"
	"bug-tracker/usecase"
	"os"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authController *controller.AuthController
	bugController  *controller.BugController
	authUseCase    usecase.AuthUseCaseInterface
}

func NewRouter(authController *controller.AuthController, bugController *controller.BugController, authUseCase usecase.AuthUseCaseInterface) *Router {
	return &Router{
		authController: authController,
		bugController:  bugController,
		authUseCase:    authUseCase,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Allow requests from these origins
		allowedOrigins := []string{
			"https://bug-tracker-frontend.vercel.app",
			"https://bug-tracker-frontend-kenean-r.vercel.app", // Deployed frontend URL
			"http://localhost:5174",                            // Local development ports
			"http://localhost:5173",
			"http://localhost:3000",
		}

		// Check if the origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// If no allowed origin matches, allow any origin in development
		if !allowed && (os.Getenv("GIN_MODE") != "release") {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", r.authController.Register)
		auth.POST("/login", r.authController.Login)
		auth.GET("/developers", r.authController.GetDevelopers)
	}

	// Bug routes (protected)
	bugs := router.Group("/api/bugs")
	bugs.Use(AuthMiddleware(r.authUseCase))
	{
		bugs.POST("", r.bugController.CreateBug)
		bugs.GET("", r.bugController.GetBugs)
		bugs.GET("/:id", r.bugController.GetBugByID)
		bugs.PUT("/:id", r.bugController.UpdateBug)
		bugs.DELETE("/:id", r.bugController.DeleteBug)
		bugs.PATCH("/:id/status", r.bugController.UpdateBugStatus)
		bugs.POST("/:id/assign", r.bugController.AssignBug)
	}

	return router
}

// AuthMiddleware validates the JWT token
func AuthMiddleware(authUseCase usecase.AuthUseCaseInterface) gin.HandlerFunc {
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
