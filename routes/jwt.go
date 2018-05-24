package routes
import (
	"fmt"
	"report/middleware/jwt"
	"report/controller"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)
func Engine() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/login",controller.LoginHandler)

	//验证
	authorize := r.Group("/", jwtauth.JWTAuth())
	{
		authorize.Any("user", func(c *gin.Context) {
			claims := c.MustGet("claims").(*jwtauth.CustomClaims)
			fmt.Println(claims.Name)
			c.String(http.StatusOK, claims.Name)
		})
	}
	//app使用次数 时长相关统计路由
	application := r.Group("/app",jwtauth.JWTAuth())
	{
		application.GET("time_length_count",controller.TimeLengthCountHandler) //app总的使用时长统计 按天
		application.GET("times_count",controller.TimesCountHandler)  //app使用次数统计 按天
	}
	//机器人使用时长，次数，待机时间统计路由
	robot := r.Group("/robot",jwtauth.JWTAuth())
	{
		robot.GET("monitor",controller.GetRobotActivityCount) //机器人监控 活跃数
		robot.GET("time_span_count",controller.RobotTimeSpanCount) //机器人待机时长统计
		robot.GET("time_spread",controller.GetTimeLengthHandler) //机器人使用时长分布统计
	}
	mobile := r.Group("/mobile",jwtauth.JWTAuth())
	{
		mobile.GET("monitor",controller.MobileActivityRate) //手机端监控 活跃数
		mobile.GET("time_span_count",controller.MobileTimeSpanCount) //手机端使用待机时长统计
		mobile.GET("use_times_spread",controller.MobileUserTimeSpread) //手机端使用时长分布统计
	}
	//
	collect := r.Group("/collect",jwtauth.JWTAuth())
	{
		collect.GET("app",controller.RobotCoolect) //各个APP使用时长以及使用次数统计
		collect.GET("robot/month_activity",controller.ActivityUserByMonth)  //机器端 使用高活跃 低活跃统计 月
		collect.GET("mobile/month_activity",controller.ActivityMobileByMonth) //手机端 使用高活跃 低活跃统计 月
		collect.GET("count",controller.GetCounts)  //总时长 总次数 总活跃数 统一统计
	}

	r.GET("/dologin", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(
			200,
			`<form action="/login" method="POST"><input type="text" name="name"><input type="text" name="password"><input type="submit"></form>`,
			)
	})

	return r
}