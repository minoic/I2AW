package Database

import (
	"github.com/jinzhu/gorm"
	"html/template"
)

func init() {
	DB := GetDatabase()
	DB.AutoMigrate(
		&Item{},
		&RgbItem{},
	)
	return
}

type Item struct {
	gorm.Model
	FileName   string
	Identifier string
	Value      string
}

type RgbItem struct {
	gorm.Model
	FileName   string
	Identifier string
	Value      template.HTML
}
