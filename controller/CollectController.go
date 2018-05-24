package controller

import (
	"github.com/gin-gonic/gin"
	"report/config"
	"strconv"
	"net/http"
	"fmt"
	"github.com/tealeg/xlsx"
	time2 "time"
)

type TimeCountResult struct {
	TimesCount int `json:"times_count"`
	PackageName string `json:"name"`
}

type TimeLengthCountResult struct {
	TimeLengthCount int `json:"time_length_count"`
	AppName string `json:"app_name"`
	PackageName string `json:"name"`
}

/**
* @api {POST} /collect/app 应用统计
* @apiGroup Collect
* @apiVersion 0.0.1
* @apiDescription 各个应用使用次数，時長统计(合計統計接口)
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*                ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE1ODU0LCJpc3MiOiJ0ZXN0In0.Zo0vBvD5nFBDg6_ca8N_7EVkAZ0cGdhhXZfiqv95Hwo&start_at=1525104000&end_at=1527129476&act=undefined&_=1527129467311
* @apiSuccess (200) {String} all_counts 总次数统计
* @apiSuccess (200) {String} all_length 总时长统计
* @apiSuccess (200) {String} data {"name":应用名称,"time_length":应用使用时长，"times_count":应用使用时长}
* @apiSuccessExample {json} 返回样例:
*            {"all_counts":46730,"all_length":857701356,"data":[{"name":"防沉迷自动解锁","time_length":88043611,"times_count":3299},{"name":"防沉迷手动解锁","time_length":88065446,"times_count":3548},{"name":"收听语音留言","time_length":6215152,"times_count":1432},{"name":"听音乐","time_length":129296,"times_count":31},{"name":"多屏互动","time_length":1365940,"times_count":106},{"name":"打开微商城","time_length":8925,"times_count":6},{"name":"读绘本","time_length":19119,"times_count":10},{"name":"发送语音留言","time_length":24124260,"times_count":3058},{"name":"语音交互","time_length":6678071,"times_count":365},{"name":"哆哆天才乐园","time_length":1358829,"times_count":354},{"name":"哆哆全脑思维","time_length":4989877,"times_count":499},{"name":"小咖哆哆","time_length":1823629,"times_count":281},{"name":"悟空数学","time_length":3733094,"times_count":407},{"name":"悟空识字","time_length":8918342,"times_count":596},{"name":"魔力童英语","time_length":721980,"times_count":244},{"name":"蛋生园","time_length":12523751,"times_count":2244},{"name":"蛋生世界","time_length":18323955,"times_count":2801},{"name":"相机","time_length":378725,"times_count":141},{"name":"作业帮","time_length":4286916,"times_count":391},{"name":"斑马速算","time_length":1644122,"times_count":299},{"name":"中学课堂HD","time_length":3533762,"times_count":494},{"name":"跟谁学","time_length":260995,"times_count":48},{"name":"爱手工","time_length":664863,"times_count":73},{"name":"咔哒故事","time_length":1799498,"times_count":187},{"name":"水果猜猜乐","time_length":351257,"times_count":230},{"name":"小哈操作视频","time_length":7625570,"times_count":1452},{"name":"小哈读绘本","time_length":5840755,"times_count":1090},{"name":"社区矫正机器人","time_length":67292,"times_count":2},{"name":"打开设置","time_length":52189,"times_count":6},{"name":"应用市场","time_length":46935991,"times_count":4094},{"name":"小哈故事机","time_length":13432064,"times_count":1844},{"name":"运动加加","time_length":62515,"times_count":22},{"name":"UtoVR","time_length":78889,"times_count":14},{"name":"小哈学绘画","time_length":2312208,"times_count":713},{"name":"麦田园","time_length":560902,"times_count":209},{"name":"爱奇艺动画屋","time_length":352056999,"times_count":5907},{"name":"宝宝学交通工具","time_length":7311272,"times_count":706},{"name":"中华美食","time_length":22232809,"times_count":1701},{"name":"宝宝写数字","time_length":9089429,"times_count":863},{"name":"宝宝超市","time_length":22908154,"times_count":1485},{"name":"宝宝拼拼乐","time_length":14149539,"times_count":1384},{"name":"小哈菠萝树英语","time_length":4566677,"times_count":1209},{"name":"小企鹅乐园","time_length":17190385,"times_count":371},{"name":"小哈电视助手","time_length":228529,"times_count":11},{"name":"学习助手","time_length":27166,"times_count":1},{"name":"小学课堂","time_length":26290469,"times_count":1694},{"name":"UstudyStudent","time_length":179006,"times_count":72},{"name":"爱家康","time_length":111317,"times_count":7},{"name":"Neo Notes","time_length":184963,"times_count":9},{"name":"CIBN微视听","time_length":23110123,"times_count":369},{"name":"小猪佩奇-乔治感冒了","time_length":3591,"times_count":4},{"name":"小猪佩奇-佩奇去度假","time_length":3663,"times_count":2},{"name":"小猪佩奇-踢足球","time_length":3129,"times_count":4},{"name":"小猪佩奇-第一次在朋友家过夜","time_length":3254,"times_count":5},{"name":"小猪佩奇-运动会","time_length":1129,"times_count":5},{"name":"小猪佩奇-乔治的新恐龙","time_length":572,"times_count":3},{"name":"小猪佩奇-森林小路","time_length":13932,"times_count":7},{"name":"小猪佩奇-佩奇的新邻居","time_length":568,"times_count":3},{"name":"小猪佩奇-校车旅行","time_length":268,"times_count":1},{"name":"小猪佩奇-乔治第一天上幼儿园","time_length":1638,"times_count":4},{"name":"小猪佩奇-快乐环保","time_length":289,"times_count":4},{"name":"小猪佩奇-游乐场","time_length":46,"times_count":4},{"name":"小猪佩奇-佩奇家的电脑","time_length":32374,"times_count":3},{"name":"小猪佩奇-看牙医","time_length":1422,"times_count":6},{"name":"小猪佩奇-露营去","time_length":23869,"times_count":9},{"name":"小猪佩奇-佩奇去滑雪","time_length":54729,"times_count":7},{"name":"小猪佩奇-佩奇去划船","time_length":148,"times_count":3},{"name":"小猪佩奇-去游泳","time_length":1168,"times_count":5},{"name":"小猪佩奇-佩奇的第一副眼镜","time_length":11635,"times_count":4},{"name":"小猪佩奇-乔治的气球","time_length":1554,"times_count":4},{"name":"高中同步课堂","time_length":973751,"times_count":264}],"status":200}
*/
func RobotCoolect(c *gin.Context)  {
	start_at,_:= strconv.ParseInt(c.Query("start_at"),10,64)
	end_at,_ := strconv.ParseInt(c.Query("end_at"),10,64)
	//次数统计sql
	times_count_sql := `SELECT
	NAME as package_name,
	count(NAME) AS times_count
	FROM
	(
		SELECT
			*
		FROM
			OpReport
		WHERE
			OpTime BETWEEN ?
		AND ?
		GROUP BY
			Name,
			BindRUId
	) a
	GROUP BY
	NAME`
	//时长统计sql
	time_length_count_sql := `SELECT
	app_name,
	SUM(time_count) as time_length_count,
	NAME as package_name
	FROM
	(
		SELECT
			*, Description AS app_name,
			max(OpTime) - min(OpTime) AS time_count,
			FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS date_time
		FROM
			OpReport
		LEFT JOIN LookUpValue ON LookUpValue.Value = OpReport.Name
		WHERE
			UId <> 5
		AND BindRUId = BindRUId
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			Name,
			BindRUId,
			date_time
	) a
	GROUP BY
	NAME`

	db := config.Db
	type result struct {
		Name string `json:"name"`
		TimeLength int `json:"time_length"` //时长统计
		TimesCount int `json:"times_count"` //次数统计
	}
	//使用时长统计 各应用
	var timesLengthRes [] TimeLengthCountResult
	db.Raw(time_length_count_sql,start_at,end_at).Scan(&timesLengthRes)
	m := make(map[string]result)
	for _,item := range timesLengthRes{
		var res result
		res.Name = item.AppName
		res.TimeLength = item.TimeLengthCount
		m[item.PackageName] = res
	}
	var Results []result
	//使用次数统计 各应用
	var timesCountRes  [] TimeCountResult
	var count int = 0
	var length int = 0
	db.Raw(times_count_sql,start_at,end_at).Scan(&timesCountRes)
	for _,item := range timesCountRes{
		if _,ok := m[item.PackageName];ok {
			r := m[item.PackageName]
			r.TimesCount = item.TimesCount
			Results = append(Results,r )
			count = count + r.TimesCount
			length = length + r.TimeLength
		}
	}
	act := c.Query("act")
	if act == "export"  {
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("机器人活跃数")
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		cell := row.AddCell()
		cell.Value = "应用"
		cell = row.AddCell()
		cell.Value = "使用次数"
		cell = row.AddCell()
		cell.Value = "次数占比"
		cell = row.AddCell()
		cell.Value = "使用时长"
		cell = row.AddCell()
		cell.Value = "时长占比"
		cell = row.AddCell()
		cell.Value = "总时长/s"
		cell = row.AddCell()
		cell.Value = "总次数"
		for _,v := range Results{
			row = sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度
			cell = row.AddCell()
			cell.Value = v.Name
			cell = row.AddCell()
			cell.Value = strconv.Itoa(v.TimesCount / 60 )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v.TimesCount * 100  / count)+ "%"
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v.TimeLength )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v.TimeLength * 100 / length ) + "%"
			cell = row.AddCell()
			cell.Value = strconv.Itoa( length )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( count )
		}
		time := strconv.FormatInt(time2.Now().Unix(),10)
		path := "./export/app/usage/"+time+".xlsx"
		err := file.Save(path)
		if err != nil {
			panic(err)
		}
		doExport(c,path)
		return
	}


	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"all_counts":count,"all_length":length,"data":Results})
}

type ActivityMonth struct {
	Count int
	MonthDate string
}

/**
* @api {POST} /collect/robot/month_activity 機器人活躍用戶統計
* @apiGroup Collect
* @apiVersion 0.0.1
* @apiDescription 機器人活躍用戶統計
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*               ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE1ODU0LCJpc3MiOiJ0ZXN0In0.Zo0vBvD5nFBDg6_ca8N_7EVkAZ0cGdhhXZfiqv95Hwo&start_at=1525104000&end_at=1527129779&act=undefined&_=1527129775561
* @apiSuccess (200) {String} data {"日期":{"低活跃":number,"月活跃":number,"高活跃":number}}
* @apiSuccessExample {json} 返回样例:
* 			  {"data":{"2018-05":{"低活跃":7299,"月活跃":14507,"高活跃":1847}},"status":200}
*/
func ActivityUserByMonth(c *gin.Context)  {
	start_at,_:= strconv.ParseInt(c.Query("start_at"),10,64)
	end_at,_ := strconv.ParseInt(c.Query("end_at"),10,64)
	sql := `SELECT
	count(UserId) as count,
	UserId,
	FROM_UNIXTIME(OpTime, '%Y-%m') AS month_date
	FROM
	(
		SELECT
			UserId,
			OpTime,
			FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS date_time
		FROM
			StatisticOperation
		WHERE
			Platform = 100
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			UserId,
			date_time
	) a
	GROUP BY
	UserId`
	resultMap := queryForResult(sql,start_at,end_at)
	act := c.Query("act")
	if act == "export"  {
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("机器人活跃数")
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		cell := row.AddCell()
		cell.Value = "月份"
		cell = row.AddCell()
		cell.Value = "高活跃"
		cell = row.AddCell()
		cell.Value = "低活跃"
		cell = row.AddCell()
		cell.Value = "低活跃"
		for k,v := range resultMap{
			row = sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度
			cell = row.AddCell()
			cell.Value = k
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v["高活跃"] )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v["低活跃"] )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v["低活跃"] )
		}
		time := strconv.FormatInt(time2.Now().Unix(),10)
		path := "./export/robot/month_activity/"+time+".xlsx"
		err := file.Save(path)
		if err != nil {
			panic(err)
		}
		doExport(c,path)
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"data":resultMap})
}
/**
* @api {POST} /collect/mobile/month_activity 手機段活躍用戶統計
* @apiGroup Collect
* @apiVersion 0.0.1
* @apiDescription 手機段活躍用戶統計
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*               ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE1ODU0LCJpc3MiOiJ0ZXN0In0.Zo0vBvD5nFBDg6_ca8N_7EVkAZ0cGdhhXZfiqv95Hwo&start_at=1525104000&end_at=1527129779&act=undefined&_=1527129775561
* @apiSuccess (200) {String} data {"日期":{"低活跃":number,"月活跃":number,"高活跃":number}}
* @apiSuccessExample {json} 返回样例:
* 			  {"data":{"2018-05":{"低活跃":8453,"月活跃":9630,"高活跃":121}},"status":200}
*/
func ActivityMobileByMonth(c *gin.Context)  {
	start_at,_:= strconv.ParseInt(c.Query("start_at"),10,64)
	end_at,_ := strconv.ParseInt(c.Query("end_at"),10,64)
	sql := `SELECT
	count(UserId) as count,
	UserId,
	FROM_UNIXTIME(OpTime, '%Y-%m') AS month_date
	FROM
	(
		SELECT
			UserId,
			OpTime,
			FROM_UNIXTIME(OpTime, '%Y-%m-%d') AS date_time
		FROM
			StatisticOperation
		WHERE
			Platform in(0,1)
		AND OpTime BETWEEN ?
		AND ?
		GROUP BY
			UserId,
			date_time
	) a
	GROUP BY
	UserId`
	resultMap := queryForResult(sql,start_at,end_at)
	act := c.Query("act")
	if act == "export"  {
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("手机端活跃数")
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		cell := row.AddCell()
		cell.Value = "月份"
		cell = row.AddCell()
		cell.Value = "高活跃"
		cell = row.AddCell()
		cell.Value = "低活跃"
		cell = row.AddCell()
		cell.Value = "低活跃"
		for k,v := range resultMap{
			row = sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度
			cell = row.AddCell()
			cell.Value = k
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v["高活跃"] )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v["低活跃"] )
			cell = row.AddCell()
			cell.Value = strconv.Itoa( v["低活跃"] )
		}
		time := strconv.FormatInt(time2.Now().Unix(),10)
		path := "./export/robot/month_activity/"+time+".xlsx"
		err := file.Save(path)
		if err != nil {
			panic(err)
		}
		doExport(c,path)
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"data":resultMap})
}


func queryForResult(sql string,start_at,end_at int64) map[string]map[string]int  {
	db := config.Db
	var res []ActivityMonth
	m := make(map[string][]int)
	fmt.Println(start_at)
	fmt.Println(end_at)
	db.Raw(sql,start_at,end_at).Scan(&res)
	for _,item := range res{
		m[item.MonthDate] = append(m[item.MonthDate], item.Count)
	}

	resultMap := make(map[string]map[string]int)
	for k,v := range m{
		arr := make(map[string]int)
		arr["高活跃"] = 0
		arr["低活跃"] = 0
		for _,count := range v{
			if count > 15 {
				arr["高活跃"] += 1
			}else if count >=1 && count <= 5{
				arr["低活跃"] += 1
			}
		}
		arr["月活跃"] = len(v)
		resultMap[k] = arr
	}
	return resultMap
}

type countRes struct {
	OnlineRobotCount int
	RobotCount int
	MobileCount int
}

/**
* @api {POST} /collect/count 机器人，手机端总数统计
* @apiGroup Collect
* @apiVersion 0.0.1
* @apiDescription  机器人，手机端总数统计，包含在线机器人数统计
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*               ?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTI3MjE1ODU0LCJpc3MiOiJ0ZXN0In0.Zo0vBvD5nFBDg6_ca8N_7EVkAZ0cGdhhXZfiqv95Hwo&_=1527129775561
* @apiSuccess (200) {String} OnlineRobotCount 在綫機器人數
* @apiSuccess (200) {String} RobotCount 机器人总数
* @apiSuccess (200) {String} MobileCount 手机端总数
* @apiSuccessExample {json} 返回样例:
* 			  {"data":{"OnlineRobotCount":1726,"RobotCount":55871,"MobileCount":32725},"status":200}
*/

func GetCounts(c *gin.Context)  {
	db := config.Db
	online_robot_count_sql := `SELECT count(*) as online_robot_count from RUserBase where IsOnline=1`
	robot_count_sql := `SELECT count(*) as robot_count from RUserBase`
	mobile_count_sql := `SELECT count(*) as mobile_count from MUserBase`
	var res countRes
	db.Raw(online_robot_count_sql).Scan(&res)
	db.Raw(robot_count_sql).Scan(&res)
	db.Raw(mobile_count_sql).Scan(&res)
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"data":res})
}

func doExport(c *gin.Context,path string)  {
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", path))//文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.File(path)
	return
}