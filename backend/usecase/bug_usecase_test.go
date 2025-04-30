package usecase

import (
	"bug-tracker/models"
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockBugRepository struct {
	bugs map[primitive.ObjectID]*models.Bug
}

func NewMockBugRepository() *MockBugRepository {
	return &MockBugRepository{
		bugs: make(map[primitive.ObjectID]*models.Bug),
	}
}

func (m *MockBugRepository) Create(ctx context.Context, bug *models.Bug) error {
	if bug.ID.IsZero() {
		bug.ID = primitive.NewObjectID()
	}
	m.bugs[bug.ID] = bug
	return nil
}

func (m *MockBugRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Bug, error) {
	if bug, exists := m.bugs[id]; exists {
		return bug, nil
	}
	return nil, nil
}

func (m *MockBugRepository) FindAll(ctx context.Context) ([]*models.Bug, error) {
	bugs := make([]*models.Bug, 0, len(m.bugs))
	for _, bug := range m.bugs {
		bugs = append(bugs, bug)
	}
	// Sort bugs by ID to ensure consistent order
	sort.Slice(bugs, func(i, j int) bool {
		return bugs[i].ID.Hex() < bugs[j].ID.Hex()
	})
	return bugs, nil
}

func (m *MockBugRepository) FindByAssignee(ctx context.Context, assigneeID primitive.ObjectID) ([]*models.Bug, error) {
	var bugs []*models.Bug
	for _, bug := range m.bugs {
		if bug.AssignedTo == assigneeID {
			bugs = append(bugs, bug)
		}
	}
	return bugs, nil
}

func (m *MockBugRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	bug, exists := m.bugs[id]
	if !exists {
		return errors.New("bug not found")
	}
	bug.Status = status
	return nil
}

func (m *MockBugRepository) AssignToDeveloper(ctx context.Context, bugID, developerID primitive.ObjectID) error {
	bug, exists := m.bugs[bugID]
	if !exists {
		return errors.New("bug not found")
	}
	bug.AssignedTo = developerID
	return nil
}

func (m *MockBugRepository) Update(ctx context.Context, bug *models.Bug) error {
	if _, exists := m.bugs[bug.ID]; !exists {
		return errors.New("bug not found")
	}
	m.bugs[bug.ID] = bug
	return nil
}

func (m *MockBugRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if _, exists := m.bugs[id]; !exists {
		return errors.New("bug not found")
	}
	delete(m.bugs, id)
	return nil
}

func TestCreateBug(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create a test reporter
	reporterID := primitive.NewObjectID()
	reporter := &models.User{
		ID:    reporterID,
		Name:  "Test Reporter",
		Email: "reporter@example.com",
		Role:  "developer",
	}
	_ = mockUserRepo.Create(context.Background(), reporter)

	t.Run("successful bug creation", func(t *testing.T) {
		req := models.CreateBugRequest{
			Title:       "Test Bug",
			Description: "This is a test bug",
			Priority:    "high",
		}

		response, err := bugUseCase.CreateBug(context.Background(), req, reporterID)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Title, response.Title)
		assert.Equal(t, req.Description, response.Description)
		assert.Equal(t, req.Priority, response.Priority)
		assert.Equal(t, "open", response.Status)
		assert.Equal(t, reporterID, response.ReportedBy.ID)
	})
}

func TestGetBugByID(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create a test bug
	bugID := primitive.NewObjectID()
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		ID:          bugID,
		Title:       "Test Bug",
		Description: "This is a test bug",
		Priority:    "high",
		Status:      "open",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = mockBugRepo.Create(context.Background(), bug)

	// Create a test reporter
	reporter := &models.User{
		ID:    reporterID,
		Name:  "Test Reporter",
		Email: "reporter@example.com",
		Role:  "developer",
	}
	_ = mockUserRepo.Create(context.Background(), reporter)

	t.Run("successful bug retrieval", func(t *testing.T) {
		response, err := bugUseCase.GetBugByID(context.Background(), bugID)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, bugID, response.ID)
		assert.Equal(t, bug.Title, response.Title)
		assert.Equal(t, bug.Description, response.Description)
		assert.Equal(t, bug.Priority, response.Priority)
		assert.Equal(t, bug.Status, response.Status)
		assert.Equal(t, reporterID, response.ReportedBy.ID)
	})

	t.Run("bug not found", func(t *testing.T) {
		response, err := bugUseCase.GetBugByID(context.Background(), primitive.NewObjectID())
		assert.Error(t, err)
		assert.Equal(t, ErrBugNotFound, err)
		assert.Nil(t, response)
	})
}

func TestGetAllBugs(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create test bugs
	reporterID := primitive.NewObjectID()
	reporter := &models.User{
		ID:    reporterID,
		Name:  "Test Reporter",
		Email: "reporter@example.com",
		Role:  "developer",
	}
	_ = mockUserRepo.Create(context.Background(), reporter)

	bug1ID := primitive.NewObjectID()
	bug2ID := primitive.NewObjectID()
	bugs := []*models.Bug{
		{
			ID:          bug1ID,
			Title:       "Bug 1",
			Description: "First test bug",
			Priority:    "high",
			Status:      "open",
			ReportedBy:  reporterID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          bug2ID,
			Title:       "Bug 2",
			Description: "Second test bug",
			Priority:    "medium",
			Status:      "in-progress",
			ReportedBy:  reporterID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, bug := range bugs {
		_ = mockBugRepo.Create(context.Background(), bug)
	}

	t.Run("get all bugs", func(t *testing.T) {
		responses, err := bugUseCase.GetAllBugs(context.Background())
		assert.NoError(t, err)
		assert.Len(t, responses, len(bugs))

		// Sort responses by ID to match the order of bugs
		sort.Slice(responses, func(i, j int) bool {
			return responses[i].ID.Hex() < responses[j].ID.Hex()
		})

		for i, response := range responses {
			assert.Equal(t, bugs[i].ID, response.ID)
			assert.Equal(t, bugs[i].Title, response.Title)
			assert.Equal(t, bugs[i].Description, response.Description)
			assert.Equal(t, bugs[i].Priority, response.Priority)
			assert.Equal(t, bugs[i].Status, response.Status)
			assert.Equal(t, reporterID, response.ReportedBy.ID)
		}
	})
}

func TestUpdateBugStatus(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create a test bug
	bugID := primitive.NewObjectID()
	reporterID := primitive.NewObjectID()
	developerID := primitive.NewObjectID()
	bug := &models.Bug{
		ID:          bugID,
		Title:       "Test Bug",
		Description: "This is a test bug",
		Priority:    "high",
		Status:      "open",
		ReportedBy:  reporterID,
		AssignedTo:  developerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = mockBugRepo.Create(context.Background(), bug)

	// Create test users
	reporter := &models.User{
		ID:    reporterID,
		Name:  "Test Reporter",
		Email: "reporter@example.com",
		Role:  "developer",
	}
	developer := &models.User{
		ID:    developerID,
		Name:  "Test Developer",
		Email: "developer@example.com",
		Role:  "developer",
	}
	_ = mockUserRepo.Create(context.Background(), reporter)
	_ = mockUserRepo.Create(context.Background(), developer)

	t.Run("successful status update", func(t *testing.T) {
		response, err := bugUseCase.UpdateBugStatus(context.Background(), bugID, "in-progress", developerID)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "in-progress", response.Status)
	})

	t.Run("unauthorized update", func(t *testing.T) {
		response, err := bugUseCase.UpdateBugStatus(context.Background(), bugID, "in-progress", reporterID)
		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorized, err)
		assert.Nil(t, response)
	})

	t.Run("bug not found", func(t *testing.T) {
		response, err := bugUseCase.UpdateBugStatus(context.Background(), primitive.NewObjectID(), "in-progress", developerID)
		assert.Error(t, err)
		assert.Equal(t, ErrBugNotFound, err)
		assert.Nil(t, response)
	})
}

func TestAssignBug(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create a test bug
	bugID := primitive.NewObjectID()
	reporterID := primitive.NewObjectID()
	developerID := primitive.NewObjectID()
	invalidDeveloperID := primitive.NewObjectID()
	bug := &models.Bug{
		ID:          bugID,
		Title:       "Test Bug",
		Description: "This is a test bug",
		Priority:    "high",
		Status:      "open",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = mockBugRepo.Create(context.Background(), bug)

	// Create test users
	reporter := &models.User{
		ID:    reporterID,
		Name:  "Test Reporter",
		Email: "reporter@example.com",
		Role:  "developer",
	}
	developer := &models.User{
		ID:    developerID,
		Name:  "Test Developer",
		Email: "developer@example.com",
		Role:  "developer",
	}
	_ = mockUserRepo.Create(context.Background(), reporter)
	_ = mockUserRepo.Create(context.Background(), developer)

	t.Run("successful bug assignment", func(t *testing.T) {
		response, err := bugUseCase.AssignBug(context.Background(), bugID, developer.ID)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, developer.ID, response.AssignedTo.ID)
	})

	t.Run("invalid developer", func(t *testing.T) {
		response, err := bugUseCase.AssignBug(context.Background(), bugID, invalidDeveloperID)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, response)
	})

	t.Run("bug not found", func(t *testing.T) {
		response, err := bugUseCase.AssignBug(context.Background(), primitive.NewObjectID(), developer.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrBugNotFound, err)
		assert.Nil(t, response)
	})
}

func TestUpdateBug(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create a test bug
	bugID := primitive.NewObjectID()
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		ID:          bugID,
		Title:       "Test Bug",
		Description: "This is a test bug",
		Priority:    "high",
		Status:      "open",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = mockBugRepo.Create(context.Background(), bug)

	// Create test users
	reporter := &models.User{
		ID:    reporterID,
		Name:  "Test Reporter",
		Email: "reporter@example.com",
		Role:  "developer",
	}
	manager := &models.User{
		ID:    primitive.NewObjectID(),
		Name:  "Test Manager",
		Email: "manager@example.com",
		Role:  "manager",
	}
	_ = mockUserRepo.Create(context.Background(), reporter)
	_ = mockUserRepo.Create(context.Background(), manager)

	t.Run("successful bug update by reporter", func(t *testing.T) {
		req := models.UpdateBugRequest{
			Title:       "Updated Bug Title",
			Description: "Updated bug description",
			Priority:    "medium",
		}

		response, err := bugUseCase.UpdateBug(context.Background(), bugID, req, reporter)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Title, response.Title)
		assert.Equal(t, req.Description, response.Description)
		assert.Equal(t, req.Priority, response.Priority)
	})

	t.Run("successful bug update by manager", func(t *testing.T) {
		req := models.UpdateBugRequest{
			Title:       "Updated Bug Title by Manager",
			Description: "Updated bug description by manager",
			Priority:    "low",
		}

		response, err := bugUseCase.UpdateBug(context.Background(), bugID, req, manager)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Title, response.Title)
		assert.Equal(t, req.Description, response.Description)
		assert.Equal(t, req.Priority, response.Priority)
	})

	t.Run("unauthorized update", func(t *testing.T) {
		req := models.UpdateBugRequest{
			Title: "Unauthorized Update",
		}

		unauthorizedUser := &models.User{
			ID:    primitive.NewObjectID(),
			Name:  "Unauthorized User",
			Email: "unauthorized@example.com",
			Role:  "developer",
		}

		response, err := bugUseCase.UpdateBug(context.Background(), bugID, req, unauthorizedUser)
		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorized, err)
		assert.Nil(t, response)
	})

	t.Run("bug not found", func(t *testing.T) {
		req := models.UpdateBugRequest{
			Title: "Non-existent Bug Update",
		}

		response, err := bugUseCase.UpdateBug(context.Background(), primitive.NewObjectID(), req, reporter)
		assert.Error(t, err)
		assert.Equal(t, ErrBugNotFound, err)
		assert.Nil(t, response)
	})
}

func TestDeleteBug(t *testing.T) {
	mockBugRepo := NewMockBugRepository()
	mockUserRepo := NewMockUserRepository()
	bugUseCase := NewBugUseCase(mockBugRepo, mockUserRepo)

	// Create a test bug
	bugID := primitive.NewObjectID()
	reporterID := primitive.NewObjectID()
	bug := &models.Bug{
		ID:          bugID,
		Title:       "Test Bug",
		Description: "This is a test bug",
		Priority:    "high",
		Status:      "open",
		ReportedBy:  reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = mockBugRepo.Create(context.Background(), bug)

	t.Run("successful bug deletion", func(t *testing.T) {
		err := bugUseCase.DeleteBug(context.Background(), bugID)
		assert.NoError(t, err)

		// Verify bug is deleted
		_, err = bugUseCase.GetBugByID(context.Background(), bugID)
		assert.Error(t, err)
		assert.Equal(t, ErrBugNotFound, err)
	})

	t.Run("bug not found", func(t *testing.T) {
		err := bugUseCase.DeleteBug(context.Background(), primitive.NewObjectID())
		assert.Error(t, err)
		assert.Equal(t, ErrBugNotFound, err)
	})
}
