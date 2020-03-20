package Database

import (
	"github.com/jinzhu/gorm"
	"html/template"
)

func init() {
	DB := GetDatabase()
	DB.AutoMigrate(
		&RgbItem{},
		&Status{},
	)
	DB.FirstOrCreate(&Status{
		Model:    gorm.Model{},
		Sessions: 0,
		Items:    0,
	}, "id = ?", 1)
	return
}

type RgbItem struct {
	gorm.Model
	FileName   string        `json:"file_name"`
	Identifier string        `json:"identifier"`
	SessionID  string        `json:"-"`
	SrcHeight  int           `json:"src_height"`
	SrcWidth   int           `json:"src_width"`
	DstHeight  int           `json:"dst_height"`
	DstWidth   int           `json:"dst_width"`
	Value      template.HTML `json:"-"`
}

type Status struct {
	gorm.Model `json:"-"`
	Sessions   uint `json:"sessions"`
	Items      uint `json:"items"`
}
