package usecase

import (
	"bug-tracker/models"
	"bug-tracker/repository"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBugNotFound  = errors.New("bug not found")
	ErrUnauthorized = errors.New("unauthorized action")
)

// BugUseCaseInterface defines the interface for bug use cases
type BugUseCaseInterface interface {
	CreateBug(ctx context.Context, req models.CreateBugRequest, reporterID primitive.ObjectID) (*models.BugResponse, error)
	GetBugByID(ctx context.Context, id primitive.ObjectID) (*models.BugResponse, error)
	GetAllBugs(ctx context.Context) ([]*models.BugResponse, error)
	GetBugsByDeveloper(ctx context.Context, developerID primitive.ObjectID) ([]*models.BugResponse, error)
	UpdateBugStatus(ctx context.Context, bugID primitive.ObjectID, status string, userID primitive.ObjectID) (*models.BugResponse, error)
	AssignBug(ctx context.Context, bugID, developerID primitive.ObjectID) (*models.BugResponse, error)
	UpdateBug(ctx context.Context, id primitive.ObjectID, req models.UpdateBugRequest, user *models.User) (*models.BugResponse, error)
	DeleteBug(ctx context.Context, id primitive.ObjectID) error
}

type BugUseCase struct {
	bugRepo  repository.BugRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

func NewBugUseCase(bugRepo repository.BugRepositoryInterface, userRepo repository.UserRepositoryInterface) *BugUseCase {
	return &BugUseCase{
		bugRepo:  bugRepo,
		userRepo: userRepo,
	}
}

func (uc *BugUseCase) CreateBug(ctx context.Context, req models.CreateBugRequest, reporterID primitive.ObjectID) (*models.BugResponse, error) {
	bug := &models.Bug{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		ReportedBy:  reporterID,
		Status:      "open",
	}

	if err := uc.bugRepo.Create(ctx, bug); err != nil {
		return nil, err
	}

	return uc.getBugResponse(ctx, bug)
}

func (uc *BugUseCase) GetAllBugs(ctx context.Context) ([]*models.BugResponse, error) {
	bugs, err := uc.bugRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.BugResponse, len(bugs))
	for i, bug := range bugs {
		response, err := uc.getBugResponse(ctx, bug)
		if err != nil {
			return nil, err
		}
		responses[i] = response
	}

	return responses, nil
}

func (uc *BugUseCase) GetBugsByDeveloper(ctx context.Context, developerID primitive.ObjectID) ([]*models.BugResponse, error) {
	bugs, err := uc.bugRepo.FindByAssignee(ctx, developerID)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.BugResponse, len(bugs))
	for i, bug := range bugs {
		response, err := uc.getBugResponse(ctx, bug)
		if err != nil {
			return nil, err
		}
		responses[i] = response
	}

	return responses, nil
}

func (uc *BugUseCase) UpdateBugStatus(ctx context.Context, bugID primitive.ObjectID, status string, userID primitive.ObjectID) (*models.BugResponse, error) {
	bug, err := uc.bugRepo.FindByID(ctx, bugID)
	if err != nil {
		return nil, err
	}
	if bug == nil {
		return nil, ErrBugNotFound
	}

	// Only assigned developer can update status
	if bug.AssignedTo != userID {
		return nil, ErrUnauthorized
	}

	if err := uc.bugRepo.UpdateStatus(ctx, bugID, status); err != nil {
		return nil, err
	}

	bug.Status = status
	return uc.getBugResponse(ctx, bug)
}

func (uc *BugUseCase) AssignBug(ctx context.Context, bugID, developerID primitive.ObjectID) (*models.BugResponse, error) {
	// Find the bug
	bug, err := uc.bugRepo.FindByID(ctx, bugID)
	if err != nil {
		return nil, err
	}
	if bug == nil {
		return nil, errors.New("bug not found")
	}

	// Find the developer
	developer, err := uc.userRepo.FindByID(ctx, developerID)
	if err != nil {
		return nil, err
	}
	if developer == nil {
		return nil, errors.New("user not found")
	}
	if developer.Role != "developer" {
		return nil, errors.New("invalid developer role")
	}

	// Assign the bug to the developer
	bug.AssignedTo = developerID
	err = uc.bugRepo.Update(ctx, bug)
	if err != nil {
		return nil, err
	}

	// Get the updated bug response
	response, err := uc.GetBugByID(ctx, bugID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (uc *BugUseCase) GetBugByID(ctx context.Context, id primitive.ObjectID) (*models.BugResponse, error) {
	bug, err := uc.bugRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if bug == nil {
		return nil, ErrBugNotFound
	}

	return uc.getBugResponse(ctx, bug)
}

func (uc *BugUseCase) UpdateBug(ctx context.Context, id primitive.ObjectID, req models.UpdateBugRequest, user *models.User) (*models.BugResponse, error) {
	bug, err := uc.bugRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if bug == nil {
		return nil, ErrBugNotFound
	}

	// Check permissions
	if user.Role != "admin" && user.Role != "manager" &&
		bug.ReportedBy != user.ID && bug.AssignedTo != user.ID {
		return nil, ErrUnauthorized
	}

	// Update fields if provided
	if req.Title != "" {
		bug.Title = req.Title
	}
	if req.Description != "" {
		bug.Description = req.Description
	}
	if req.Priority != "" {
		bug.Priority = req.Priority
	}

	if err := uc.bugRepo.Update(ctx, bug); err != nil {
		return nil, err
	}

	return uc.getBugResponse(ctx, bug)
}

func (uc *BugUseCase) DeleteBug(ctx context.Context, id primitive.ObjectID) error {
	bug, err := uc.bugRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if bug == nil {
		return ErrBugNotFound
	}

	return uc.bugRepo.Delete(ctx, id)
}

func (uc *BugUseCase) getBugResponse(ctx context.Context, bug *models.Bug) (*models.BugResponse, error) {
	reporter, err := uc.userRepo.FindByID(ctx, bug.ReportedBy)
	if err != nil {
		return nil, err
	}

	response := &models.BugResponse{
		ID:          bug.ID,
		Title:       bug.Title,
		Description: bug.Description,
		Status:      bug.Status,
		Priority:    bug.Priority,
		ReportedBy:  reporter.ToResponse(),
		CreatedAt:   bug.CreatedAt,
		UpdatedAt:   bug.UpdatedAt,
	}

	if !bug.AssignedTo.IsZero() {
		assignee, err := uc.userRepo.FindByID(ctx, bug.AssignedTo)
		if err != nil {
			return nil, err
		}
		assigneeResponse := assignee.ToResponse()
		response.AssignedTo = &assigneeResponse
	}

	return response, nil
}
