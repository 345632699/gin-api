package routes
import (
	"fmt"
	"../middleware/jwt"
	"../controller"
	"net/http"
	"github.com/gin-gonic/gin"
)
func Engine() *gin.Engine {
	r := gin.Default()

	r.POST("/login",controller.LoginHandler)

	r.GET("/jwt", func(c *gin.Context) {

	})
	authorize := r.Group("/", jwtauth.JWTAuth())
	{
		authorize.GET("user", func(c *gin.Context) {
			claims := c.MustGet("claims").(*jwtauth.CustomClaims)
			fmt.Println(claims.Name)
			c.String(http.StatusOK, claims.Name)
		})
	}

	application := r.Group("/app",jwtauth.JWTAuth())
	{
		application.GET("times_count",controller.TimeCountHandler)
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