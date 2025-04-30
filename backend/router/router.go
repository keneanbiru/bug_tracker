package router

import (
	"bug-tracker/controller"
	"bug-tracker/usecase"

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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

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
