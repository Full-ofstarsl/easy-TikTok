package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavListReq struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type FavListRes struct {
	StatusCode string  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户点赞视频列表
}

func DoFavoriteList(req *FavListReq) *FavListRes {

	uid := req.UserID

	//获取用户点赞视频列表
	videoList := GetFavoriteListByUserID(uid)
	return &FavListRes{
		StatusCode: "0",
		StatusMsg:  "成功返回喜欢列表",
		VideoList:  videoList,
	}
}

func FavoriteList(c *gin.Context) {
	var req FavListReq
	var res *FavListRes
	req.UserID = c.Query("user_id")
	req.Token = c.Query("token")
	res = DoFavoriteList(&req)
	c.JSON(200, res)
}

// 通过用户id获取用户点赞视频列表
func GetFavoriteListByUserID(uid string) []Video {
	// 查询点赞表获取用户点赞的视频列表
	var userFavorites []Favorite // 假设 Favorite 是点赞表的模型结构体
	fmt.Println("用户id", uid)
	DB.Where("user_id = ?", uid).Find(&userFavorites)
	//获取视频id列表
	var vidList []string
	for _, favorite := range userFavorites {
		id := strconv.Itoa(int(favorite.VideoID))
		vidList = append(vidList, id)
	}
	fmt.Println("喜欢的视频列表", vidList)
	//查询视频表获取视频列表
	var videoList []Video // 假设 Video 是视频表的模型结构体
	for _, vid := range vidList {
		var video Video
		DB.Where("id = ?", vid).Find(&video)
		video.IsFavorite = true
		videoList = append(videoList, video)
	}
	return videoList
}




// 通过用户id和视频id判断用户是否点赞，若点赞则修改视频的IsFavorite字段为true
func IsFavorite(uid string, vid string, v *Video) {
	//查询点赞表
	// var favorite Favorite
	// result := DB.Where("user_id = ? AND video_id = ?", uid, vid).Find(&favorite)
	// //若查询结果为空，则返回false
	// if result.Error != nil && result.Error.Error() == "record not found" {
	// 	v.IsFavorite = false
	// 	return
	// }
	// //若查询结果不为空，则返回true
	// if result.Error == nil {
	// 	v.IsFavorite = true
	// }

  
  //通过redis查询用户是否点赞过该视频
	flag, _ := RDB.SIsMember(CTX, fmt.Sprintf("video:%s:favorite", vid), uid).Result()
	if flag {
		v.IsFavorite = true
	} else {
		v.IsFavorite = false
	}
}


