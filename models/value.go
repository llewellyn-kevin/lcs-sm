package models

import(
	"github.com/jinzhu/gorm"
)

type Value struct {
	gorm.Model
	TeamID  int
	SplitID int
	Week 	int
}