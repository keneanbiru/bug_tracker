package controller

import (
	"net/http"

	"bug-tracker/models"
	"bug-tracker/usecase"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase usecase.AuthUseCaseInterface
}

func NewAuthController(authUseCase usecase.AuthUseCaseInterface) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

type RegisterResponse struct {
	User models.UserResponse `json:"user"`
}

type LoginResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authUseCase.Register(ctx, req)
	if err != nil {
		switch err {
		case usecase.ErrEmailAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{User: *user})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := c.authUseCase.Login(ctx, req)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		case usecase.ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}

func (c *AuthController) GetDevelopers(ctx *gin.Context) {
	developers, err := c.authUseCase.GetDevelopers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch developers"})
		return
	}

	ctx.JSON(http.StatusOK, developers)
}
