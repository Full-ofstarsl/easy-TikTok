package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type chatres struct {
	MessageList []Message `json:"message_list"` // 用户列表
	StatusCode  string    `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string    `json:"status_msg"`   // 返回状态描述
}

func MessageChat(c *gin.Context) {
	token := c.Query("token")
	to_id, _ := strconv.Atoi(c.Query("to_user_id"))
	time := c.Query("pre_msg_time")
	fmt.Println(time)
	var usermaster Usermaster
	DB.Where("token=?", token).First(&usermaster)
	var messagelist []Message
	result := DB.Where("from_user_id IN ? AND to_user_id IN ? AND create_time>?", []int64{usermaster.ID, int64(to_id)}, []int64{usermaster.ID, int64(to_id)}, time).Find(&messagelist)
	if result.Error != nil {
		fmt.Println("查找数据失败")
	}
	if len(messagelist) == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "message is empty",
		})
		return
	}
	c.JSON(http.StatusOK, chatres{
		StatusCode:  "0",
		StatusMsg:   "success",
		MessageList: messagelist,
	})
}
