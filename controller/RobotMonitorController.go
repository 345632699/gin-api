package controller

import (
	"github.com/gin-gonic/gin"
	"report/middleware/jwt"
	"net/http"
	"fmt"
	"report/config"
	"strconv"
	"report/model"
)

// Scan
type Result struct {
	Datetime string
	ActivityCount  int
}

func TestHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwtauth.CustomClaims)
	s := c.Query("test")
	res := Res{ID:11,Msg:s,UserName:claims.Name}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": res})
	fmt.Println(s)
}
//获取每天开机机器人的活跃数
func GetRobotActivityCount(c *gin.Context){
	c.MustGet("claims")
	//string 类型转换为int类型值
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	sql := "SELECT FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,count(1) * 14 AS activity_count FROM(SELECT * FROM StatisticOperation WHERE LookUpFunctionValueId <> 5 AND UserId = UserId AND Platform = 100  AND OpTime BETWEEN ? AND ? GROUP BY UserId ) a"
	db := config.Db
	var result Result
	var _result []Result
	for start_at < end_at {
		day_end := start_at + 86400
		db.Raw(sql, start_at,day_end).Scan(&result)
		_result = append(_result,result)
		start_at += 86400
	}
	var rUserBase model.RUserBase
	db.Select("count(*) as count").First(&rUserBase)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK,"counts": rUserBase.Count, "data": _result})

}