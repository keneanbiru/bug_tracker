package usecase

import (
	"bug-tracker/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) FindByRole(ctx context.Context, role string) ([]*models.User, error) {
	var developers []*models.User
	for _, user := range m.users {
		if user.Role == role {
			developers = append(developers, user)
		}
	}
	if len(developers) == 0 {
		return nil, errors.New("no developers found")
	}
	return developers, nil
}

func TestRegister(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authUseCase := NewAuthUseCase(mockRepo, "test-secret")

	t.Run("successful registration", func(t *testing.T) {
		req := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
			Role:     "developer",
		}

		response, err := authUseCase.Register(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Email, response.Email)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Role, response.Role)
	})

	t.Run("duplicate email", func(t *testing.T) {
		req := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
			Role:     "developer",
		}

		response, err := authUseCase.Register(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, ErrEmailAlreadyExists, err)
		assert.Nil(t, response)
	})
}

func TestLogin(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authUseCase := NewAuthUseCase(mockRepo, "test-secret")

	// Register a test user first
	req := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "developer",
	}
	_, _ = authUseCase.Register(context.Background(), req)

	t.Run("successful login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		token, response, err := authUseCase.Login(context.Background(), loginReq)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, response)
		assert.Equal(t, req.Email, response.Email)
	})

	t.Run("user not found", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		token, response, err := authUseCase.Login(context.Background(), loginReq)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		assert.Empty(t, token)
		assert.Nil(t, response)
	})

	t.Run("invalid password", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		token, response, err := authUseCase.Login(context.Background(), loginReq)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPassword, err)
		assert.Empty(t, token)
		assert.Nil(t, response)
	})
}

func TestValidateToken(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authUseCase := NewAuthUseCase(mockRepo, "test-secret")

	// Register and login a test user to get a valid token
	req := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "developer",
	}
	_, _ = authUseCase.Register(context.Background(), req)

	loginReq := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	token, _, _ := authUseCase.Login(context.Background(), loginReq)

	t.Run("valid token", func(t *testing.T) {
		user, err := authUseCase.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("invalid token", func(t *testing.T) {
		user, err := authUseCase.ValidateToken("invalid.token.here")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestGetDevelopers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authUseCase := NewAuthUseCase(mockRepo, "test-secret")

	// Register some test developers
	developers := []models.RegisterRequest{
		{
			Email:    "dev1@example.com",
			Password: "password123",
			Name:     "Developer 1",
			Role:     "developer",
		},
		{
			Email:    "dev2@example.com",
			Password: "password123",
			Name:     "Developer 2",
			Role:     "developer",
		},
	}

	for _, dev := range developers {
		_, _ = authUseCase.Register(context.Background(), dev)
	}

	t.Run("get all developers", func(t *testing.T) {
		devs, err := authUseCase.GetDevelopers(context.Background())
		assert.NoError(t, err)
		assert.Len(t, devs, 2)
	})

	t.Run("no developers found", func(t *testing.T) {
		// Clear the mock repository
		mockRepo.users = make(map[string]*models.User)

		devs, err := authUseCase.GetDevelopers(context.Background())
		assert.Error(t, err)
		assert.Nil(t, devs)
	})
}
