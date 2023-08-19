package main

import (
	"fmt"
	"github.com/easy-TikTok/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		var usermaster controller.Usermaster
		result := controller.DB.Where("token=?", token).First(&usermaster)
		if result.Error != nil {
			fmt.Println("鉴权失败")
			c.AbortWithStatusJSON(http.StatusOK, controller.Response{StatusCode: 1, StatusMsg: "token is error!"})
		} else {
			fmt.Println("鉴权成功")
			//通过中间件以实现鉴权
		}
	}
}

func main() {
	controller.Init()

	a := gin.Default()
	//静态服务器的建立
	a.Static("/public", "./public")

	//基本功能接口
	c := a.Group("/douyin")
	c.GET("/feed/", controller.Feed)
	c.POST("/user/register/", controller.Register)
	c.POST("/user/login/", controller.Login)
	c.POST("/publish/action/", controller.Publish)

	c.Use(Token()) //全局方式下使用中间件，上面三个路由不会调用
	c.GET("/user", controller.User_info)
	c.GET("/publish/list/", controller.PublishList)
	//互动接口
	c.POST("favorite/action/", controller.FavoriteAction)
	c.GET("/favorite/list/", controller.FavoriteList)
	c.POST("/comment/action/", controller.CommentAction)
	c.GET("/comment/list/", controller.CommentList)
	//社交接口
	c.POST("/relation/action/", controller.RelationAction)
	c.GET("/relation/follow/list/", controller.FollowList)
	c.GET("/relation/follower/list/", controller.FollowerList)
	c.GET("/relation/friend/list/", controller.FriendList)
	c.GET("/message/chat/", controller.MessageChat)
	c.POST("/message/action/", controller.MessageAction) //功能已测试
	controller.Sync()
	a.Run()

}
