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