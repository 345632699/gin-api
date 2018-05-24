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

/**
* @api {POST} /robot/monitor 机器人活跃率统计
* @apiGroup Robot
* @apiVersion 0.0.1
* @apiDescription 机器人活跃率统计，折线图统计
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*                ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE3Nzc5LCJpc3MiOiJ0ZXN0In0.6Gc46YfdqzEgl3FnJ6OigMC_zDZnSxJe7t63IwxgS_I&start_at=1527004800&end_at=1527091200&act=undefined&_=1527131381769
* @apiSuccess (200) {String} counts 机器人总数
* @apiSuccess (200) {String} Datetime 日期
* @apiSuccess (200) {String} ActivityCount 活跃数
* @apiSuccess (200) {String} SumTime 使用时长统计 活跃率统计接口时值为0
* @apiSuccessExample {json} 返回样例:
*				{"counts":55871,"data":[{"Datetime":"2018-05-22","ActivityCount":3741,"SumTime":0},{"Datetime":"2018-05-23","ActivityCount":3649,"SumTime":0},{"Datetime":"2018-05-24","ActivityCount":445,"SumTime":0}],"status":200}
*/
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

/**
* @api {POST} /robot/time_span_count 机器人待机时间统计
* @apiGroup Robot
* @apiVersion 0.0.1
* @apiDescription 机器人待机时间统计，折线图统计
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*                ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE3Nzc5LCJpc3MiOiJ0ZXN0In0.6Gc46YfdqzEgl3FnJ6OigMC_zDZnSxJe7t63IwxgS_I&start_at=1527004800&end_at=1527091200&act=undefined&_=1527131381769
* @apiSuccess (200) {String} Datetime 日期
* @apiSuccess (200) {String} ActivityCount 活跃数 使用时长统计接口时值为0
* @apiSuccess (200) {String} SumTime 使用时长统计
* @apiSuccessExample {json} 返回样例:
*				{"data":[{"Datetime":"2018-05-22","ActivityCount":0,"SumTime":79549164},{"Datetime":"2018-05-23","ActivityCount":0,"SumTime":76173190},{"Datetime":"2018-05-24","ActivityCount":0,"SumTime":1354214}],"status":200}
*/
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

/**
* @api {POST}   /mobile/use_times_spread 机器人在线时长每天分布统计
* @apiGroup Robot
* @apiVersion 0.0.1
* @apiDescription 机器人在线时长每天分布统计，扇形图接口
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*             ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE3Nzc5LCJpc3MiOiJ0ZXN0In0.6Gc46YfdqzEgl3FnJ6OigMC_zDZnSxJe7t63IwxgS_I&start_at=1526745600&end_at=1527091200&act=undefined&_=1527131751786
* @apiSuccess (200) {String} data {"日期":{detail.key:detail.num}}
* @apiSuccess (200) {String} detail.key-0 小于1小时
* @apiSuccess (200) {String} detail.key-1 大于一小时 小于三小时
* @apiSuccess (200) {String} detail.key-2 大于三小时 小于五小时
* @apiSuccess (200) {String} detail.key-3 大于5小时 小于8小时
* @apiSuccess (200) {String} detail.key-4 大于8小时
* @apiSuccessExample {json} 返回样例:
*			  {"data":[{"Datetime":"2018-05-22","Map":{"0":1792,"1":538,"2":220,"3":233,"4":958}}],"status":200}
*/
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