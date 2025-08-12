package repository

import (
	"axolotl-cloud/internal/app/model"
	"context"

	"gorm.io/gorm"
)

type ContainerRepository struct {
	DB *gorm.DB
}

func (repo *ContainerRepository) Create(ctx context.Context, container *model.Container) error {
	return repo.DB.WithContext(ctx).Create(container).Error
}

func (repo *ContainerRepository) FindAllByProjectID(ctx context.Context, projectID uint) ([]model.Container, error) {
	var containers []model.Container
	err := repo.DB.WithContext(ctx).Preload("LastJob").Where("project_id = ?", projectID).Find(&containers).Error
	return containers, err
}

func (repo *ContainerRepository) FindByID(ctx context.Context, id uint) (*model.Container, error) {
	var container model.Container
	err := repo.DB.WithContext(ctx).Preload("LastJob").First(&container, id).Error
	if err != nil {
		return nil, err
	}
	return &container, nil
}

func (repo *ContainerRepository) Save(ctx context.Context, container *model.Container) error {
	return repo.DB.WithContext(ctx).Save(container).Error
}

func (repo *ContainerRepository) Delete(ctx context.Context, id uint) error {
	return repo.DB.WithContext(ctx).Delete(&model.Container{}, id).Error
}

func (repo *ContainerRepository) GetAllContainers(ctx context.Context) ([]model.Container, error) {
	var containers []model.Container
	err := repo.DB.WithContext(ctx).Find(&containers).Error
	if err != nil {
		return nil, err
	}
	return containers, nil
}
