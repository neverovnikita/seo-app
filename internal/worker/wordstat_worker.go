package worker

import (
	"log"
	"seo-app/internal/service"
	"time"
)

type WordstatWorker struct {
	jobs     *service.JobService
	projects *service.ProjectService
	wordstat *service.WordstatService
	stopChan chan bool
}

func NewWordstatWorker(jobs *service.JobService, projects *service.ProjectService, w *service.WordstatService) *WordstatWorker {
	return &WordstatWorker{jobs: jobs, projects: projects, wordstat: w}
}

func (w *WordstatWorker) Run() {
	log.Println("Wordstat Worker запущен")
	for {
		select {
		case <-w.stopChan:
			log.Println("Wordstat Worker остановлен")
			return
		default:
			w.processOne()
			time.Sleep(time.Second * 1)
		}
	}
}

func (w *WordstatWorker) Stop() {
	w.stopChan <- true
}

func (w *WordstatWorker) processOne() {
	id, err := w.jobs.GetJob("wordstat_processing")
	if err != nil {
		return
	}
	log.Printf("Обрабатываю проект %s через Wordstat", id)

	project, err := w.projects.GetProjectById(id.String())
	if err != nil {
		log.Printf("Ошибка получения проекта %s: %v", id, err)
		return
	}

	request := service.TopRequestsRequest{}
	if project.AiResultData != nil {
		if v, ok := project.AiResultData["keywords"]; ok {
			switch val := v.(type) {
			case []interface{}:
				phrases := make([]string, len(val))
				for i, item := range val {
					if str, ok := item.(string); ok {
						phrases[i] = str
					} else {
						log.Printf("Элемент %d не является строкой: %v", i, item)
						return
					}
				}
				request.Phrases = phrases
			case []string:
				request.Phrases = val
			default:
				log.Printf("Неожиданный тип keywords: %T", v)
				return
			}
		} else {
			log.Println("Поле 'keywords' не найдено в AiResultData")
			return
		}
	} else {
		log.Println("AiResultData пустой или nil")
		return
	}
	result, err := w.wordstat.RequestWordstat(&request)
	if err != nil {
		log.Printf("Ошибка ИИ для проекта %s: %v", id, err)
		return
	}

	err = w.projects.UpdateSeoResult(id, result)
	if err != nil {
		log.Printf("Ошибка сохранения результата для %s: %v", id, err)
		return
	}
	err = w.jobs.RemoveJob(id, "wordstat_processing")

	err = w.jobs.AddJob(id, "clusterizing")
	if err != nil {
		log.Printf("Ошибка добавления задачи clusterizing для %s: %v", id, err)
	}

	log.Printf("Конец обработки проекта %s через Wordstat", id)
}
