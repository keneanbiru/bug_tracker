package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bug-tracker/models"
	"bug-tracker/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockAuthUseCase is a mock implementation of the AuthUseCaseInterface
type MockAuthUseCase struct {
	mock.Mock
}

// Ensure MockAuthUseCase implements the interface
var _ usecase.AuthUseCaseInterface = (*MockAuthUseCase)(nil)

func (m *MockAuthUseCase) Login(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error) {
	args := m.Called(ctx, req)
	token, _ := args.Get(0).(string)
	user, _ := args.Get(1).(*models.UserResponse)
	return token, user, args.Error(2)
}

func (m *MockAuthUseCase) Register(ctx context.Context, req models.RegisterRequest) (*models.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockAuthUseCase) GetDevelopers(ctx context.Context) ([]models.UserResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.UserResponse), args.Error(1)
}

func (m *MockAuthUseCase) ValidateToken(tokenString string) (*models.User, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestLogin(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        models.LoginRequest
		mockResponse   func(*MockAuthUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful Login",
			payload: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockResponse: func(m *MockAuthUseCase) {
				userID, err := primitive.ObjectIDFromHex("680f741010571194baf681b1")
				if err != nil {
					t.Fatal(err)
				}
				m.On("Login", mock.Anything, models.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				}).Return("token123", &models.UserResponse{
					ID:    userID,
					Name:  "Test User",
					Email: "test@example.com",
					Role:  "developer",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"token": "token123",
				"user": map[string]interface{}{
					"id":    "680f741010571194baf681b1",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
			},
		},
		{
			name: "Invalid Credentials",
			payload: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockResponse: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, models.LoginRequest{
					Email:    "test@example.com",
					Password: "wrongpassword",
				}).Return("", nil, usecase.ErrInvalidPassword)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid email or password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			mockAuth := new(MockAuthUseCase)
			tt.mockResponse(mockAuth)

			// Create the controller with the mock
			controller := NewAuthController(mockAuth)

			// Register the route
			router.POST("/login", controller.Login)

			// Create the request
			payload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestRegister(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        models.RegisterRequest
		mockResponse   func(*MockAuthUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful Registration",
			payload: models.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "developer",
			},
			mockResponse: func(m *MockAuthUseCase) {
				userID, err := primitive.ObjectIDFromHex("680f741010571194baf681b1")
				if err != nil {
					t.Fatal(err)
				}
				m.On("Register", mock.Anything, models.RegisterRequest{
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password123",
					Role:     "developer",
				}).Return(&models.UserResponse{
					ID:    userID,
					Name:  "Test User",
					Email: "test@example.com",
					Role:  "developer",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"user": map[string]interface{}{
					"id":    "680f741010571194baf681b1",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
			},
		},
		{
			name: "Email Already Exists",
			payload: models.RegisterRequest{
				Name:     "Test User",
				Email:    "existing@example.com",
				Password: "password123",
				Role:     "developer",
			},
			mockResponse: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, models.RegisterRequest{
					Name:     "Test User",
					Email:    "existing@example.com",
					Password: "password123",
					Role:     "developer",
				}).Return(nil, usecase.ErrEmailAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "Email already exists",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			mockAuth := new(MockAuthUseCase)
			tt.mockResponse(mockAuth)

			// Create the controller with the mock
			controller := NewAuthController(mockAuth)

			// Register the route
			router.POST("/register", controller.Register)

			// Create the request
			payload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
