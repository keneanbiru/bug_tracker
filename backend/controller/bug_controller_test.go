package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"bug-tracker/models"
	"bug-tracker/usecase"
)

// MockBugUseCase is a mock implementation of the BugUseCaseInterface
type MockBugUseCase struct {
	mock.Mock
}

// Ensure MockBugUseCase implements BugUseCaseInterface
var _ usecase.BugUseCaseInterface = (*MockBugUseCase)(nil)

func (m *MockBugUseCase) CreateBug(ctx context.Context, req models.CreateBugRequest, reporterID primitive.ObjectID) (*models.BugResponse, error) {
	args := m.Called(ctx, req, reporterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) GetBugByID(ctx context.Context, id primitive.ObjectID) (*models.BugResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) GetAllBugs(ctx context.Context) ([]*models.BugResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) GetBugsByDeveloper(ctx context.Context, developerID primitive.ObjectID) ([]*models.BugResponse, error) {
	args := m.Called(ctx, developerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) UpdateBugStatus(ctx context.Context, bugID primitive.ObjectID, status string, userID primitive.ObjectID) (*models.BugResponse, error) {
	args := m.Called(ctx, bugID, status, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) AssignBug(ctx context.Context, bugID, developerID primitive.ObjectID) (*models.BugResponse, error) {
	args := m.Called(ctx, bugID, developerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) UpdateBug(ctx context.Context, id primitive.ObjectID, req models.UpdateBugRequest, user *models.User) (*models.BugResponse, error) {
	args := m.Called(ctx, id, req, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BugResponse), args.Error(1)
}

func (m *MockBugUseCase) DeleteBug(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateBug(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		payload        models.CreateBugRequest
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful Bug Creation",
			payload: models.CreateBugRequest{
				Title:       "Test Bug",
				Description: "This is a test bug",
				Priority:    "high",
			},
			mockResponse: func(m *MockBugUseCase) {
				m.On("CreateBug", mock.Anything, models.CreateBugRequest{
					Title:       "Test Bug",
					Description: "This is a test bug",
					Priority:    "high",
				}, fixedUserID).Return(&models.BugResponse{
					ID:          fixedBugID,
					Title:       "Test Bug",
					Description: "This is a test bug",
					Priority:    "high",
					Status:      "open",
					ReportedBy: models.UserResponse{
						ID:    fixedUserID,
						Name:  "Test User",
						Email: "test@example.com",
						Role:  "developer",
					},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          "680f74774848325f4e61925c",
				"title":       "Test Bug",
				"description": "This is a test bug",
				"priority":    "high",
				"status":      "open",
				"reported_by": map[string]interface{}{
					"id":    "680f74774848325f4e61925d",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			// Add middleware to set user in context
			router.Use(func(c *gin.Context) {
				c.Set("user", &models.User{ID: fixedUserID})
				c.Next()
			})

			router.POST("/bugs", bugController.CreateBug)

			// Create a request
			payload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/bugs", bytes.NewBuffer(payload))
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

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}

func TestGetBugByID(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		bugID          primitive.ObjectID
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Successful Bug Retrieval",
			bugID: fixedBugID,
			mockResponse: func(m *MockBugUseCase) {
				m.On("GetBugByID", mock.Anything, fixedBugID).Return(&models.BugResponse{
					ID:          fixedBugID,
					Title:       "Test Bug",
					Description: "This is a test bug",
					Priority:    "high",
					Status:      "open",
					ReportedBy: models.UserResponse{
						ID:    fixedUserID,
						Name:  "Test User",
						Email: "test@example.com",
						Role:  "developer",
					},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          "680f74774848325f4e61925c",
				"title":       "Test Bug",
				"description": "This is a test bug",
				"priority":    "high",
				"status":      "open",
				"reported_by": map[string]interface{}{
					"id":    "680f74774848325f4e61925d",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
		{
			name:  "Bug Not Found",
			bugID: fixedBugID,
			mockResponse: func(m *MockBugUseCase) {
				m.On("GetBugByID", mock.Anything, fixedBugID).Return(nil, usecase.ErrBugNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Bug not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			router.GET("/bugs/:id", bugController.GetBugByID)

			// Create a request
			req, _ := http.NewRequest("GET", "/bugs/"+tt.bugID.Hex(), nil)

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

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}

func TestGetBugs(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID1, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedBugID2, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925e")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		userRole       string
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "Get All Bugs as Manager",
			userRole: "manager",
			mockResponse: func(m *MockBugUseCase) {
				m.On("GetAllBugs", mock.Anything).Return([]*models.BugResponse{
					{
						ID:          fixedBugID1,
						Title:       "Test Bug 1",
						Description: "This is test bug 1",
						Priority:    "high",
						Status:      "open",
						ReportedBy: models.UserResponse{
							ID:    fixedUserID,
							Name:  "Test User",
							Email: "test@example.com",
							Role:  "developer",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					},
					{
						ID:          fixedBugID2,
						Title:       "Test Bug 2",
						Description: "This is test bug 2",
						Priority:    "medium",
						Status:      "in_progress",
						ReportedBy: models.UserResponse{
							ID:    fixedUserID,
							Name:  "Test User",
							Email: "test@example.com",
							Role:  "developer",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":          "680f74774848325f4e61925c",
					"title":       "Test Bug 1",
					"description": "This is test bug 1",
					"priority":    "high",
					"status":      "open",
					"reported_by": map[string]interface{}{
						"id":    "680f74774848325f4e61925e",
						"name":  "Test User",
						"email": "test@example.com",
						"role":  "developer",
					},
					"created_at": "0001-01-01T00:00:00Z",
					"updated_at": "0001-01-01T00:00:00Z",
				},
				map[string]interface{}{
					"id":          "680f74774848325f4e61925d",
					"title":       "Test Bug 2",
					"description": "This is test bug 2",
					"priority":    "medium",
					"status":      "in_progress",
					"reported_by": map[string]interface{}{
						"id":    "680f74774848325f4e61925e",
						"name":  "Test User",
						"email": "test@example.com",
						"role":  "developer",
					},
					"created_at": "0001-01-01T00:00:00Z",
					"updated_at": "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:     "Get Developer's Bugs",
			userRole: "developer",
			mockResponse: func(m *MockBugUseCase) {
				m.On("GetBugsByDeveloper", mock.Anything, fixedUserID).Return([]*models.BugResponse{
					{
						ID:          fixedBugID1,
						Title:       "Test Bug 1",
						Description: "This is test bug 1",
						Priority:    "high",
						Status:      "open",
						ReportedBy: models.UserResponse{
							ID:    fixedUserID,
							Name:  "Test User",
							Email: "test@example.com",
							Role:  "developer",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":          "680f74774848325f4e61925c",
					"title":       "Test Bug 1",
					"description": "This is test bug 1",
					"priority":    "high",
					"status":      "open",
					"reported_by": map[string]interface{}{
						"id":    "680f74774848325f4e61925e",
						"name":  "Test User",
						"email": "test@example.com",
						"role":  "developer",
					},
					"created_at": "0001-01-01T00:00:00Z",
					"updated_at": "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:     "Error Getting Bugs",
			userRole: "manager",
			mockResponse: func(m *MockBugUseCase) {
				m.On("GetAllBugs", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to fetch bugs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			// Add middleware to set user in context
			router.Use(func(c *gin.Context) {
				c.Set("user", &models.User{
					ID:   fixedUserID,
					Role: tt.userRole,
				})
				c.Next()
			})

			router.GET("/bugs", bugController.GetBugs)

			// Create a request
			req, _ := http.NewRequest("GET", "/bugs", nil)

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}

func TestUpdateBugStatus(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		bugID          primitive.ObjectID
		payload        models.UpdateBugStatusRequest
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Successful Status Update",
			bugID: fixedBugID,
			payload: models.UpdateBugStatusRequest{
				Status: "in-progress",
			},
			mockResponse: func(m *MockBugUseCase) {
				m.On("UpdateBugStatus", mock.Anything, fixedBugID, "in-progress", fixedUserID).Return(&models.BugResponse{
					ID:          fixedBugID,
					Title:       "Test Bug",
					Description: "This is a test bug",
					Priority:    "high",
					Status:      "in-progress",
					ReportedBy: models.UserResponse{
						ID:    fixedUserID,
						Name:  "Test User",
						Email: "test@example.com",
						Role:  "developer",
					},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          "680f74774848325f4e61925c",
				"title":       "Test Bug",
				"description": "This is a test bug",
				"priority":    "high",
				"status":      "in-progress",
				"reported_by": map[string]interface{}{
					"id":    "680f74774848325f4e61925d",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
		{
			name:  "Bug Not Found",
			bugID: fixedBugID,
			payload: models.UpdateBugStatusRequest{
				Status: "in-progress",
			},
			mockResponse: func(m *MockBugUseCase) {
				m.On("UpdateBugStatus", mock.Anything, fixedBugID, "in-progress", fixedUserID).Return(nil, usecase.ErrBugNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Bug not found",
			},
		},
		{
			name:  "Unauthorized Update",
			bugID: fixedBugID,
			payload: models.UpdateBugStatusRequest{
				Status: "in-progress",
			},
			mockResponse: func(m *MockBugUseCase) {
				m.On("UpdateBugStatus", mock.Anything, fixedBugID, "in-progress", fixedUserID).Return(nil, usecase.ErrUnauthorized)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"error": "Not authorized to update this bug",
			},
		},
		{
			name:  "Invalid Bug ID",
			bugID: primitive.ObjectID{},
			payload: models.UpdateBugStatusRequest{
				Status: "in-progress",
			},
			mockResponse: func(m *MockBugUseCase) {
				// No mock needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid bug ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			// Add middleware to set user in context
			router.Use(func(c *gin.Context) {
				c.Set("user", &models.User{ID: fixedUserID})
				c.Next()
			})

			router.PATCH("/bugs/:id/status", bugController.UpdateBugStatus)

			// Create a request
			payload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("PATCH", "/bugs/"+tt.bugID.Hex()+"/status", bytes.NewBuffer(payload))
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

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}

func TestAssignBug(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}
	fixedDeveloperID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925e")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		bugID          primitive.ObjectID
		developerID    primitive.ObjectID
		userRole       string
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Successful Bug Assignment",
			bugID:       fixedBugID,
			developerID: fixedDeveloperID,
			userRole:    "manager",
			mockResponse: func(m *MockBugUseCase) {
				m.On("AssignBug", mock.Anything, fixedBugID, fixedDeveloperID).Return(&models.BugResponse{
					ID:          fixedBugID,
					Title:       "Test Bug",
					Description: "This is a test bug",
					Priority:    "high",
					Status:      "open",
					ReportedBy: models.UserResponse{
						ID:    fixedUserID,
						Name:  "Test User",
						Email: "test@example.com",
						Role:  "developer",
					},
					AssignedTo: &models.UserResponse{
						ID:    fixedDeveloperID,
						Name:  "Test Developer",
						Email: "developer@example.com",
						Role:  "developer",
					},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          "680f74774848325f4e61925c",
				"title":       "Test Bug",
				"description": "This is a test bug",
				"priority":    "high",
				"status":      "open",
				"reported_by": map[string]interface{}{
					"id":    "680f74774848325f4e61925d",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
				"assigned_to": map[string]interface{}{
					"id":    "680f74774848325f4e61925e",
					"name":  "Test Developer",
					"email": "developer@example.com",
					"role":  "developer",
				},
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
		{
			name:        "Unauthorized Assignment",
			bugID:       fixedBugID,
			developerID: fixedDeveloperID,
			userRole:    "developer",
			mockResponse: func(m *MockBugUseCase) {
				// No mock needed for this case
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"error": "Only managers and admins can assign bugs",
			},
		},
		{
			name:        "Bug Not Found",
			bugID:       fixedBugID,
			developerID: fixedDeveloperID,
			userRole:    "manager",
			mockResponse: func(m *MockBugUseCase) {
				m.On("AssignBug", mock.Anything, fixedBugID, fixedDeveloperID).Return(nil, usecase.ErrBugNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Bug not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			// Add middleware to set user in context
			router.Use(func(c *gin.Context) {
				c.Set("user", &models.User{
					ID:   fixedUserID,
					Role: tt.userRole,
				})
				c.Next()
			})

			router.POST("/bugs/:id/assign", bugController.AssignBug)

			// Create a request
			payload, _ := json.Marshal(map[string]string{
				"developer_id": tt.developerID.Hex(),
			})
			req, _ := http.NewRequest("POST", "/bugs/"+tt.bugID.Hex()+"/assign", bytes.NewBuffer(payload))
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

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}

func TestUpdateBug(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		bugID          primitive.ObjectID
		payload        models.UpdateBugRequest
		userRole       string
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "Successful Bug Update by Reporter",
			bugID: fixedBugID,
			payload: models.UpdateBugRequest{
				Title:       "Updated Bug",
				Description: "This is an updated bug",
				Priority:    "medium",
			},
			userRole: "developer",
			mockResponse: func(m *MockBugUseCase) {
				m.On("UpdateBug", mock.Anything, fixedBugID, models.UpdateBugRequest{
					Title:       "Updated Bug",
					Description: "This is an updated bug",
					Priority:    "medium",
				}, mock.Anything).Return(&models.BugResponse{
					ID:          fixedBugID,
					Title:       "Updated Bug",
					Description: "This is an updated bug",
					Priority:    "medium",
					Status:      "open",
					ReportedBy: models.UserResponse{
						ID:    fixedUserID,
						Name:  "Test User",
						Email: "test@example.com",
						Role:  "developer",
					},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          "680f74774848325f4e61925c",
				"title":       "Updated Bug",
				"description": "This is an updated bug",
				"priority":    "medium",
				"status":      "open",
				"reported_by": map[string]interface{}{
					"id":    "680f74774848325f4e61925d",
					"name":  "Test User",
					"email": "test@example.com",
					"role":  "developer",
				},
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
		{
			name:  "Unauthorized Update",
			bugID: fixedBugID,
			payload: models.UpdateBugRequest{
				Title: "Unauthorized Update",
			},
			userRole: "developer",
			mockResponse: func(m *MockBugUseCase) {
				m.On("UpdateBug", mock.Anything, fixedBugID, models.UpdateBugRequest{
					Title: "Unauthorized Update",
				}, mock.Anything).Return(nil, usecase.ErrUnauthorized)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"error": "Not authorized to update this bug",
			},
		},
		{
			name:  "Bug Not Found",
			bugID: fixedBugID,
			payload: models.UpdateBugRequest{
				Title: "Non-existent Bug Update",
			},
			userRole: "developer",
			mockResponse: func(m *MockBugUseCase) {
				m.On("UpdateBug", mock.Anything, fixedBugID, models.UpdateBugRequest{
					Title: "Non-existent Bug Update",
				}, mock.Anything).Return(nil, usecase.ErrBugNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Bug not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			// Add middleware to set user in context
			router.Use(func(c *gin.Context) {
				c.Set("user", &models.User{
					ID:   fixedUserID,
					Role: tt.userRole,
				})
				c.Next()
			})

			router.PUT("/bugs/:id", bugController.UpdateBug)

			// Create a request
			payload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("PUT", "/bugs/"+tt.bugID.Hex(), bytes.NewBuffer(payload))
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

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}

func TestDeleteBug(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create fixed ObjectIDs for testing
	fixedBugID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925c")
	if err != nil {
		t.Fatal(err)
	}
	fixedUserID, err := primitive.ObjectIDFromHex("680f74774848325f4e61925d")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		bugID          primitive.ObjectID
		userRole       string
		mockResponse   func(*MockBugUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "Successful Bug Deletion",
			bugID:    fixedBugID,
			userRole: "manager",
			mockResponse: func(m *MockBugUseCase) {
				m.On("DeleteBug", mock.Anything, fixedBugID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Bug deleted successfully",
			},
		},
		{
			name:     "Unauthorized Deletion",
			bugID:    fixedBugID,
			userRole: "developer",
			mockResponse: func(m *MockBugUseCase) {
				// No mock needed for this case
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"error": "Only managers and admins can delete bugs",
			},
		},
		{
			name:     "Bug Not Found",
			bugID:    fixedBugID,
			userRole: "manager",
			mockResponse: func(m *MockBugUseCase) {
				m.On("DeleteBug", mock.Anything, fixedBugID).Return(usecase.ErrBugNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Bug not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock
			mockBugUseCase := new(MockBugUseCase)
			tt.mockResponse(mockBugUseCase)

			// Create a new controller with the mock
			bugController := NewBugController(mockBugUseCase)

			// Create a new Gin router
			router := gin.New()

			// Add middleware to set user in context
			router.Use(func(c *gin.Context) {
				c.Set("user", &models.User{
					ID:   fixedUserID,
					Role: tt.userRole,
				})
				c.Next()
			})

			router.DELETE("/bugs/:id", bugController.DeleteBug)

			// Create a request
			req, _ := http.NewRequest("DELETE", "/bugs/"+tt.bugID.Hex(), nil)

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

			// Verify that the mock was called as expected
			mockBugUseCase.AssertExpectations(t)
		})
	}
}
