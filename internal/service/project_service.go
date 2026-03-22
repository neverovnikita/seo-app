package service

import (
	"errors"
	"log"
	"seo-app/internal/dto"
	"seo-app/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repository.ProjectRepo
	jobs *JobService
	ai   *AIService
}

func NewProjectService(repo *repository.ProjectRepo, ai *AIService, jobs *JobService) *ProjectService {
	return &ProjectService{repo: repo, ai: ai, jobs: jobs}
}

func (s *ProjectService) CreateProject(input dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	if len(input.Name) < 3 {
		return nil, errors.New("name is too short")
	}
	if len(input.Description) < 100 {
		return nil, errors.New("description is too short")
	}
	if len(input.BaseKeywords) < 10 {
		return nil, errors.New("very few keywords, must be at least 10")
	}
	project, err := s.repo.Create(input)
	if err != nil {
		return nil, err
	}

	response := dto.ToProjectResponse(project)

	err = s.jobs.AddJob(project.ID, "ai_processing")
	if err != nil {
		log.Println("Не удалось добавить задачу в очередь: " + err.Error())
	}
	return response, nil
}

func (s *ProjectService) GetProjectById(id string) (*dto.ProjectResponse, error) {
	corrId, err := uuid.Parse(id)
	if err != nil {
		log.Println("project id is incorrect")
		return nil, err
	}
	project, err := s.repo.GetById(corrId)
	response := dto.ToProjectResponse(project)
	return response, nil
}

func (s *ProjectService) GetAllProjects() (*[]dto.ProjectResponse, error) {
	projects, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	response := dto.ToProjectListResponse(*projects)
	return &response, nil
}

func (s *ProjectService) UpdateStatus(id uuid.UUID, status string) error {
	err := s.repo.UpdateStatus(id, status)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) UpdateAIResult(id uuid.UUID, result KeywordsResult) error {
	data := map[string]interface{}{
		"keywords":   result.Keywords,
		"count":      result.Count,
		"updated_at": time.Now(),
	}
	err := s.repo.UpdateAIResult(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) UpdateSeoResult(id uuid.UUID, result *WordstatResponse) error {
	data := map[string]interface{}{
		"Results":   result.Results,
		"UpdatedAt": time.Now(),
	}
	err := s.repo.UpdateSeoResult(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) MarkCompleted(id uuid.UUID) error {
	err := s.repo.MarkCompleted(id)
	if err != nil {
		return err
	}
	return nil
}
