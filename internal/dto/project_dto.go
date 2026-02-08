package dto

import (
	"seo-app/internal/models"
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description" validate:"required"`
	BaseKeywords []string  `json:"base_keywords" validate:"required"`
}

type ProjectResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	BaseKeywords []string               `json:"base_keywords"`
	Status       models.ProjectStatus   `json:"status"` // "pending", "processing", "completed", "failed"
	ResultData   map[string]interface{} `json:"result_data,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
