package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// RegisterReq 用户注册请求
type RegisterReq struct {
	Password string `json:"password"` // 密码，最长32个字符
	Username string `json:"username"` // 注册用户名，最长32个字符
}

// RegisterRes 用户登录返回
type RegisterRes struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

func hashSHA256(input string) string {
	// 创建一个 SHA-256 哈希对象
	hasher := sha256.New()
	// 将输入字符串转换为字节数组，并写入哈希对象
	hasher.Write([]byte(input))
	// 计算哈希值
	hashedBytes := hasher.Sum(nil)
	// 将哈希值转换为十六进制字符串
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}

// 数据库操作
// 插入用户
func InsertUsermaster(user *Usermaster) (int64, error) {
	// 插入数据
	result := DB.Create(user)
	return user.ID, result.Error
}

// 具体操作
func DoRegister(req *RegisterReq) *RegisterRes {
	// 插入用户
	fmt.Println(req.Username, req.Password)
	hashed := hashSHA256(req.Password + req.Username)
	newuser := Usermaster{
		Username: req.Username,
		Password: req.Password,
		Token:    hashed,
	}
	fmt.Println("要插入的用户为：", newuser)

	uid, err := InsertUsermaster(&newuser)

	//插入到user表中
	newuser2 := User{
		Name:         req.Username,
		UsermasterID: uid,
	}
	DB.Create(&newuser2)

	fmt.Println("插入的用户ID为：", uid)
	if err != nil {
		//判断是否是用户名重复
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			// 判断是否是唯一约束冲突
			if mysqlErr.Number == 1062 {
				return &RegisterRes{
					StatusCode: 1,
					StatusMsg:  "注册失败,用户名已存在",
					Token:      "nil",
					UserID:     -1,
				}
			} else {
				return &RegisterRes{
					StatusCode: 1,
					StatusMsg:  "注册失败," + err.Error(),
					Token:      "nil",
					UserID:     -1,
				}
			}
		}
	}
	//返回
	return &RegisterRes{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		Token:      hashed,
		UserID:     uid,
	}
}

func Register(c *gin.Context) {
	var req RegisterReq
	//	c.BindJSON(&req)
	req.Password = c.Query("password")
	req.Username = c.Query("username")

	fmt.Println("发送的消息为", req.Password, req.Username)

	//判断用户名和密码是否过长
	if len(req.Username) > 32 || len(req.Password) > 32 {
		c.JSON(200, RegisterRes{
			StatusCode: 1,
			StatusMsg:  "注册失败,用户名或密码过长,最长32个字符",
			Token:      "nil",
			UserID:     -1,
		})
		return
	}

	//进行注册
	res := DoRegister(&req)
	c.JSON(200, res)
}
