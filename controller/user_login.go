package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// LoginReq 用户登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRes 用户登录返回
type LoginRes struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

// DoLogin 具体操作
func DoLogin(req *LoginReq) (res *LoginRes) {
	var usermaster Usermaster
	// 通过用户名查询用户
	result := DB.Where("username=?", req.Username).Order("").First(&usermaster)

	//若用户不存在
	if result.Error != nil {
		err := result.Error
		if err.Error() == "record not found" {
			return &LoginRes{
				StatusCode: 1,
				StatusMsg:  "登录失败,用户不存在",
				Token:      "nil",
				UserID:     -1,
			}
		} else {
			return &LoginRes{
				StatusCode: 1,
				StatusMsg:  err.Error(),
				Token:      "nil",
				UserID:     -1,
			}
		}

	}
	//若用户存在,查询密码是否正确
	fmt.Println("查询到的用户为：", usermaster)
	fmt.Println("查询到的用户密码为：", usermaster.Password)
	if usermaster.Password == req.Password {
		return &LoginRes{
			StatusCode: 0,
			StatusMsg:  "登录成功",
			Token:      usermaster.Token,
			UserID:     usermaster.ID,
		}
	} else {
		return &LoginRes{
			StatusCode: 1,
			StatusMsg:  "登录失败,密码错误",
			Token:      "nil",
			UserID:     -1,
		}
	}
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginReq
	//c.BindJSON(&req)
	req.Password = c.Query("password")
	req.Username = c.Query("username")
	fmt.Println("发送的消息为", req.Password, req.Username)
	//进行登录
	res := DoLogin(&req)
	c.JSON(200, res)

}
