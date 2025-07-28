package project

import (
	"context"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func (repo *ProjectRepository) Create(ctx context.Context, p *Project) error {
	return repo.DB.WithContext(ctx).Create(p).Error
}

func (repo *ProjectRepository) FindAll(ctx context.Context) ([]Project, error) {
	var projects []Project
	err := repo.DB.WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (repo *ProjectRepository) FindByID(ctx context.Context, id uint) (*Project, error) {
	var project Project
	err := repo.DB.WithContext(ctx).First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (repo *ProjectRepository) Save(ctx context.Context, p *Project) error {
	return repo.DB.WithContext(ctx).Save(p).Error
}

func (repo *ProjectRepository) Delete(ctx context.Context, id uint) error {
	return repo.DB.WithContext(ctx).Delete(&Project{}, id).Error
}
