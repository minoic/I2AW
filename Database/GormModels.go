package Database

import (
	"github.com/jinzhu/gorm"
)

func init() {
	DB := GetDatabase()
	DB.AutoMigrate(
		&Item{},
	)
	return
}

type Item struct {
	gorm.Model
	FileName   string
	Identifier string
	Value      string
}
