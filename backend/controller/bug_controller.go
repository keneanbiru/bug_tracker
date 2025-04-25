package controller

import (
	"bug-tracker/models"
	"bug-tracker/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BugController struct {
	bugUseCase *usecase.BugUseCase
}

func NewBugController(bugUseCase *usecase.BugUseCase) *BugController {
	return &BugController{
		bugUseCase: bugUseCase,
	}
}

func (c *BugController) CreateBug(ctx *gin.Context) {
	var req models.CreateBugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.MustGet("user").(*models.User)

	bug, err := c.bugUseCase.CreateBug(ctx, req, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bug"})
		return
	}

	ctx.JSON(http.StatusCreated, bug)
}

func (c *BugController) GetBugs(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	var bugs []*models.BugResponse
	var err error

	if user.Role == "developer" {
		// Developers only see bugs assigned to them
		bugs, err = c.bugUseCase.GetBugsByDeveloper(ctx, user.ID)
	} else {
		// Managers and admins see all bugs
		bugs, err = c.bugUseCase.GetAllBugs(ctx)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bugs"})
		return
	}

	ctx.JSON(http.StatusOK, bugs)
}

func (c *BugController) UpdateBugStatus(ctx *gin.Context) {
	var req models.UpdateBugStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bugID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bug ID"})
		return
	}

	user := ctx.MustGet("user").(*models.User)

	bug, err := c.bugUseCase.UpdateBugStatus(ctx, bugID, req.Status, user.ID)
	if err != nil {
		switch err {
		case usecase.ErrBugNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Bug not found"})
		case usecase.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this bug"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bug status"})
		}
		return
	}

	ctx.JSON(http.StatusOK, bug)
}

func (c *BugController) AssignBug(ctx *gin.Context) {
	var req models.AssignBugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bugID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bug ID"})
		return
	}

	user := ctx.MustGet("user").(*models.User)
	if user.Role != "manager" && user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Only managers and admins can assign bugs"})
		return
	}

	bug, err := c.bugUseCase.AssignBug(ctx, bugID, req.DeveloperID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign bug"})
		return
	}

	ctx.JSON(http.StatusOK, bug)
}

func (c *BugController) GetBugByID(ctx *gin.Context) {
	bugID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bug ID"})
		return
	}

	bug, err := c.bugUseCase.GetBugByID(ctx, bugID)
	if err != nil {
		switch err {
		case usecase.ErrBugNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Bug not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bug"})
		}
		return
	}

	ctx.JSON(http.StatusOK, bug)
}

func (c *BugController) UpdateBug(ctx *gin.Context) {
	var req models.UpdateBugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bugID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bug ID"})
		return
	}

	user := ctx.MustGet("user").(*models.User)

	bug, err := c.bugUseCase.UpdateBug(ctx, bugID, req, user)
	if err != nil {
		switch err {
		case usecase.ErrBugNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Bug not found"})
		case usecase.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this bug"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bug"})
		}
		return
	}

	ctx.JSON(http.StatusOK, bug)
}

func (c *BugController) DeleteBug(ctx *gin.Context) {
	bugID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bug ID"})
		return
	}

	user := ctx.MustGet("user").(*models.User)

	// Only allow deletion by managers or admins
	if user.Role != "manager" && user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Only managers and admins can delete bugs"})
		return
	}

	err = c.bugUseCase.DeleteBug(ctx, bugID)
	if err != nil {
		switch err {
		case usecase.ErrBugNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Bug not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bug"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bug deleted successfully"})
}
