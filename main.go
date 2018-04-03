package main

import (
	//这里我们导入已经集成的 mysql 驱动，当然也可以导入原版的 import _ "github.com/go-sql-driver/mysql" 一样的
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"./routes"
	"net/http"
)


func main() {
	http.ListenAndServe(":8081", routes.Engine())
}

