package controller

import (
	"github.com/gin-gonic/gin"
	"report/config"
	"net/http"
	"strconv"
	"report/model"
	"time"
)

type Res struct {
	ID        uint   `json:"id"`
	DateTime     string `json:"datetime"`
	Name     string `json:"name"`
	TimeLengthCount string `json:"time_length_count"`
	TimesCount string `json:"times_count"`
}

//各應用使用時常統計
func TimeLengthCountHandler(c *gin.Context) {
	sql := `SELECT
		LookUpFunctionValueId as id,
		Description as name,
		sum(time_count) as time_length_count,
		FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS date_time
	FROM
	(
		SELECT
	UserId,
		LookUpFunctionValueId,
		OpTime,
		max(OpTime) - min(OpTime) AS time_count,
		lv.Description
	FROM
	StatisticOperation
	LEFT JOIN LookUpValue lv ON lv.UId = LookUpFunctionValueId
	WHERE
	LookUpFunctionValueId <> 5
	AND UserId = UserId
	AND Platform = 100
	AND OpTime BETWEEN ?
	AND ?
	GROUP BY
	LookUpFunctionValueId,
		UserId
	) a
	GROUP BY
	a.LookUpFunctionValueId`
	db := config.Db
	start_at,_:= strconv.ParseInt(c.Query("start_at"),10,64)
	end_at,_ := strconv.ParseInt(c.Query("end_at"),10,64)
	resMaps := make(map[string]map[string]string)
	var res Res
	var dateArr [] string
	for start_at < end_at {
		var results []Res
		day_end := start_at + 86400
		db.Raw(sql, start_at,day_end).Scan(&results)
		resMap := make(map[string]string)
		var datetime string
		for _,result := range results{
			resMap[result.Name] = result.TimeLengthCount
		}
		for _,app := range res.allApp(){
			if _, ok := resMap[app.Name]; ok == false {
				resMap[app.Name] = "0"
			}
		}
		tm := time.Unix(start_at,0)
		datetime = tm.Format("2006-01-02")
		resMaps[datetime] = resMap
		dateArr = append(dateArr, datetime)
		start_at += 86400
	}
	count_arr := res.formatRes(resMaps)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": gin.H{"date_arr":dateArr,"count_arr":count_arr},"count_num":res.allCountNum()})

}

func TimesCountHandler(c *gin.Context)  {
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS date_time,
	Description AS name,
	count AS times_count
FROM
	(
		SELECT
			UserId,
			LookUpFunctionValueId,
			OpTime,
			count(LookUpFunctionValueId) AS count,
			lv.Description
		FROM
			StatisticOperation
		LEFT JOIN LookUpValue lv ON lv.UId = LookUpFunctionValueId
		WHERE
			LookUpFunctionValueId <> 5
		AND UserId = UserId
		AND Platform = 100
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			LookUpFunctionValueId
	) a`
	db := config.Db
	start_at,_:= strconv.ParseInt(c.Query("start_at"),10,64)
	end_at,_ := strconv.ParseInt(c.Query("end_at"),10,64)
	resMaps := make(map[string]map[string]string)
	var res Res
	var dateArr [] string
	for start_at < end_at {
		var results []Res
		day_end := start_at + 86400
		db.Raw(sql, start_at,day_end).Scan(&results)
		resMap := make(map[string]string)
		var datetime string
		for _,result := range results{
			resMap[result.Name] = result.TimesCount
		}
		for _,app := range res.allApp(){
			if _, ok := resMap[app.Name]; ok == false {
				resMap[app.Name] = "0"
			}
		}
		tm := time.Unix(start_at,0)
		datetime = tm.Format("2006-01-02")
		resMaps[datetime] = resMap
		dateArr = append(dateArr, datetime)
		start_at += 86400
	}
	//拼接前端格式化数据
	m := make(map[string] string)
	for _,item := range resMaps{
		for k,v := range item{
			if m[k] == "" {
				m[k] = v
				continue
			}
			m[k] = m[k] + "," + v
		}
	}
	count_arr := res.formatRes(resMaps)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": gin.H{"date_arr":dateArr,"count_arr":count_arr},"count_num":res.allCountNum()})
}

func (res Res)allApp() [] model.LookUpValue  {
	db := config.Db
	var app  [] model.LookUpValue
	db.Select("Description as name").Where("TypeCode = ? AND Scope = ?","FUNCTION","ROBOT").Find(&app)
	return app
}

func (res Res)allCountNum() map[string]int {
	count_num := make(map[string]int)
	db := config.Db
	var mUserBase model.MUserBase
	db.Select("count(*) as count").First(&mUserBase)
	count_num["mobile_count_num"] = mUserBase.Count
	var rUserBase model.RUserBase
	db.Select("count(*) as count").First(&rUserBase)
	count_num["robot_count_num"] = rUserBase.Count
	return count_num
}

func (res Res)formatRes(resMaps map[string]map[string]string) map[string] string{
	m := make(map[string] string)
	for _,item := range resMaps{
		for k,v := range item{
			if m[k] == "" {
				m[k] = v
				continue
			}
			m[k] = m[k] + "," + v
		}
	}
	return m
}