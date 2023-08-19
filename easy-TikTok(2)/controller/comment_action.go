package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type comment_res struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	Comment    Comment `json:"comment"`     // 评论成功返回评论内容，不需要重新拉取整个列表
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	action_type := c.Query("action_type")
	var usermaster Usermaster
	DB.Where("token=?", token).First(&usermaster)
	if action_type == "1" {
		comment_text := c.Query("comment_text")
		comment := Comment{
			CreateDate: time.Now().Unix(),
			VideoID:    uint(video_id),
			UserID:     uint(usermaster.ID),
			Content:    comment_text,
		}
		var video Video
		DB.First(&video, video_id)
		video.CommentCount += 1
		DB.Save(&video)
		DB.Create(&comment)
		Sum(token, &comment.User)
		c.JSON(http.StatusOK, comment_res{
			StatusCode: 0,
			StatusMsg:  "success",
			Comment:    comment,
		})
	} else if action_type == "2" {
		comment_id := c.Query("comment_id")
		DB.Delete(&Comment{}, comment_id)
		var video Video
		DB.First(&video, video_id)
		video.CommentCount -= 1
		DB.Save(&video)
		c.JSON(http.StatusOK, comment_res{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "error",
		})
	}

}
