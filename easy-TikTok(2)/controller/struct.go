package controller

import "time"

type Response struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// Video
type Video struct {
	User          User      `json:"author"`         // 视频作者信息
	CommentCount  int64     `json:"comment_count"`  // 视频的评论总数
	CoverURL      string    `json:"cover_url"`      // 视频封面地址
	FavoriteCount int       `json:"favorite_count"` // 视频的点赞总数
	ID            uint      `json:"id"`             // 视频唯一标识
	IsFavorite    bool      `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string    `json:"play_url"`       // 视频播放地址
	Title         string    `json:"title"`          // 视频标题
	Time          time.Time `json:"-"`              //作品时间，在返回的时候做一个unix的时间戳转换
	UserID        uint      `json:"-"`              //外键,同时在创建记录的时候可以通过定义外键使得视频和用户绑定起来
	Comment       []Comment `json:"-"`              //与评论是一对多关系
}

// User
type User struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数量
	WorkCount       int    `json:"work_count"`       // 作品数
	//建立一对多关系
	Video   []Video   `gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE; " json:"-"`
	Comment []Comment `gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE; " json:"-"`
	//用于创建一对一关系
	UsermasterID int64 `json:"-"`
}

// Message
type Message struct {
	Content    string `json:"content"`      // 消息内容
	CreateTime int64  `json:"create_time"`  // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID int64  `json:"from_user_id"` // 消息发送者id
	ID         int64  `json:"id"`           // 消息id
	ToUserID   int64  `json:"to_user_id"`   // 消息接收者id

}

// Comment
type Comment struct {
	Content    string `json:"content"`     // 评论内容
	CreateDate int64  `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64  `json:"id"`          // 评论id
	User       User   `json:"user"`        // 评论用户信息,调用预加载
	//建立关系
	Video   Video `json:"-"` //一对多关系
	VideoID uint  `json:"-"` //外键
	UserID  uint  `json:"-"` //外键,同时在创建记录的时候可以通过定义外键使得视频和用户绑定起来
}

// user表,在调用的时候创建新的response响应，使用高级查询
type Usermaster struct {
	ID       int64  `gorm:"primarykey;unique;autoIncrement" json:"user_id"` // 用户id
	Username string `gorm:"unique" json:"-"`                                //用户名
	Password string `json:"-"`                                              //用户密码
	Token    string `json:"token"`                                          // 用户鉴权token
	//用于创建一对一关系
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// 点赞表
type Favorite struct {
	UserID  uint
	VideoID uint
}

// 关注表
type Follow struct {
	UserID   uint
	AuthorID uint
}

//还有各个表之间的关联关系(已解决)
//video表还需要一个字段存储时间戳，修改结构体添加字段，设置json为"-”，然后在调用feed.go的时候直接赋值
//现在还差一个点赞关注列表，创建一个结构体，使用gorm的高级查询重新创建需要返回的response
//点赞关注表分为两个部分，关于用户本身的点赞关注，视频的点赞和评论数
//用户鉴权使用哈希函数进行加密，在前端返回的时候进行id和token的校验
//关于高性能，点赞关注评论三个数据的更新，使用redis来缓存，
