package model

import "github.com/jinzhu/gorm"

type (
	Todo struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	RUserBase struct {
		Count string `json:"count"`
	}

	TransformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)
// 设置RUserBase的表名为`RUserBase`
func (RUserBase) TableName() string {
	return "RUserBase"
}