package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Userinfores struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	User       User   `json:"user,omitempty"`
}

func User_info(c *gin.Context) {
	id := c.Query("user_id")
	token := c.Query("token")
	var user User
	result := DB.Where("usermaster_id=?", id).First(&user)
	Sum(token, &user)
	if result.Error != nil {
		c.JSON(http.StatusOK, Userinfores{
			StatusCode: 1,
			StatusMsg:  "user not find",
		})
	} else {
		c.JSON(http.StatusOK, Userinfores{
			StatusCode: 0,
			StatusMsg:  "success",
			User:       user,
		})
	}
}
