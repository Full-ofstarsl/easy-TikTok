package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_id := c.Query("to_user_id")
	action_type, _ := strconv.Atoi(c.Query("action_type"))
	var usermaster Usermaster
	DB.Where("token=?", token).First(&usermaster)
	//执行逻辑为将用户关注的用户和被关注的用户写成两个集合
	if action_type == 1 {
		RDB.SAdd(CTX, fmt.Sprintf("%sfo", strconv.FormatInt(usermaster.ID, 10)), to_id)

		RDB.SAdd(CTX, fmt.Sprintf("%sed", to_id), usermaster.ID) //用户被关注列表即粉丝列表
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "add success",
		})
	} else if action_type == 2 {
		RDB.SRem(CTX, fmt.Sprintf("%sfo", string(usermaster.ID)), to_id)
		RDB.SRem(CTX, fmt.Sprintf("%sed", to_id), usermaster.ID) //用户被关注列表
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "del success",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "fail",
		})
	}

}
