package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func MessageAction(c *gin.Context) {
	token := c.Query("token")
	to_id, _ := strconv.Atoi(c.Query("to_user_id"))
	action_type := c.Query("action_type")
	content := c.Query("content")
	var usermaster Usermaster
	DB.Where("token=?", token).First(&usermaster)
	if action_type == "1" {
		send_id := usermaster.ID
		message := Message{
			Content:    content,
			CreateTime: time.Now().Unix(),
			FromUserID: send_id,
			ToUserID:   int64(to_id),
		}
		DB.Create(&message)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "success",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "fail",
		})
	}

}
