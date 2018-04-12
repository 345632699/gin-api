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

/**
* @api {POST} /app/time_length_count 各个应用使用时长统计
* @apiGroup Application
* @apiVersion 0.0.1
* @apiDescription 各个应用使用时长统计
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*                ?start_at=1522886400&end_at=1523059200&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTIzNjA1NTU0LCJpc3MiOiJ0ZXN0In0.f1Glpi2jLZe5nUnpOWbB-II6NtMom5D6Nq6oTYuF5nA
* @apiSuccess (200) {String} count_num {"mobile_count_num":手机客户端总数,"robot_count_num":机器人总数}
* @apiSuccess (200) {String} count_arr 各个应用使用时长列表，每天数据用,分割
* @apiSuccess (200) {String} date_arr 日期列表
* @apiSuccessExample {json} 返回样例:
*				{"count_num":{"mobile_count_num":27825,"robot_count_num":53649},"data":{"count_arr":{"CIBN微视听":"0,8683613","UtoVR":"0,0","中华美食":"0,3146688","中学课堂HD":"0,966390","互联急救":"0,0","作业帮":"0,2691380","发送语音留言":"0,4778008","听音乐":"0,48","咔哒故事":"0,428353","哆哆全脑思维":"0,747438","哆哆天才乐园":"0,418270","多屏互动":"0,86624","学习助手":"0,0","宝宝写数字":"0,1180141","宝宝学交通工具":"0,1172986","宝宝拼拼乐":"0,1774545","宝宝超市":"0,3646006","家庭医生签约(用户版)":"0,0","小企鹅乐园":"0,3267913","小咖哆哆":"0,527634","小哈学绘画":"0,504627","小哈操作视频":"0,1202072","小哈故事机":"0,1517020","小哈水族馆":"0,0","小哈菠萝树英语":"0,756895","小哈读绘本":"0,934429","小学课堂":"0,5529475","小猪佩奇-乔治感冒了":"0,0","小猪佩奇-乔治的新恐龙":"0,0","小猪佩奇-乔治的气球":"0,0","小猪佩奇-乔治第一天上幼儿园":"0,0","小猪佩奇-佩奇去划船":"0,0","小猪佩奇-佩奇去度假":"0,0","小猪佩奇-佩奇去滑雪":"0,0","小猪佩奇-佩奇家的电脑":"0,0","小猪佩奇-佩奇的新邻居":"0,0","小猪佩奇-佩奇的第一副眼镜":"0,0","小猪佩奇-去游泳":"0,0","小猪佩奇-快乐环保":"0,0","小猪佩奇-校车旅行":"0,0","小猪佩奇-森林小路":"0,0","小猪佩奇-游乐场":"0,0","小猪佩奇-看牙医":"0,0","小猪佩奇-第一次在朋友家过夜":"0,0","小猪佩奇-踢足球":"0,0","小猪佩奇-运动会":"0,0","小猪佩奇-露营去":"0,0","应用市场":"0,8177876","悟空数学":"0,809800","悟空识字":"0,1102720","打开微商城":"184000","收听语音留言":"0,1097344","斑马速算":"0,406978","水果猜猜乐":"0,90709","爱奇艺动画屋":"0,66422861","爱手工":"0,618438","相机":"0,74866","蛋生世界":"0,2451195","蛋生园":"0,1522829","语音交互":"0,610084","读绘本":"0,29","跟谁学":"0,528","运动加加":"0,0","防沉迷手动解锁":"0,19998129","防沉迷自动解锁":"0,19678251","高中同步课堂":"0,122310","魔力童英语":"0,39644","麦田园":"0,210904"},"date_arr":["2018-04-05","2018-04-06"]},"status":200}
*/
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
/**
* @api {POST} /app/times_count 各个应用使用次数统计
* @apiGroup Application
* @apiVersion 0.0.1
* @apiDescription 各个应用使用次数统计
* @apiParam {String} start_at 开始时间
* @apiParam {String} end_at 结束时间
* @apiParam {String} token token令牌
* @apiParamExample {json} 请求样例：
*                ?start_at=1522886400&end_at=1523059200&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTIzNjA1NTU0LCJpc3MiOiJ0ZXN0In0.f1Glpi2jLZe5nUnpOWbB-II6NtMom5D6Nq6oTYuF5nA
* @apiSuccess (200) {String} count_num {"mobile_count_num":手机客户端总数,"robot_count_num":机器人总数}
* @apiSuccess (200) {String} count_arr 各个应用使用次数列表，每天数据用,分割
* @apiSuccess (200) {String} date_arr 日期列表
* @apiSuccessExample {json} 返回样例:
*                {"count_num":{"mobile_count_num":27824,"robot_count_num":53649},"data":{"count_arr":{"CIBN微视听":"2347,2428","UtoVR":"0,0","中华美食":"1536,1520","中学课堂HD":"561,578","互联急救":"0,0","作业帮":"1459,1449","发送语音留言":"2700,2867","听音乐":"4,10","咔哒故事":"191,184","哆哆全脑思维":"424,433","哆哆天才乐园":"222,307","多屏互动":"41,19","学习助手":"0,0","宝宝写数字":"663,687","宝宝学交通工具":"593,547","宝宝拼拼乐":"1019,1108","宝宝超市":"1722,1708","家庭医生签约(用户版)":"0,0","小企鹅乐园":"1590,1319","小咖哆哆":"205,273","小哈学绘画":"366,388","小哈操作视频":"958,886","小哈故事机":"1128,1253","小哈水族馆":"0,0","小哈菠萝树英语":"715,727","小哈读绘本":"1053,964","小学课堂":"3372,3990","小猪佩奇-乔治感冒了":"0,0","小猪佩奇-乔治的新恐龙":"0,0","小猪佩奇-乔治的气球":"0,0","小猪佩奇-乔治第一天上幼儿园":"0,0","小猪佩奇-佩奇去划船":"0,0","小猪佩奇-佩奇去度假":"0,0","小猪佩奇-佩奇去滑雪":"0,0","小猪佩奇-佩奇家的电脑":"0,0","小猪佩奇-佩奇的新邻居":"1,1","小猪佩奇-佩奇的第一副眼镜":"0,0","小猪佩奇-去游泳":"2,0","小猪佩奇-快乐环保":"0,0","小猪佩奇-校车旅行":"0,0","小猪佩奇-森林小路":"0,0","小猪佩奇-游乐场":"1,0","小猪佩奇-看牙医":"0,0","小猪佩奇-第一次在朋友家过夜":"0,0","小猪佩奇-踢足球":"0,0","小猪佩奇-运动会":"0,0","小猪佩奇-露营去":"0,1","应用市场":"6239,5790","悟空数学":"487,398","悟空识字":"719,675","打开微商城":"418,430","收听语音留言":"1179,1065","斑马速算":"319,339","水果猜猜乐":"143,139","爱奇艺动画屋":"21770,22006","爱手工":"374,360","相机":"54,73","蛋生世界":"1966,2002","蛋生园":"1266,1348","语音交互":"664,725","读绘本":"4,2","跟谁学":"17,8","运动加加":"2,0","防沉迷手动解锁":"4483,4513","防沉迷自动解锁":"4293,4195","高中同步课堂":"276,160","魔力童英语":"154,122","麦田园":"133,114"},"date_arr":["2018-04-05","2018-04-06"]},"status":200}
*/
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
		for k,_ := range item{
			if m[k] == "" {
				m[k] = item[k]
				continue
			}
			m[k] = m[k] + "," + item[k]
		}
	}
	return m
}