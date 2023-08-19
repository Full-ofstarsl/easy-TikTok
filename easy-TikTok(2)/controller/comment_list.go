package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type commentlist_res struct {
	CommentList []Comment `json:"comment_list"` // 评论列表
	StatusCode  int64     `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string    `json:"status_msg"`   // 返回状态描述
}

func CommentList(c *gin.Context) {
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	var commentlist []Comment
	DB.Where("video_id=?", video_id).Find(&commentlist)

	length := len(commentlist)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if commentlist[j].CreateDate > commentlist[j+1].CreateDate {
				commentlist[j], commentlist[j+1] = commentlist[j+1], commentlist[j]
			}
		}
	}
	c.JSON(http.StatusOK, commentlist_res{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: commentlist,
	})
}
