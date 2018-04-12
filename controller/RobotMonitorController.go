package controller

import (
	"github.com/gin-gonic/gin"
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
	SumTime int
}


//获取每天开机机器人的活跃数
func GetRobotActivityCount(c *gin.Context){
	//string 类型转换为int类型值
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,
	count(1) AS activity_count
	FROM
	(
		SELECT
			*
		FROM
			StatisticOperation
		WHERE
			LookUpFunctionValueId <> 5
		AND UserId = UserId
		AND Platform = 100
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			UserId
	) a`
	db := config.Db
	var res Result
	_result := res.doQuery(start_at,end_at,sql)
	var rUserBase model.RUserBase
	db.Select("count(*) as count").First(&rUserBase)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK,"counts": rUserBase.Count, "data": _result})
}

//机器人待机时间
func RobotTimeSpanCount(c *gin.Context) {
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	var res Result
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,
	SUM(count_time) AS sum_time
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
		AND Platform = 100
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			LookUpFunctionValueId,
			UserId
	) a`
	_result := res.doQuery(start_at,end_at,sql)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _result})

}

//机器人在线时长每天分布统计
type Spread struct {
	Datetime string
	Map map[int]int
}
type SumRes struct {
	Datetime string
	SumCountTime float64
}
func GetTimeLengthHandler(c *gin.Context) {
	start_at,_:= strconv.Atoi(c.Query("start_at"))
	end_at,_ := strconv.Atoi(c.Query("end_at"))
	sql := `SELECT
	FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS datetime,
	SUM(count_time) AS sum_count_time
	FROM
	(
		SELECT
			UserId AS user_id,
			LookUpFunctionValueId,
			OpTime,
			(max(OpTime) - min(OpTime)) / 3600 AS count_time
		FROM
			StatisticOperation
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
	user_id`
	db := config.Db
	var spread Spread
	var spreads []Spread

	for start_at < end_at {
		var results []SumRes
		day_end := start_at + 86400
		db.Raw(sql, start_at,day_end).Scan(&results)
		m := make(map[int]int)
		for _, result := range results {
			sum_count := result.SumCountTime
			// 0 小于1小时
			// 1 大于一小时 小于三小时
			// 2 大于三小时 小于五小时
			// 3 大于5小时 小于8小时
			// 4 大于8小时
			switch {
			case sum_count <= 1:
				m[0] += 1
			case 1 < sum_count && sum_count <= 3:
				m[1] += 1
			case 3 < sum_count && sum_count <= 5:
				m[2] += 1
			case 5 < sum_count && sum_count <= 8:
				m[3] += 1
			default:
				m[4] += 1
			}
			spread.Datetime = result.Datetime
		}
		spread.Map = m
		spreads = append(spreads, spread)
		start_at += 86400
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": spreads})
}

func (c Result) doQuery(start_at int,end_at int,sql string) []Result {
	db := config.Db
	var result Result
	var _result []Result
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