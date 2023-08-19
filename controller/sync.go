package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 定时任务，每隔一段时间将redis中的数据同步到mysql中
func Sync() {
	tricker := time.NewTicker(time.Second * 5)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-tricker.C:
				//将redis中的点赞数据同步到mysql中
				SyncFavorite()
			case <-quit:
				tricker.Stop()
				return
			}
		}
	}()

	fmt.Println("同步结束")
}

var i int = 0

// 每隔一秒将redis中的点赞数据同步到mysql中
func SyncFavorite() {
	CTX := context.Background()
	//遍历点赞视频id
	keys, _ := RDB.Keys(CTX, "video:*:favorite").Result()
	videoIDs := make([]string, 0)
	for _, key := range keys {
		// 将命名空间和 ":favorite" 部分去除，只保留视频 ID
		videoID := strings.TrimPrefix(key, "video:")
		videoID = strings.TrimSuffix(videoID, ":favorite")
		videoIDs = append(videoIDs, videoID)
	}
	//fmt.Println("点赞视频id：", videoIDs)
	//写入mysql
	for _, videoID := range videoIDs {
		key := "video:" + videoID + ":favorite"
		userIDs, _ := RDB.SMembers(CTX, key).Result()
		//fmt.Println("点赞用户id：", userIDs)
		//遍历点赞用户id
		for _, userID := range userIDs {
			uid, _ := strconv.ParseUint(userID, 10, 64)
			vid, _ := strconv.ParseUint(videoID, 10, 64)
			//将点赞数据同步增加到点赞表中
			favorite := Favorite{
				UserID:  uint(uid),
				VideoID: uint(vid),
			}
			//判断是否存在，不存在则创建
			res := DB.First(&favorite, "user_id = ? AND video_id = ?", uid, vid)

			if res.Error != nil && res.Error.Error() == "record not found" {
				res = DB.Create(&favorite)
				if res.Error != nil {
					fmt.Println(res.Error)
				}
			}
			//将点赞总数据同步到视频表中
			sum, err := RDB.SCard(CTX, "video:"+videoID+":favorite").Result()
			//fmt.Println("点赞总数：", sum)
			if err != nil {
				fmt.Println(err)
			}
			DB.Model(&Video{}).Where("id = ?", vid).Update("favorite_count", sum)
		}
	}

	//遍历取消点赞视频id
	keys, _ = RDB.Keys(CTX, "video:*:unfavorite").Result()
	videoIDs = make([]string, 0)
	for _, key := range keys {
		// 将命名空间和 ":unfavorite" 部分去除，只保留视频 ID
		videoID := strings.TrimPrefix(key, "video:")
		videoID = strings.TrimSuffix(videoID, ":unfavorite")
		videoIDs = append(videoIDs, videoID)
	}
	//fmt.Println("取消点赞视频id：", videoIDs)
	//写入mysql
	for _, videoID := range videoIDs {
		key := "video:" + videoID + ":unfavorite"
		userIDs, _ := RDB.SMembers(CTX, key).Result()
		//fmt.Println("取消点赞用户id：", userIDs)
		//遍历取消点赞用户id
		for _, userID := range userIDs {
			uid, _ := strconv.ParseUint(userID, 10, 64)
			vid, _ := strconv.ParseUint(videoID, 10, 64)
			//将取消点赞数据同步增加到取消点赞表中
			favorite := Favorite{
				UserID:  uint(uid),
				VideoID: uint(vid),
			}
			//判断是否存在，存在则删除
			res := DB.First(&favorite, "user_id = ? AND video_id = ?", uid, vid)
			if res.Error != nil && res.Error.Error() != "record not found" {
				res = DB.Where("user_id = ? AND video_id = ?", uid, vid).Delete(&favorite)
				if res.Error != nil {
					fmt.Println(res.Error)
				}
			}
			//将点赞总数据同步到视频表中
			sum, err := RDB.SCard(CTX, "video:"+videoID+":favorite").Result()
			//fmt.Println("点赞总数：", sum)
			if err != nil {
				fmt.Println(err)
			}
			DB.Model(&Video{}).Where("id = ?", vid).Update("favorite_count", sum)
		}
	}

	//同步用户关注数据
	keys, _ = RDB.Keys(CTX, "*fo").Result()
	for _, key := range keys {
		// 将命名空间和 "fo" 部分去除，只保留用户 ID
		userID := strings.TrimSuffix(key, "fo")
		//关注的用户的数量
		count, _ := RDB.SCard(CTX, key).Result()

		//将关注数据同步增加到用户表中
		uid, _ := strconv.ParseUint(userID, 10, 64)
		DB.Model(&User{}).Where("id = ?", uid).Update("follow_count", count)
	}

	//同步用户粉丝数据
	keys, _ = RDB.Keys(CTX, "*ed").Result()
	for _, key := range keys {
		// 将命名空间和 "ed" 部分去除，只保留用户 ID
		userID := strings.TrimSuffix(key, "ed")
		//粉丝的数量
		count, _ := RDB.SCard(CTX, key).Result()

		//将粉丝数据同步增加到用户表中
		uid, _ := strconv.ParseUint(userID, 10, 64)
		DB.Model(&User{}).Where("id = ?", uid).Update("follower_count", count)
	}

	//同步用户关注数据到关注表中
	keys, _ = RDB.Keys(CTX, "*fo").Result()
	//清除关注表中的记录
	DB.Exec("TRUNCATE TABLE follow")
	for _, key := range keys {
		// 将命名空间和 "fo" 部分去除，只保留用户 ID
		userID := strings.TrimSuffix(key, "fo")
		//遍历关注的用户id
		userIDs, _ := RDB.SMembers(CTX, key).Result()
		for _, id := range userIDs {
			uid, _ := strconv.ParseUint(userID, 10, 64)
			authorID, _ := strconv.ParseUint(id, 10, 64)
			//将关注数据同步增加到关注表中
			follow := Follow{
				UserID:   uint(uid),
				AuthorID: uint(authorID),
			}
			//创建关注记录
			res := DB.Create(&follow)
			if res.Error != nil {
				fmt.Println(res.Error)
			}
		}

	}

	i++
	fmt.Println("同步次数", i)
}
