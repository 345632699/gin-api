package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"../middleware/jwt"
)

var users = gin.H{
	"admin":    "admin123",
}

func LoginHandler(c *gin.Context){
	name := c.PostForm("name")
	password := c.PostForm("password")

	if _, ok := users[name]; ok &&  users[name] == password {
		j := &jwtauth.JWT{
			[]byte("test"),
		}
		claims := jwtauth.CustomClaims{
			name,
			password,
			jwt.StandardClaims{
				//ExpiresAt: 15000, //time.Now().Add(24 * time.Hour).Unix()
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
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
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "token": token})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusUnauthorized, "msg": "账号或者密码错误"})
	}

}