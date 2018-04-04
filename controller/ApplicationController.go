package controller

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"report/middleware/jwt"
)

type Res struct {
	ID        uint   `json:"id"`
	Msg     string `json:"title"`
	UserName string ` user name`
}

func TimeCountHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwtauth.CustomClaims)
	s := c.Query("test")
	res := Res{ID:11,Msg:s,UserName:claims.Name}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": res})
	fmt.Println(s)
}
