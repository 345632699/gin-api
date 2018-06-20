package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"report/middleware/jwt"
	"report/config"
	"report/model"
)

/**
* @api {POST} /login 用户登录
* @apiGroup Users
* @apiVersion 0.0.1
* @apiDescription 获取token验证值
* @apiParam {String} name 用户名
* @apiParam {String} password 密码
* @apiParamExample {json} 请求样例：
*                ?name=admin&password=admin123
* @apiSuccess (200) {String} code 200 代表无错误 1代表有错误
* @apiSuccess (200) {String} token 验证token
* @apiSuccessExample {json} 返回样例:
*             {"status":200,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwYXNzd29yZCI6ImFkbWluMTIzIiwiZXhwIjoxNTIzNjA1NTU0LCJpc3MiOiJ0ZXN0In0.f1Glpi2jLZe5nUnpOWbB-II6NtMom5D6Nq6oTYuF5nA"}
*/
func LoginHandler(c *gin.Context){
	var user model.User
	name := c.PostForm("name")
	password := c.PostForm("password")

	db := config.Db
	db.Where("name = ?",name).Find(&user)

	if  user.Name != "" && user.Password == password {
		j := &jwtauth.JWT{
			[]byte("test"),
		}
		claims := jwtauth.CustomClaims{
			name,
			password,
			jwt.StandardClaims{
				//ExpiresAt: 15000, //time.Now().Add(24 * time.Hour).Unix()
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				Issuer: "test",
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			c.String(http.StatusOK, err.Error())
			c.Abort()
		}
		//过期则重新输出
		//res, err := j.ParseToken(token)
		if err != nil {
			if err == jwtauth.TokenExpired {
				c.JSON(http.StatusOK, gin.H{"status": http.StatusUnauthorized, "msg": "登录验证超时，请重新登录"})
				c.Abort()
				return
				newToken, err := j.RefreshToken(token)
				if err != nil {
					c.String(http.StatusOK, err.Error())
				} else {
					c.String(http.StatusOK, newToken)
				}
			} else {
				c.String(http.StatusOK, err.Error())
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "token": token,})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusUnauthorized, "msg": "账号或者密码错误"})
	}

}