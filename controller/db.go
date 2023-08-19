package controller

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

const (
	HOST     = "127.0.0.1"
	PORT     = "3306"
	USER     = "root"
	PASS     = "admin"
	DATABASE = "test"
)

func Init() {
	//dsn := "root:FyOATQYp@tcp(172.16.32.75:51060)/temp?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DATABASE)
	var err error
	fmt.Println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:    true,                                       //开启预编译语句缓存，提高性能
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, //设置生成表的时候表名不为复数
	})
	if err != nil {
		fmt.Println("连接数据库失败")
	}
	//数据库表结构生成
	err = DB.AutoMigrate(&Usermaster{}, &User{}, &Video{}, &Comment{}, &Message{}, &Favorite{}, &Follow{})
	if err != nil {
		fmt.Println("数据库表结构生成失败")
	}

	Redisinit()

	DB.Create(&Usermaster{
		ID:       1,
		Username: "admin",
		Password: "admin123",
		Token:    "admin",
	})
	DB.Create(&User{
		ID:              1,
		Avatar:          "http://192.168.0.102:8080/public/head_1.jpg",
		BackgroundImage: "http://192.168.0.102:8080/public/bg_1.jpg",
		Name:            "admin",
		Signature:       "this is a test text",
		WorkCount:       1,
		UsermasterID:    1,
	})
	DB.Create(&Video{
		ID:       1,
		Title:    "test_video",
		CoverURL: "http://192.168.0.102:8080/public/video_1.jpg",
		PlayURL:  "http://192.168.0.102:8080/public/video_1.mp4",
		UserID:   1,
		Time:     time.Now(),
	})
	RDB.SAdd(CTX, "author1", "1")
}
