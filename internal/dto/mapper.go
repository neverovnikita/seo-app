package dto

import (
	"encoding/json"
	"seo-app/internal/models"
)

func ToProjectResponse(project *models.Project) *ProjectResponse {
	if project == nil {
		return nil
	}
	var keywords []string
	if project.BaseKeywords != nil {
		keywords = project.BaseKeywords
	}

	response := ProjectResponse{
		ID:           project.ID,
		Name:         project.Name,
		Description:  project.Description,
		BaseKeywords: keywords,
		Status:       project.Status,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
	}

	if project.ResultData != nil && len(*project.ResultData) > 0 {
		var data map[string]interface{}
		if err := json.Unmarshal(*project.ResultData, &data); err == nil {
			response.ResultData = data
		}
	}

	if project.AiResultData != nil && len(*project.AiResultData) > 0 {
		var data map[string]interface{}
		if err := json.Unmarshal(*project.AiResultData, &data); err == nil {
			response.AiResultData = data
		}
	}

	if project.SeoResultData != nil && len(*project.SeoResultData) > 0 {
		var data map[string]interface{}
		if err := json.Unmarshal(*project.SeoResultData, &data); err == nil {
			response.SeoResultData = data
		}
	}

	return &response
}

func ToProjectListResponse(projects []models.Project) []ProjectResponse {
	responses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *ToProjectResponse(&project)
	}
	return responses
}
