package db

import (
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/internal/app/model"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(shared.GetEnv("DATABASE_PATH")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	ScanFK(db)

	err = db.AutoMigrate(
		&model.Project{},
		&model.Container{},
		&model.Job{},
		&model.JobLog{},
		&model.Setting{},
	)

	return db, err
}

func ScanFK(db *gorm.DB) {
	// 2) vérifier que PRAGMA est bien activé
	var fkOn int
	db.Raw("PRAGMA foreign_keys;").Scan(&fkOn)
	fmt.Println("foreign_keys =", fkOn) // -> doit être 1

	// 3) inspecter les FKs créées par GORM
	type FKRow struct {
		Id       int    `gorm:"column:id"`
		Seq      int    `gorm:"column:seq"`
		Table    string `gorm:"column:table"`
		From     string `gorm:"column:from"`
		To       string `gorm:"column:to"`
		OnUpdate string `gorm:"column:on_update"`
		OnDelete string `gorm:"column:on_delete"`
		Match    string `gorm:"column:match"`
	}
	var fks []FKRow
	db.Raw("PRAGMA foreign_key_list('jobs');").Scan(&fks)
	fmt.Printf("jobs FKs: %+v\n", fks)

	var fks2 []FKRow
	db.Raw("PRAGMA foreign_key_list('job_logs');").Scan(&fks2)
	fmt.Printf("job_logs FKs: %+v\n", fks2)

	// 4) et afficher le SQL de création si besoin
	var sql string
	db.Raw("SELECT sql FROM sqlite_master WHERE type='table' AND name='jobs';").Scan(&sql)
	fmt.Println("jobs create sql:", sql)
}
