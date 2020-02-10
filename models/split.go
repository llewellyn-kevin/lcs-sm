package models

import(
	"github.com/jinzhu/gorm"
)

type Split struct {
	gorm.Model
	League string
	Season string
	Year   int
}