package repository

import (
	"encoding/json"
	"log"
	"seo-app/internal/dto"
	"seo-app/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProjectRepo struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) GetAll() (*[]models.Project, error) {
	var projects []models.Project
	query := `select * from projects`
	err := r.db.Select(&projects, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &projects, nil
}
func (r *ProjectRepo) Create(input dto.CreateProjectRequest) (*models.Project, error) {
	var project models.Project
	query := `
		insert into projects (user_id, name, description,base_keywords, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) 
        returning id, user_id, name, description, base_keywords, status, created_at, updated_at`
	now := time.Now()

	err := r.db.QueryRowx(query, input.UserID, input.Name, input.Description,
		pq.StringArray(input.BaseKeywords), now, now).StructScan(&project)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) GetById(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	query := `
		select * from projects 
		where id = $1`

	err := r.db.Get(&project, query, id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) UpdateStatus(id uuid.UUID, status string) error {
	_, err := r.db.Exec(
		"UPDATE projects SET status = $1 WHERE id = $2",
		status, id,
	)
	return err
}

func (r *ProjectRepo) UpdateAIResult(id uuid.UUID, result map[string]interface{}) error {
	resultJSON, _ := json.Marshal(result)
	_, err := r.db.Exec(
		"UPDATE projects SET ai_result_data = $1 WHERE id = $2",
		resultJSON, id,
	)
	return err
}

func (r *ProjectRepo) UpdateSeoResult(id uuid.UUID, result map[string]interface{}) error {
	resultJSON, _ := json.Marshal(result)
	_, err := r.db.Exec(
		"UPDATE projects SET seo_result_data = $1 WHERE id = $2",
		resultJSON, id,
	)
	return err
}

func (r *ProjectRepo) MarkCompleted(id uuid.UUID) error {
	_, err := r.db.Exec(
		"UPDATE projects SET status = 'completed' WHERE id = $1",
		id,
	)
	return err
}
