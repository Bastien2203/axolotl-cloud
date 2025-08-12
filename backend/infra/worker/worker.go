package worker

import (
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/infra/websocket"
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/internal/app/repository"
	"context"
	"fmt"
)

type Worker struct {
	Jobs chan *model.Job
	Repo *repository.JobRepository
}

func NewWorker(queueSize int, jobRepo *repository.JobRepository) *Worker {
	return &Worker{
		Jobs: make(chan *model.Job, queueSize),
		Repo: jobRepo,
	}
}

func (w *Worker) AddJob(job *model.Job, containerId *uint) (uint, error) {
	job.Status = model.JobStatusPending
	if containerId != nil {
		job.ContainerID = containerId
	}
	err := w.Repo.Save(job)
	if err != nil {
		logger.Error("Failed to save job:", err)
		return 0, fmt.Errorf("failed to save job: %w", err)
	}
	select {
	case w.Jobs <- job:
		return job.ID, nil
	default:
		fmt.Println("Job queue is full, dropping job:", job.Name)
		w.Repo.UpdateStatus(job.ID, model.JobStatusFailed)
		w.Repo.AddLog(job.ID, fmt.Sprintf("[ERROR] Job '%s' dropped due to full queue", job.Name))
		return 0, fmt.Errorf("job queue is full")
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case job := <-w.Jobs:
				go RunJob(ctx, job, w.Repo)
			}
		}
	}()
}

func RunJob(ctx context.Context, j *model.Job, repo *repository.JobRepository) {
	repo.UpdateStatus(j.ID, model.JobStatusRunning)
	repo.AddLog(j.ID, fmt.Sprintf("[INFO] Starting job: %s", j.Name))
	topicName := fmt.Sprintf("job:%d", j.ID)

	jobLogger := logger.NewLogger(func(level logger.LogLevel, msg string, args ...any) {
		line := fmt.Sprintf("[%s] %s", level, fmt.Sprintf(msg, args...))
		log, err := repo.AddLog(j.ID, line)
		if err != nil {
			logger.Error("Failed to add job log:", err)
			return
		}

		websocket.SendMessageToTopic(topicName, websocket.WSMessage[websocket.JobLogPayload]{
			Type: websocket.JobLogUpdateMessageType,
			Data: websocket.JobLogPayload{
				JobID: j.ID,
				Log:   *log,
			},
		})
	})

	if err := j.Run(ctx, jobLogger); err != nil {
		jobLogger.Error("Job failed", err.Error())
		repo.UpdateStatus(j.ID, model.JobStatusFailed)
		return
	}

	log, err := repo.AddLog(j.ID, fmt.Sprintf("[SUCCESS] Job '%s' completed", j.Name))
	if err != nil {
		jobLogger.Error("Failed to add job completion log:", err)
		return
	}
	websocket.SendMessageToTopic(topicName, websocket.WSMessage[websocket.JobLogPayload]{
		Type: websocket.JobLogUpdateMessageType,
		Data: websocket.JobLogPayload{
			JobID: j.ID,
			Log:   *log,
		},
	})
	repo.UpdateStatus(j.ID, model.JobStatusCompleted)
	websocket.UnsubscribeEveryoneFromTopic(topicName)
}
