package models

import(
	"github.com/jinzhu/gorm"
)

type StockValue struct {
	gorm.Model
	TeamID  	int
	Team 		Team
	SplitID 	int
	Split		Split
	Week 		int
	Value		int
}