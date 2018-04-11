package controller

import (
	"fmt"
	"report/config"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"report/model"
)

type MobileActivityRes struct {
	Datetime string
	ActivityCount int
	SpanTimeLength int
}

//手機用戶活躍率
func MobileActivityRate(c *gin.Context)  {
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,
	count(OpTime) AS activity_count
	FROM
	(
		SELECT
			OpTime
		FROM
			StatisticOperation so
		WHERE
			so.Platform IN (0, 1)
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			UserId
	) app_count`
	var mobile MobileActivityRes
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	_result := mobile.doQuery(start_at,end_at,sql)
	db := config.Db
	var mUserBase model.MUserBase
	db.Select("count(*) as count").First(&mUserBase)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "mobile_counts": mUserBase.Count, "data": _result})
}

//手機使用時長
func MobileTimeSpanCount(c *gin.Context)  {
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	var res MobileActivityRes
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,
	SUM(count_time) AS span_time_length
	FROM
	(
		SELECT
			UserId,
			LookUpFunctionValueId,
			OpTime,
			max(OpTime) - min(OpTime) AS count_time
		FROM
			StatisticOperation
		WHERE
			LookUpFunctionValueId <> 5
		AND UserId = UserId
		AND Platform IN (0, 1)
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			LookUpFunctionValueId,
			UserId
	) a`
	_result := res.doQuery(start_at,end_at,sql)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _result})
}

//手機使用次數分佈
type TimeCount struct {
	Datetime string
	UseTimes int
}
func MobileUserTimeSpread(c *gin.Context)  {
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	db := config.Db
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,
	count(UserId) AS use_times
	FROM
	(
		SELECT
			OpTime,
			UserId
		FROM
			StatisticOperation so
		WHERE
			so.Platform IN (0, 1)
		AND OpTime BETWEEN ?
		AND ?
	) a
	GROUP BY
	UserId`
	var time_counts []TimeCount
	res := make(map[string]map[int]int)
	for start_at < end_at {
		m := make(map[int]int)
		day_end := start_at + 86400
		db.Raw(sql,start_at,day_end).Scan(&time_counts)
		for _,time_count := range time_counts{
			use_times := time_count.UseTimes
			switch {
			case use_times <= 1:
				m[0] += 1
			case 2 < use_times && use_times <= 5:
				m[1] += 1
			case 6 < use_times && use_times <= 10:
				m[2] += 1
			default:
				m[3] += 1
			}
		}
		datetime := time_counts[0].Datetime
		res[datetime] = m
		start_at += 86400
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": res})
}


func (c MobileActivityRes) doQuery(start_at int,end_at int,sql string) []MobileActivityRes {
	db := config.Db
	var result MobileActivityRes
	var _result []MobileActivityRes
	for start_at <= end_at {
		day_end := start_at + 86400
		db.Raw(sql, start_at,day_end).Scan(&result)
		fmt.Println(sql)
		fmt.Println(result)
		_result = append(_result,result)
		start_at += 86400
	}
	return _result
}