package models

import(
	"github.com/jinzhu/gorm"
)

type StockValue struct {
	gorm.Model
	TeamID  int
	SplitID int
	Week 	int
	Value	int
}