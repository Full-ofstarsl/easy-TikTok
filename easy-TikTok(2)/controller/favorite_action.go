package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 点赞请求
type FavReq struct {
	ActionType string `json:"action_type"` // 1-点赞，2-取消点赞
	Token      string `json:"token"`       // 用户鉴权token
	VideoID    string `json:"video_id"`    // 视频id
}

// 点赞响应
type FavRes struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// 具体点赞操作
func DoFavorite(req *FavReq) (res *FavRes) {
	//通过token获取用户id
	user := GetUserIDByToken(req.Token)
	uid := user.ID
	fmt.Println(uid)
	//若状态码为1，则点赞，否则取消点赞
	//CTX := context.Background()
	if req.ActionType == "1" {
		RDB.SAdd(CTX, "video:"+req.VideoID+":favorite", uid)
		RDB.SRem(CTX, "video:"+req.VideoID+":unfavorite", uid)
	} else {
		RDB.SRem(CTX, "video:"+req.VideoID+":favorite", uid)
		RDB.SAdd(CTX, "video:"+req.VideoID+":unfavorite", uid)
	}
	//返回成功

	return &FavRes{
		StatusCode: 0,
		StatusMsg:  "success",
	}

}

// 点赞操作
func FavoriteAction(c *gin.Context) {
	var req FavReq
	var res *FavRes
	req.ActionType = c.Query("action_type")
	req.Token = c.Query("token")
	req.VideoID = c.Query("video_id")
	res = DoFavorite(&req)
	c.JSON(200, res)
	fmt.Println("点赞成功")

}

// 通过token获取用户id
func GetUserIDByToken(token string) (user Usermaster) {
	DB.Where("token = ?", token).Find(&user)
	return
}
