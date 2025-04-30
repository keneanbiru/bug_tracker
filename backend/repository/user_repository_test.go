package repository

import (
	"bug-tracker/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Test case 1: Successful creation
	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "developer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)

		// Verify the user was created in the database
		var savedUser models.User
		err = db.Collection("users").FindOne(ctx, bson.M{"_id": user.ID}).Decode(&savedUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, savedUser.Name)
		assert.Equal(t, user.Email, savedUser.Email)
		assert.Equal(t, user.Password, savedUser.Password)
		assert.Equal(t, user.Role, savedUser.Role)
	})

	// Test case 2: Duplicate email
	t.Run("Duplicate Email", func(t *testing.T) {
		user1 := &models.User{
			Name:      "User 1",
			Email:     "same@example.com",
			Password:  "hashedpassword1",
			Role:      "developer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &models.User{
			Name:      "User 2",
			Email:     "same@example.com",
			Password:  "hashedpassword2",
			Role:      "developer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = repo.Create(ctx, user2)
		assert.Error(t, err) // Should fail due to unique email constraint
	})
}

func TestUserFindByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "developer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test case 1: Successful find
	t.Run("Success", func(t *testing.T) {
		foundUser, err := repo.FindByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test case 2: Not found
	t.Run("Not Found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		foundUser, err := repo.FindByID(ctx, nonExistentID)
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserFindByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "developer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test case 1: Successful find
	t.Run("Success", func(t *testing.T) {
		foundUser, err := repo.FindByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test case 2: Not found
	t.Run("Not Found", func(t *testing.T) {
		foundUser, err := repo.FindByEmail(ctx, "nonexistent@example.com")
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserFindByRole(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	users := []*models.User{
		{
			Name:      "Developer 1",
			Email:     "dev1@example.com",
			Password:  "hashedpassword1",
			Role:      "developer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Developer 2",
			Email:     "dev2@example.com",
			Password:  "hashedpassword2",
			Role:      "developer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Manager",
			Email:     "manager@example.com",
			Password:  "hashedpassword3",
			Role:      "manager",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Test case 1: Find developers
	t.Run("Find Developers", func(t *testing.T) {
		developers, err := repo.FindByRole(ctx, "developer")
		assert.NoError(t, err)
		assert.Len(t, developers, 2)
		for _, dev := range developers {
			assert.Equal(t, "developer", dev.Role)
		}
	})

	// Test case 2: Find managers
	t.Run("Find Managers", func(t *testing.T) {
		managers, err := repo.FindByRole(ctx, "manager")
		assert.NoError(t, err)
		assert.Len(t, managers, 1)
		assert.Equal(t, "manager", managers[0].Role)
	})

	// Test case 3: No users with role
	t.Run("No Users With Role", func(t *testing.T) {
		admins, err := repo.FindByRole(ctx, "admin")
		assert.NoError(t, err)
		assert.Empty(t, admins)
	})
}

func TestUserUpdate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		Name:      "Original Name",
		Email:     "original@example.com",
		Password:  "hashedpassword",
		Role:      "developer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test case 1: Successful update
	t.Run("Success", func(t *testing.T) {
		user.Name = "Updated Name"
		user.Email = "updated@example.com"

		err := repo.Update(ctx, user)
		assert.NoError(t, err)

		// Verify the update
		updatedUser, err := repo.FindByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.Name)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	// Test case 2: Update non-existent user
	t.Run("Not Found", func(t *testing.T) {
		nonExistentUser := &models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Non-existent",
			Email:     "nonexistent@example.com",
			Password:  "hashedpassword",
			Role:      "developer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Update(ctx, nonExistentUser)
		assert.NoError(t, err) // MongoDB's UpdateOne doesn't return error for non-existent documents
	})
}
