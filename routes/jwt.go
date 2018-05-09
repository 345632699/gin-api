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

	authorize := r.Group("/", jwtauth.JWTAuth())
	{
		authorize.Any("user", func(c *gin.Context) {
			claims := c.MustGet("claims").(*jwtauth.CustomClaims)
			fmt.Println(claims.Name)
			c.String(http.StatusOK, claims.Name)
		})
	}

	application := r.Group("/app",jwtauth.JWTAuth())
	{
		application.Any("time_length_count",controller.TimeLengthCountHandler)
		application.Any("times_count",controller.TimesCountHandler)
	}
	robot := r.Group("/robot",jwtauth.JWTAuth())
	{
		robot.Any("monitor",controller.GetRobotActivityCount)
		robot.Any("time_span_count",controller.RobotTimeSpanCount)
		robot.Any("time_spread",controller.GetTimeLengthHandler)
	}
	mobile := r.Group("/mobile",jwtauth.JWTAuth())
	{
		mobile.Any("monitor",controller.MobileActivityRate)
		mobile.Any("time_span_count",controller.MobileTimeSpanCount)
		mobile.Any("use_times_spread",controller.MobileUserTimeSpread)
	}
	collect := r.Group("/collect",jwtauth.JWTAuth())
	{
		collect.Any("app",controller.RobotCoolect)
		collect.Any("robot/month_activity",controller.ActivityUserByMonth)
		collect.Any("mobile/month_activity",controller.ActivityMobileByMonth)
		collect.Any("count",controller.GetCounts)
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