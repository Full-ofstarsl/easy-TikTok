package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func Publish(c *gin.Context) {
	//鉴权
	token := c.PostForm("token")
	var usermaster Usermaster
	result := DB.Where("token=?", token).First(&usermaster)
	if result.Error != nil {
		//fmt.Println("鉴权失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token is error!",
		})
	} else {
		//上传文件保存
		title := c.PostForm("title")
		data, _ := c.FormFile("data")
		filename := filepath.Base(data.Filename)
		saveFile := filepath.Join("./public/", filename)
		c.SaveUploadedFile(data, saveFile)
		Picture(fmt.Sprintf("./public/%s", filename), fmt.Sprintf("./public/%s_cover.jpg", strconv.FormatInt(usermaster.ID, 10)))

		//文件保存完成后对数据库video表增加和user表作品数的更新
		url := fmt.Sprintf("http://192.168.0.102:8080/public/%s", filename)
		coverurl := fmt.Sprintf("http://192.168.0.102:8080/public/%s_cover.jpg", strconv.FormatInt(usermaster.ID, 10))
		var video = Video{
			Title:    title,
			PlayURL:  url,
			CoverURL: coverurl,
			UserID:   uint(usermaster.ID),
			Time:     time.Now(),
		}
		DB.Create(&video)

		var user User
		DB.Where("usermaster_id = ?", usermaster.ID).First(&user)
		user.WorkCount += 1
		DB.Save(&user)
		RDB.SAdd(CTX, fmt.Sprintf("author%s", strconv.FormatInt(user.ID, 10)), strconv.FormatInt(int64(video.ID), 10))
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "文件上传成功"})
	}

}
