package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JobQueueRepo struct {
	db *sqlx.DB
}

func NewQueueRepo(db *sqlx.DB) *JobQueueRepo {
	return &JobQueueRepo{db: db}
}

func (r *JobQueueRepo) AddJob(entityID uuid.UUID, stage string) error {
	_, err := r.db.Exec("INSERT INTO job_queue (entity_id, stage) VALUES ($1, $2)", entityID, stage)
	return err
}

func (r *JobQueueRepo) GetJob(stage string) (uuid.UUID, error) {
	var entityID uuid.UUID
	err := r.db.QueryRow("SELECT entity_id FROM job_queue WHERE stage = $1 order by created_at limit 1",
		stage).Scan(&entityID)
	return entityID, err
}

func (r *JobQueueRepo) RemoveJob(entityID uuid.UUID, stage string) error {
	_, err := r.db.Exec("DELETE FROM job_queue WHERE entity_id = $1 and stage = $2", entityID, stage)
	return err
}
