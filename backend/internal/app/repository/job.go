package repository

import (
	"axolotl-cloud/internal/app/model"

	"gorm.io/gorm"
)

type JobRepository struct {
	DB *gorm.DB
}

func (r *JobRepository) Save(job *model.Job) error {
	return r.DB.Save(job).Error
}

func (r *JobRepository) AddLog(jobID uint, line string) (*model.JobLog, error) {
	log := model.JobLog{JobID: jobID, Line: line}
	if err := r.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *JobRepository) UpdateStatus(jobID uint, status model.JobStatus) error {
	return r.DB.Model(&model.Job{}).Where("id = ?", jobID).Update("status", status).Error
}

func (r *JobRepository) GetAll() ([]model.Job, error) {
	var jobs []model.Job
	err := r.DB.Preload("Logs").Order("created_at desc").Find(&jobs).Error
	return jobs, err
}

func (r *JobRepository) GetByID(id uint) (*model.Job, error) {
	var job model.Job
	err := r.DB.Preload("Logs").First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) RemoveByID(id uint) error {
	return r.DB.Delete(&model.Job{}, id).Error
}
