package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	//"time"
)

type feedResponse struct {
	NextTime   int64   `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 视频列表
}

// 用于在返回用户信息时动态的调用redis中的数据才生成关注数，点赞数，是否已关注等相关信息
func Sum(token string, user *User) {
	var usermaster Usermaster
	DB.Where("token=?", token).First(&usermaster)
	flag, _ := RDB.SIsMember(CTX, fmt.Sprintf("%sfo", strconv.FormatInt(usermaster.ID, 10)), strconv.FormatInt(user.ID, 10)).Result()
	fmt.Println("检查：", fmt.Sprintf("%sfo", strconv.FormatInt(usermaster.ID, 10)))
	if flag || (usermaster.ID == user.ID) {
		user.IsFollow = true
	}
	user.FollowCount, _ = RDB.SCard(CTX, fmt.Sprintf("%sfo", strconv.FormatInt(user.ID, 10))).Result()
	user.FollowerCount, _ = RDB.SCard(CTX, fmt.Sprintf("%sed", strconv.FormatInt(user.ID, 10))).Result()
	user.FavoriteCount = int64(len(GetFavoriteListByUserID(strconv.FormatInt(user.ID, 10))))
	var sum int64 = 0
	videoids, _ := RDB.SMembers(CTX, fmt.Sprintf("author%s", strconv.FormatInt(user.ID, 10))).Result()
	for _, videoid := range videoids {
		count, _ := RDB.SCard(CTX, fmt.Sprintf("video:%s:favorite", videoid)).Result()
		fmt.Println(fmt.Sprintf("video:%s:favorite", videoid))
		sum += count
	}
	user.TotalFavorited = sum
}

func Feed(c *gin.Context) {
	token := c.Query("token")
	var videolist []Video
	DB.Preload("User").Find(&videolist) //通过preload实现预加载，完成结构体嵌套的调用
	for i := 0; i < len(videolist); i++ {
		Sum(token, &videolist[i].User)
		//通过token获得用户的id
		user := GetUserIDByToken(token)
		uid := strconv.Itoa(int(user.ID))
		vid := strconv.Itoa(int(videolist[i].ID))
		//通过用户id和视频id判断该用户是否点赞过该视频进行修改
		IsFavorite(uid, vid, &videolist[i])

	}
	c.JSON(http.StatusOK, feedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videolist,
		NextTime:   videolist[0].Time.Unix(),
	})
}
