package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bug struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"`     // "open", "in-progress", "resolved"
	Priority    string             `bson:"priority" json:"priority"` // "low", "medium", "high", "critical"
	ReportedBy  primitive.ObjectID `bson:"reported_by" json:"reported_by"`
	AssignedTo  primitive.ObjectID `bson:"assigned_to,omitempty" json:"assigned_to,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateBugRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high critical"`
}

type UpdateBugRequest struct {
	Title       string `json:"title" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high critical"`
}

type UpdateBugStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=open in-progress resolved"`
}

type AssignBugRequest struct {
	DeveloperID primitive.ObjectID `json:"developer_id" binding:"required"`
}

type BugResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Priority    string             `json:"priority"`
	ReportedBy  UserResponse       `json:"reported_by"`
	AssignedTo  *UserResponse      `json:"assigned_to,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
