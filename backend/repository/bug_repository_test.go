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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDBURI = "mongodb+srv://keneanbiru:Godislove33.@cluster0.fek5tj1.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testDBURI))
	require.NoError(t, err)

	db := client.Database("bug_tracker_test")

	// Return a cleanup function
	cleanup := func() {
		// Drop collections instead of the entire database
		err := db.Collection("bugs").Drop(ctx)
		if err != nil {
			t.Logf("Warning: Failed to drop bugs collection: %v", err)
		}
		err = db.Collection("users").Drop(ctx)
		if err != nil {
			t.Logf("Warning: Failed to drop users collection: %v", err)
		}
		err = client.Disconnect(ctx)
		require.NoError(t, err)
	}

	return db, cleanup
}

func TestCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create test data
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		Title:       "Test Bug",
		Description: "This is a test bug",
		Status:      "open",
		Priority:    "high",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test case 1: Successful creation
	t.Run("Success", func(t *testing.T) {
		err := repo.Create(ctx, bug)
		assert.NoError(t, err)
		assert.NotEmpty(t, bug.ID)

		// Verify the bug was created in the database
		var savedBug models.Bug
		err = db.Collection("bugs").FindOne(ctx, bson.M{"_id": bug.ID}).Decode(&savedBug)
		assert.NoError(t, err)
		assert.Equal(t, bug.Title, savedBug.Title)
		assert.Equal(t, bug.Description, savedBug.Description)
		assert.Equal(t, bug.Status, savedBug.Status)
		assert.Equal(t, bug.Priority, savedBug.Priority)
		assert.Equal(t, bug.ReportedBy, savedBug.ReportedBy)
	})

	// Test case 2: Duplicate bug (same title and reporter)
	t.Run("Duplicate", func(t *testing.T) {
		dupBug := &models.Bug{
			Title:       bug.Title,
			Description: "Another description",
			Status:      "open",
			Priority:    "medium",
			ReportedBy:  reporterID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(ctx, dupBug)
		assert.NoError(t, err) // MongoDB doesn't enforce unique constraints by default
	})
}

func TestFindByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create a test bug
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		Title:       "Test Bug",
		Description: "This is a test bug",
		Status:      "open",
		Priority:    "high",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, bug)
	require.NoError(t, err)

	// Test case 1: Successful find
	t.Run("Success", func(t *testing.T) {
		foundBug, err := repo.FindByID(ctx, bug.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundBug)
		assert.Equal(t, bug.ID, foundBug.ID)
		assert.Equal(t, bug.Title, foundBug.Title)
		assert.Equal(t, bug.Description, foundBug.Description)
	})

	// Test case 2: Not found
	t.Run("Not Found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		foundBug, err := repo.FindByID(ctx, nonExistentID)
		assert.NoError(t, err)
		assert.Nil(t, foundBug)
	})
}

func TestFindAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create test bugs
	reporterID := primitive.NewObjectID()
	bugs := []*models.Bug{
		{
			Title:       "Bug 1",
			Description: "First bug",
			Status:      "open",
			Priority:    "high",
			ReportedBy:  reporterID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Bug 2",
			Description: "Second bug",
			Status:      "in-progress",
			Priority:    "medium",
			ReportedBy:  reporterID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, bug := range bugs {
		err := repo.Create(ctx, bug)
		require.NoError(t, err)
	}

	// Test case 1: Find all bugs
	t.Run("Success", func(t *testing.T) {
		foundBugs, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, foundBugs, len(bugs))
	})

	// Test case 2: Empty collection
	t.Run("Empty Collection", func(t *testing.T) {
		// Clear the collection
		err := db.Collection("bugs").Drop(ctx)
		require.NoError(t, err)

		foundBugs, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Empty(t, foundBugs)
	})
}

func TestFindByAssignee(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create test bugs with different assignees
	reporterID := primitive.NewObjectID()
	assigneeID1 := primitive.NewObjectID()
	assigneeID2 := primitive.NewObjectID()

	bugs := []*models.Bug{
		{
			Title:       "Bug 1",
			Description: "First bug",
			Status:      "open",
			Priority:    "high",
			ReportedBy:  reporterID,
			AssignedTo:  assigneeID1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Bug 2",
			Description: "Second bug",
			Status:      "in-progress",
			Priority:    "medium",
			ReportedBy:  reporterID,
			AssignedTo:  assigneeID1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Bug 3",
			Description: "Third bug",
			Status:      "open",
			Priority:    "low",
			ReportedBy:  reporterID,
			AssignedTo:  assigneeID2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, bug := range bugs {
		err := repo.Create(ctx, bug)
		require.NoError(t, err)
	}

	// Test case 1: Find bugs assigned to developer 1
	t.Run("Find Bugs for Developer 1", func(t *testing.T) {
		foundBugs, err := repo.FindByAssignee(ctx, assigneeID1)
		assert.NoError(t, err)
		assert.Len(t, foundBugs, 2)
		for _, bug := range foundBugs {
			assert.Equal(t, assigneeID1, bug.AssignedTo)
		}
	})

	// Test case 2: Find bugs assigned to developer 2
	t.Run("Find Bugs for Developer 2", func(t *testing.T) {
		foundBugs, err := repo.FindByAssignee(ctx, assigneeID2)
		assert.NoError(t, err)
		assert.Len(t, foundBugs, 1)
		assert.Equal(t, assigneeID2, foundBugs[0].AssignedTo)
	})

	// Test case 3: Find bugs for non-existent developer
	t.Run("Find Bugs for Non-existent Developer", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		foundBugs, err := repo.FindByAssignee(ctx, nonExistentID)
		assert.NoError(t, err)
		assert.Empty(t, foundBugs)
	})
}

func TestUpdate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create a test bug
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		Title:       "Original Title",
		Description: "Original description",
		Status:      "open",
		Priority:    "high",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, bug)
	require.NoError(t, err)

	// Test case 1: Successful update
	t.Run("Success", func(t *testing.T) {
		bug.Title = "Updated Title"
		bug.Description = "Updated description"
		bug.Priority = "medium"

		err := repo.Update(ctx, bug)
		assert.NoError(t, err)

		// Verify the update
		updatedBug, err := repo.FindByID(ctx, bug.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", updatedBug.Title)
		assert.Equal(t, "Updated description", updatedBug.Description)
		assert.Equal(t, "medium", updatedBug.Priority)
	})

	// Test case 2: Update non-existent bug
	t.Run("Not Found", func(t *testing.T) {
		nonExistentBug := &models.Bug{
			ID:          primitive.NewObjectID(),
			Title:       "Non-existent",
			Description: "This bug doesn't exist",
			Status:      "open",
			Priority:    "low",
			ReportedBy:  reporterID,
		}

		err := repo.Update(ctx, nonExistentBug)
		assert.NoError(t, err) // MongoDB's UpdateOne doesn't return error for non-existent documents
	})
}

func TestDelete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create a test bug
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		Title:       "Bug to Delete",
		Description: "This bug will be deleted",
		Status:      "open",
		Priority:    "high",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, bug)
	require.NoError(t, err)

	// Test case 1: Successful deletion
	t.Run("Success", func(t *testing.T) {
		err := repo.Delete(ctx, bug.ID)
		assert.NoError(t, err)

		// Verify the bug was deleted
		deletedBug, err := repo.FindByID(ctx, bug.ID)
		assert.NoError(t, err)
		assert.Nil(t, deletedBug)
	})

	// Test case 2: Delete non-existent bug
	t.Run("Not Found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.NoError(t, err) // MongoDB's DeleteOne doesn't return error for non-existent documents
	})
}

func TestUpdateStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create a test bug
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		Title:       "Test Bug",
		Description: "This is a test bug",
		Status:      "open",
		Priority:    "high",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, bug)
	require.NoError(t, err)

	// Test case 1: Successful status update
	t.Run("Success", func(t *testing.T) {
		newStatus := "in-progress"
		err := repo.UpdateStatus(ctx, bug.ID, newStatus)
		assert.NoError(t, err)

		// Verify the update
		updatedBug, err := repo.FindByID(ctx, bug.ID)
		assert.NoError(t, err)
		assert.Equal(t, newStatus, updatedBug.Status)
	})

	// Test case 2: Update non-existent bug
	t.Run("Not Found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.UpdateStatus(ctx, nonExistentID, "in-progress")
		assert.NoError(t, err) // MongoDB's UpdateOne doesn't return error for non-existent documents
	})
}

func TestAssignToDeveloper(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBugRepository(db)
	ctx := context.Background()

	// Create a test bug
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		Title:       "Test Bug",
		Description: "This is a test bug",
		Status:      "open",
		Priority:    "high",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, bug)
	require.NoError(t, err)

	// Test case 1: Successful assignment
	t.Run("Success", func(t *testing.T) {
		developerID := primitive.NewObjectID()
		err := repo.AssignToDeveloper(ctx, bug.ID, developerID)
		assert.NoError(t, err)

		// Verify the assignment
		updatedBug, err := repo.FindByID(ctx, bug.ID)
		assert.NoError(t, err)
		assert.Equal(t, developerID, updatedBug.AssignedTo)
	})

	// Test case 2: Assign to non-existent bug
	t.Run("Not Found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		developerID := primitive.NewObjectID()
		err := repo.AssignToDeveloper(ctx, nonExistentID, developerID)
		assert.NoError(t, err) // MongoDB's UpdateOne doesn't return error for non-existent documents
	})
}
