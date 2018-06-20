package model

import "github.com/jinzhu/gorm"

type (
	Todo struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	RUserBase struct {
		Count int `json:"count"`
	}

	MUserBase struct {
		Count int `json:"count"`
	}

	LookUpValue struct {
		Name string `json:"name"`
	}

	TransformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}


	User struct {
		Name string `json:"name"`
		Password string `json:"password"`
		RoleType int `json:"role_type"`
	}
)
// 设置RUserBase的表名为`RUserBase`
func (RUserBase) TableName() string {
	return "RUserBase"
}

func (MUserBase) TableName() string {
	return "MUserBase"
}

func (LookUpValue) TableName() string {
	return "LookUpValue"
}

func (User) TableName() string {
	return "users"
}

