package model

import "github.com/jinzhu/gorm"

type (
	Todo struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	TransformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)
