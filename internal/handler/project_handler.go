package handler

import (
	"encoding/json"
	"net/http"
	"seo-app/internal/dto"
	"seo-app/internal/service"
	"strings"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: *service}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Description == "" {
		respondWithError(w, http.StatusBadRequest, "Description is required")
		return
	}

	project, err := h.service.CreateProject(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id := parts[len(parts)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	project, err := h.service.GetProjectById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	respondWithJSON(w, http.StatusOK, project)
}
func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetAllProjects()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}
	respondWithJSON(w, http.StatusOK, projects)
}
