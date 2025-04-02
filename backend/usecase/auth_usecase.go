package usecase

import (
	"context"
	"errors"
	"time"

	"bug-tracker/models"
	"bug-tracker/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthUseCase struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthUseCase(userRepo *repository.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req models.RegisterRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create new user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Return user response
	response := user.ToResponse()
	return &response, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrUserNotFound
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return "", nil, ErrInvalidPassword
	}

	// Generate JWT token
	token, err := uc.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	// Return token and user response
	response := user.ToResponse()
	return token, &response, nil
}

func (uc *AuthUseCase) generateToken(user *models.User) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*models.User, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get user ID
	userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	// Get user from repository
	user, err := uc.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
