package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"seo-app/internal/database"
	"seo-app/internal/handler"
	"seo-app/internal/repository"
	"seo-app/internal/service"
	"seo-app/internal/worker"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("C:\\Users\\Никита\\Desktop\\КУРСЫ\\seo-app\\.env"); err != nil {
		log.Println("No .env file found, using system environment")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/my_db?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port:  %d", serverPort)

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	log.Println("Database connected")
	apiKey := os.Getenv("AI_API_KEY")
	log.Println("API Key:", apiKey)
	apiURL := os.Getenv("AI_API_URL")
	model := os.Getenv("AI_MODEL")

	// make repo
	projectRepo := repository.NewProjectRepository(db)
	jobRepo := repository.NewQueueRepo(db)

	//make services
	jobService := service.NewJobsService(jobRepo)
	aiService := service.NewAIService(apiKey, apiURL, model, nil)
	projectService := service.NewProjectService(projectRepo, aiService, jobService)

	aiWorker := worker.NewAiWorker(jobService, projectService, aiService)
	go aiWorker.Run()

	projectHandler := handler.NewProjectHandler(projectService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/projects", methodHandler(projectHandler.CreateProject, "POST"))
	mux.HandleFunc("GET /api/v1/projects", methodHandler(projectHandler.GetAllProjects, "GET"))
	mux.HandleFunc("GET /api/v1/projects/{id}", methodHandler(projectHandler.GetProject, "GET"))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	loggedMux := loggingMiddleware(mux)

	serverAddr := ":" + serverPort
	log.Printf("Starting server on port %s", serverPort)

	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully started server on port: ", serverPort)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}
