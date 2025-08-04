package db

import (
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/internal/app/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(shared.GetEnv("DATABASE_PATH")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.Project{},
		&model.Container{},
		&model.Job{},
		&model.JobLog{},
	)
	db.Exec("PRAGMA foreign_keys = ON")
	return db, err
}
