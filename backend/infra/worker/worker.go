package worker

import (
	"axolotl-cloud/infra/logger"
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

func (w *Worker) AddJob(job *model.Job, containerId uint) (uint, error) {
	job.Status = model.JobStatusPending
	job.ContainerID = containerId
	w.Repo.Save(job)
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

	log := func(logLine string) {
		repo.AddLog(j.ID, fmt.Sprintf("[INFO] %s", logLine))
		logger.Info(logLine)
	}

	if err := j.Run(ctx, log); err != nil {
		fmt.Println("Job failed:", j.Name, err)
		fmt.Printf("[ERROR] Job '%s' (id: %d) failed: %v\n", j.Name, j.ID, err)
		repo.AddLog(j.ID, fmt.Sprintf("[ERROR] Job '%s' failed: %v", j.Name, err))
		repo.UpdateStatus(j.ID, model.JobStatusFailed)
		return
	}

	repo.AddLog(j.ID, fmt.Sprintf("[SUCCESS] Job '%s' completed", j.Name))
	repo.UpdateStatus(j.ID, model.JobStatusCompleted)
}
