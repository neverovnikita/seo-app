package service

import (
	"seo-app/internal/repository"

	"github.com/google/uuid"
)

type JobService struct {
	repo *repository.JobQueueRepo
}

func NewJobsService(repo *repository.JobQueueRepo) *JobService {
	return &JobService{repo: repo}
}

func (s *JobService) GetJob(stage string) (uuid.UUID, error) {
	id, err := s.repo.GetJob(stage)
	return id, err
}

func (s *JobService) AddJob(entityID uuid.UUID, stage string) error {
	err := s.repo.AddJob(entityID, stage)
	return err
}

func (s *JobService) RemoveJob(entityID uuid.UUID, stage string) error {
	err := s.repo.RemoveJob(entityID, stage)
	return err
}
