package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type follow_res struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserList   []User `json:"user_list"`   // 用户信息列表
}

func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	ids, _ := RDB.SMembers(CTX, fmt.Sprintf("%sfo", user_id)).Result()
	var userlist []User
	var user User
	for _, id := range ids {
		result := DB.Where("usermaster_id=?", id).First(&user)
		if result.Error != nil {
			fmt.Println("查询用户失败")
		}
		userlist = append(userlist, user)
	}
	for i := 0; i < len(userlist); i++ {
		Sum(token, &userlist[i])
	}
	c.JSON(http.StatusOK, follow_res{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userlist,
	})
}
