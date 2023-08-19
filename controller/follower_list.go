package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type follower_res struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserList   []User `json:"user_list"`   // 用户信息列表
}

func FollowerList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	user_id = fmt.Sprintf("%sed", user_id)
	ids, _ := RDB.SMembers(CTX, user_id).Result()
	fmt.Println(ids)
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
	c.JSON(http.StatusOK, follower_res{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userlist,
	})
}
