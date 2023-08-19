package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FriendList_res struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	UserList   []User `json:"user_list"`   // 用户信息列表
}

func FriendList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	ids, _ := RDB.SMembers(CTX, fmt.Sprintf("%sfo", user_id)).Result()
	var userlist []User
	var user User
	for _, id := range ids {
		flag, _ := RDB.SIsMember(CTX, fmt.Sprintf("%sfo", id), user_id).Result() //只有双方互相关注才为朋友关系，时间复杂度为o(n*n)
		if flag {
			result := DB.Where("usermaster_id=?", id).First(&user)
			if result.Error != nil {
				fmt.Println("查询friend用户失败")
			}
			userlist = append(userlist, user)
		}
	}
	for i := 0; i < len(userlist); i++ {
		Sum(token, &userlist[i])
	}
	c.JSON(http.StatusOK, FriendList_res{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userlist,
	})
}
