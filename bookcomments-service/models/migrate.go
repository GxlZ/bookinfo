package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

func Migrate(db *gorm.DB) {
	fmt.Println("db migrate...")
	table(db)

	data(db)
}

func table(db *gorm.DB) {
	db.AutoMigrate(
		&Comments{},
	)
}

func data(db *gorm.DB) {
	for _, item := range comments {
		db.Create(&item)
	}
}