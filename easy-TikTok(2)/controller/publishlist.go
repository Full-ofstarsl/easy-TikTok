package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type publishlistrs struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户发布的视频列表
}

func PublishList(c *gin.Context) {
	id := c.Query("user_id")
	token := c.Query("token")
	//查找用户以及关联视频
	var user User
	result := DB.Where("usermaster_id=?", id).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusOK, publishlistrs{
			StatusCode: 0,
			StatusMsg:  "empty",
			VideoList:  nil,
		})
		return
	}
	var videolist []Video
	DB.Where("user_id=?", user.ID).Preload("User").Find(&videolist)
	for i := 0; i < len(videolist); i++ {
		Sum(token, &videolist[i].User)
	}
	c.JSON(http.StatusOK, publishlistrs{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videolist,
	})
}
