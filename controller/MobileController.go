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

/**
* @api {POST} /mobile/monitor 手机用户活跃率
* @apiGroup Mobile
* @apiVersion 0.0.1
* @apiDescription 手机用户活跃率，折线图统计
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*                ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE3Nzc5LCJpc3MiOiJ0ZXN0In0.6Gc46YfdqzEgl3FnJ6OigMC_zDZnSxJe7t63IwxgS_I&start_at=1527004800&end_at=1527091200&act=undefined&_=1527131381769
* @apiSuccess (200) {String} Datetime 日期
* @apiSuccess (200) {String} ActivityCount 活跃数
* @apiSuccess (200) {String} SpanTimeLength 使用时长统计 活跃率统计接口时值为0
* @apiSuccessExample {json} 返回样例:
*			 {"data":[{"Datetime":"2018-05-23","ActivityCount":1086,"SpanTimeLength":0},{"Datetime":"2018-05-24","ActivityCount":272,"SpanTimeLength":0}],"mobile_counts":32726,"status":200}
*/
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

/**
* @api {POST}  /mobile/time_span_count 手机用户使用时长统计
* @apiGroup Mobile
* @apiVersion 0.0.1
* @apiDescription 手机用户使用时长统计，折线图统计
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*               ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE3Nzc5LCJpc3MiOiJ0ZXN0In0.6Gc46YfdqzEgl3FnJ6OigMC_zDZnSxJe7t63IwxgS_I&start_at=1526659200&end_at=1527004800&act=undefined&_=1527131563278
* @apiSuccess (200) {String} Datetime 日期
* @apiSuccess (200) {String} ActivityCount 活跃数 使用时长统计接口时值为0
* @apiSuccess (200) {String} SpanTimeLength 使用时长
* @apiSuccessExample {json} 返回样例:
*				{"data":[{"Datetime":"2018-05-19","ActivityCount":0,"SpanTimeLength":10747479},{"Datetime":"2018-05-20","ActivityCount":0,"SpanTimeLength":9559515},{"Datetime":"2018-05-21","ActivityCount":0,"SpanTimeLength":5990751},{"Datetime":"2018-05-22","ActivityCount":0,"SpanTimeLength":6393240},{"Datetime":"2018-05-23","ActivityCount":0,"SpanTimeLength":5791652}],"status":200}
*/
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

/**
* @api {POST}   /mobile/use_times_spread 手机端打开次数分布
* @apiGroup Mobile
* @apiVersion 0.0.1
* @apiDescription 手机端打开次数分布统计，扇形图接口
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*             ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE3Nzc5LCJpc3MiOiJ0ZXN0In0.6Gc46YfdqzEgl3FnJ6OigMC_zDZnSxJe7t63IwxgS_I&start_at=1526745600&end_at=1527091200&act=undefined&_=1527131751786
* @apiSuccess (200) {String} data {"日期":{detail.key:detail.num}}
* @apiSuccess (200) {String} detail.key-0 小于1次
* @apiSuccess (200) {String} detail.key-1 大于2次小于等于5次
* @apiSuccess (200) {String} detail.key-2 大于6次小于等于20次
* @apiSuccess (200) {String} detail.key-3 大于20次
* @apiSuccessExample {json} 返回样例:
*			  {"data":{"2018-05-20":{"0":187,"1":298,"2":190,"3":803},"2018-05-21":{"0":150,"1":241,"2":169,"3":626},"2018-05-22":{"0":152,"1":224,"2":140,"3":558},"2018-05-23":{"0":163,"1":210,"2":137,"3":576}},"status":200}
*/
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