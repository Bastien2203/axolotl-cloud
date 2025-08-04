package repository

import (
	"axolotl-cloud/internal/app/model"
	"context"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func (repo *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	return repo.DB.WithContext(ctx).Create(p).Error
}

func (repo *ProjectRepository) FindAll(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	err := repo.DB.WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (repo *ProjectRepository) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	err := repo.DB.WithContext(ctx).First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (repo *ProjectRepository) Save(ctx context.Context, p *model.Project) error {
	return repo.DB.WithContext(ctx).Save(p).Error
}

func (repo *ProjectRepository) Delete(ctx context.Context, id uint) error {
	if err := repo.DeleteAssociations(ctx, &model.Project{ID: id}); err != nil {
		return err
	}
	return repo.DB.WithContext(ctx).Delete(&model.Project{}, id).Error
}

func (repo *ProjectRepository) DeleteAssociations(ctx context.Context, project *model.Project) error {
	return repo.DB.WithContext(ctx).Select("Containers").Delete(project).Error
}
