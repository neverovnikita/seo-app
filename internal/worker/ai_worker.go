package worker

import (
	"log"
	"seo-app/internal/service"
	"time"
)

type AiWorker struct {
	jobs     *service.JobService
	projects *service.ProjectService
	ai       *service.AIService
	stopChan chan bool
}

func NewAiWorker(jobs *service.JobService, projects *service.ProjectService, s *service.AIService) *AiWorker {
	return &AiWorker{jobs: jobs, projects: projects, ai: s}
}

func (w *AiWorker) Run() {
	log.Println("AI Worker запущен")
	for {
		select {
		case <-w.stopChan:
			log.Println("AI Worker остановлен")
			return
		default:
			w.processOne()
			time.Sleep(time.Second * 1)
		}
	}
}

func (w *AiWorker) Stop() {
	w.stopChan <- true
}

func (w *AiWorker) processOne() {
	id, err := w.jobs.GetJob("ai_processing")
	if err != nil {
		return
	}
	log.Printf("Обрабатываю проект %s через ИИ", id)

	project, err := w.projects.GetProjectById(id.String())
	if err != nil {
		log.Printf("Ошибка получения проекта %s: %v", id, err)
		return
	}

	keywords, err := w.ai.SendPrompt(project)
	if err != nil {
		log.Printf("Ошибка ИИ для проекта %s: %v", id, err)
		return
	}

	err = w.projects.UpdateAIResult(id, keywords)
	if err != nil {
		log.Printf("Ошибка сохранения результата для %s: %v", id, err)
		return
	}
	err = w.jobs.RemoveJob(id, "ai_processing")

	err = w.jobs.AddJob(id, "wordstat_processing")
	if err != nil {
		log.Printf("Ошибка добавления задачи wordstat для %s: %v", id, err)
	}

	log.Printf("Конец обработки проекта %s через ИИ", id)
}
