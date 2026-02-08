package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProjectStatus string

const (
	ProjectStatusPending    ProjectStatus = "pending"
	ProjectStatusProcessing ProjectStatus = "processing"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusFailed     ProjectStatus = "failed"
)

type Project struct {
	ID            uuid.UUID        `json:"id" db:"id"`
	UserID        uuid.UUID        `json:"user_id" db:"user_id"`
	Name          string           `json:"name" db:"name"`
	Description   string           `json:"description" db:"description"`
	BaseKeywords  pq.StringArray   `json:"base_keywords" db:"base_keywords"`
	Status        ProjectStatus    `json:"status" db:"status"`
	ResultData    *json.RawMessage `json:"-" db:"result_data"`
	AiResultData  *json.RawMessage `json:"-" db:"ai_result_data"`
	SeoResultData *json.RawMessage `json:"-" db:"seo_result_data"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
}
