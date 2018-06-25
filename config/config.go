package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"report/model"
)

var (
	Db *gorm.DB
	sqlConnection = "root:root@(localhost)/go?charset=utf8&parseTime=True&loc=Local"
)

func init() {
	//打开数据库连接
	var err error
	Db, err = gorm.Open("mysql", sqlConnection)
	if err != nil {
		panic("failed to connect database")
	}

	Db.AutoMigrate(&model.Todo{})
}