# 一、项目介绍
### 此项目为简易版抖音后端，如需食用，还需配合字节跳动青训营官方前端文件使用

> 项目采用gin框架+GORM+redis技术实现
>
> <https://bc1213a8022d407d2e25bf4ba4466203-app.1024paas.com>
>
> <https://1024code.com/codecubes/ajqash7>

# 二、项目分工

| **团队成员** |                            **主要贡献**                            |
| :------: | :------------------------------------------------------------: |
|    杨越    | 在项目初期，对整体架构进行设计，包括选择合适的模块和框架作为核心，同时规划数据库结构，以及构建视频 feed 流和社交接口。 |
|    陈培煜   |       参与项目初期架构设计，评论以及点赞等主要功能的设计与测试，在项目推进过程中积极讨论并解决所遇到的问题。      |
|    王雪    |             投稿接口、评论接口等主要功能的设计与测试，后期数据的收集以及项目文档的撰写。             |

# 三、演示视频

暂时无法在飞书文档外展示此内容

# 四、项目实现

### 3.1 技术选型与相关开发文档

### 3.1.1 运用gorm建表

> gorm官方文档链接：[链接](https://gorm.io/zh_CN/docs/index.html)

Gorm是一个强大的ORM（对象关系映射）库，它提供了一种以对象的方式进行数据库操作的方法，在本次开发中通过使用gorm建表简化了数据库操作，提高其开发效率，实现了数据库连接、数据表迁移、数据创建和其他操作。

在数据库结构设计中，我们可以使用 GORM 库的 DB.AutoMigrate() 函数来将结构体转换为数据库中各个表的结构，并生成一些默认数据以便使用。同时，通过使用该库，我们还可以实现数据库中各个表之间的关系，无论是一对一、一对多还是多对多关系，只需要让表之间互相持有对方的切片或外键名称字段即可实现。

具体如下：

    //使用其GORM库，用于与数据库进行交互并生成表结构
    import "gorm.io/gorm"
    gorm.Open(mysql.Open(dsn), &amp;gorm.Config{})
    err = DB.AutoMigrate(&Usermaster{}, &User{}, &Video{}, &Comment{}, &Message{}, &Favorite{}, &Follow{})
    if err != nil {
        fmt.Println("数据库表结构生成失败")
    }



    //在结构体中定义即互相持有实现关系的建立，多余的字段可采用（json:"-"）在序列化时屏蔽。
    type User struct {
    ...
       FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` 
    ...
    }
    type Video struct {
    ...
       Users         []Usermaster `json:"-" gorm:"many2many:user_favor_videos;"`
    ...
    }

DB.where(...).First(...) DB.Preload(...).Find(...) //预编译以实现返回数据的嵌套查询 使用ORM技术，在关系型数据库和对象之间建立了映射关系，使得可以通过面向对象的方式操作数据库

### 3.1.2 运用Redis做数据缓存

抖音业务的性能关键在于关注和点赞，过多的数据库访问可能导致数据库崩溃。Redis 作为一种内存数据库，作为应用程序和后端数据存储（如数据库）之间的缓存层，可以显著提升数据读取速度，降低后端数据库的负载压力。将常用数据存储在 Redis 中，可加速视频点赞和用户关注等操作的处理。

**(为了防止缓存穿透和缓存雪崩等安全问题，我们将关注和点赞等热点数据存储在 Redis 中，并完全通过 Redis 进行读写操作。同时，我们设计了一个函数，用于定期将 Redis 中的数据同步到数据库，从而实现数据的持久化。)**

如下：

    //采用Redis数据库并利用RDB.SAdd和RDB.SRem方法来向 Redis 集合中添加或移除元素.
    //利用集合的定义来实现一对多以及集合嵌套来实现多对多的使用。
    if req.ActionType == "1" {
    RDB.SAdd(CTX, "video:"+req.VideoID+":favorite", uid)
    RDB.SRem(CTX, "video:"+req.VideoID+":unfavorite", uid)
    } 
    else {
    RDB.SRem(CTX, "video:"+req.VideoID+":favorite", uid)
    RDB.SAdd(CTX, "video:"+req.VideoID+":unfavorite", uid)
    }

通过调用 Redis 的命令，使用了 Redis 数据库进行数据查询。例如，在代码中调用RDB.SMembers(CTX, user\_id).Result() 方法，从 Redis 中获取集合中的所有成员，其中 user\_id 是从 HTTP 请求的查询参数中获取的。

### 3.1.3 中间件

\--> 用户鉴权：将用户中的身份信息加密为一个token，并在每个请求中发送给token进行身份验证，通过 token := c.Query("token") 获取 HTTP 请求中传递的 token 参数，并在 Token() 函数中使用该 token 进行用户鉴权。

Token函数实现了一个Gin中间件，通过检查请求中的token参数进行用户身份鉴权。在该中间件中，首先从请求中获取token参数，然后使用controller.DB对象执行数据库查询操作，根据token查询用户记录。如果查询结果中存在错误或用户记录不存在，则输出鉴权失败的消息并返回相应的错误响应并且中断接下来的动作。如果用户存在，则输出鉴权成功的消息，表示通过鉴权。

\--> 利用哈希函数加密：利用哈希函数进行加密是一种常见的密码存储技术，它将用户的密码转换为一个哈希值，并将该哈希值存储在数据库中，而不是明文存储用户的原始密码。这样做的目的是保护用户密码的安全性，即使数据库泄露，攻击者也无法轻易还原密码。

如下：

    func hashSHA256(input string) string {
            hasher := sha256.New()
            hasher.Write([]byte(input))
            hashedBytes := hasher.Sum(nil)
            hashedString := hex.EncodeToString(hashedBytes)
            return hashedString
    }
    //hashSHA256函数使用了SHA-256哈希算法对输入进行哈希处理，使用Go语言标准库中的crypto/sha256包来创建SHA-256哈希对象，并将输入转换为字节数组后进行哈希计算。最后，将哈希值转换为十六进制字符串表示。

以下是一般的哈希函数加密过程：

1.用户注册或更改密码时，将其输入的原始密码作为输入数据。

2.选择一个哈希函数，如SHA-256、bcrypt或Argon2等，这些函数都是常见的哈希函数算法。

3.将原始密码作为输入，使用选定的哈希函数对其进行哈希运算。哈希函数会对输入进行计算，生成固定长度的哈希值作为输出。

4.将生成的哈希值存储在数据库的密码字段中。

5.当用户登录时，系统会将其输入的密码再次经过相同的哈希函数运算，并将结果与数据库中存储的哈希值进行比对。

6.如果两个哈希值匹配，则表示用户提供的密码正确，允许其登录；否则，密码不匹配，拒绝登录。

### 3.1.4 循环和切片操作（动态生成是否关注、点赞等数据）

通过 for 循环遍历 ids 切片中的每个元素，并根据每个 id 执行数据库查询和操作。使用 append 将查询到的用户对象添加到 userlist 切片中。

如下：

    for _, id := range ids {
        result := DB.Where("usermaster_id=?", id).First(&amp;user)
        if result.Error != nil {
            fmt.Println("查询用户失败")
        }
        userlist = append(userlist, user)
    }

上述代码中的 for 循环用于遍历切片 ids 中的每个元素，并将每个元素赋值给变量 id。在每次循环迭代中，执行了以下操作：

1.通过指定条件（"usermaster\_id=?"）查询数据库。这里使用了 ORM 技术，表示查询符合 usermaster\_id 字段等于当前 id 值的记录。

2.调用 First(\&amp;user) 方法，将查询结果的第一条记录填充到 user 对象中。如果查询出错（即 result.Error != nil），则在控制台打印 "查询用户失败"。否则，将查询到的 user 对象追加到切片 userlist 中。

在这段代码中，还使用了 append 函数将 user 对象追加到 userlist 切片中，实现了动态扩展切片的功能。此外，还利用 range 关键字实现了对切片的迭代操作。

    for i := 0; i &lt; len(videolist); i++ {
        Sum(token, &amp;videolist[i].User)
    }

在这里使用了 for 循环来遍历视频列表 videolist。循环的索引 i 从 0 开始，每次迭代递增1，直到达到len(videolist) 的长度（视频列表的元素个数）为止。

在每次迭代中，使用索引 i 访问视频列表的对应元素 videolist\[i]，然后调用函数 Sum(token, \&amp;videolist\[i].User)动态生成是否已关注、点赞等数据。

### 3.2 架构设计

### 3.2.1 数据库关系

![](https://t18fcm764oh.feishu.cn/space/api/box/stream/download/asynccode/?code=NzZlMWZiMDk1MGU0ZWY5YTg3ZWY2NTYxNmU5M2Y1NThfbXRXaGFMZkc2aEdGaFg1bHJsWTRURE9SVHFISkgyWmlfVG9rZW46S3JBWGJXNU9Ub3kxWnR4cVpoNmMxR2hrbktjXzE2OTI3NTY4Mzk6MTY5Mjc2MDQzOV9WNA)

userinfo：存储用户基本信息

video：存储发布视频信息

discuss：存储评论的基本信息

### 3.2.2 初始化数据库

```go
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
        // dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DATABASE)
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
}
```

### 3.2.3 架构说明

![](https://t18fcm764oh.feishu.cn/space/api/box/stream/download/asynccode/?code=NjVjZDY0MmM3NWY0NjE4ZGQ4NzdhNDJkNTUyN2IxMDVfakY0Nk1icVpSQ1JKWDRqVjBEZzM1cEtjUmFKbUx3dFBfVG9rZW46TUFIWmJFRlRQbzd5Zll4OVFYVGNkMkxqbmtiXzE2OTI3NTY4Mzk6MTY5Mjc2MDQzOV9WNA)

例如用户注册和登录：用户通过客户端（App）进行注册和登录，提交注册信息，并进行身份验证。后端通过用户认证和权限管理模块对用户进行验证，验证通过后用户可以使用极简版抖音服务。

### 3.3 项目代码介绍

![](https://t18fcm764oh.feishu.cn/space/api/box/stream/download/asynccode/?code=NjNkMjkyMzUxY2RmMWU3OWVjMmZjZDNhYjg4NzdhM2JfYWFzRlowb3JTYUJUdVFwYWhic1NhMVdKdDNHVm9NWVZfVG9rZW46Vzk5b2IzbmZrb1kxR3J4ZXhjTWNMZTNibkx5XzE2OTI3NTY4Mzk6MTY5Mjc2MDQzOV9WNA)

comment-action.go 用于处理评论相关操作的请求。

Comment-list.go 用于获取评论列表的请求。

Db.go 用于初始化数据库连接和创建表结构，并插入一些初始数据。

Favorite-action.go 用于处理点赞操作的请求和响应。

Favorite-list.go 用于处理用户获取点赞视频列表的请求和响应。

Feed.go 关于视频流，使用了gin框架处理HTTP请求。

Follow-list.go Follower-list.go Friend-list.go

根据不同的操作结果，接收到相应的状态码、状态消息以及用户列表信息。

Message-action.go 根据客户端传递的参数进行判断和操作。如果操作类型为发送消息，则将消息存储到数据库中，并返回成功的状态码和状态消息。如果操作类型不符合预期，则返回失败的状态码和状态消息

Message-chat.go 根据客户端传递的参数进行查询操作，并将查询结果以 JSON 格式返回给客户端。如果查询成功并且消息列表不为空，返回状态码为 "0" 和消息列表。如果查询失败或消息列表为空，返回相应的状态码和状态消息。

Publish-action.go 用于处理文件上传和保存的请求。

Publishlist.go 用于处理获取用户发布视频列表的请求。

Redis.go 用于初始化 Redis 客户端并设置上下文（context）的函数。

Relation-action.go 用于处理关系操作的函数 RelationAction，针对给定的参数执行不同的逻辑。

Struct.go 定义了一系列数据结构和它们之间的关系，用于建模应用程序中的实体和关联。

Sync.go 定义数据结构和类型，用于建模应用程序中的实体和关联。

User-info.go 用于处理用户信息请求。

User-login.go 用于处理用户登录请求。

User-register.go 用于处理用户注册请求。

Main.go 主程序，是一个基于Gin框架实现的Web服务入口。

# 五、测试结果(压力测试均为10线程10循环共计100次测试)

1.  基础接口压力测试

![](https://t18fcm764oh.feishu.cn/space/api/box/stream/download/asynccode/?code=OGE5N2NlODgwZTNkOTY2ZWRhMTQxNDhjMWRkNzk0ZjdfNDdCallDNWZPSEtxV25rcTY0WlRWNXUzQVBiRjNJazlfVG9rZW46TzFXYmJCcUh6b0RrbTh4ejdrbGNNWkVsbmNmXzE2OTI3NTY4Mzk6MTY5Mjc2MDQzOV9WNA)

1.  互动接口压力测试

![](https://t18fcm764oh.feishu.cn/space/api/box/stream/download/asynccode/?code=OTA5ZDYwYzg3ODkyODM3ZDMxNzk4NjliNGVhNjI2MzdfSlJaYWJLRkViVFVNM2tGUVJOcVBZUUVaR0xJQXBnRXJfVG9rZW46Vk9CcmJwR3lOb0JsdDl4S0phcmM5M2dQblBSXzE2OTI3NTY4Mzk6MTY5Mjc2MDQzOV9WNA)

1.  社交接口压力测试

![](https://t18fcm764oh.feishu.cn/space/api/box/stream/download/asynccode/?code=YmFkYjRhMDFhOGE1NjFlYzAxNWVkNGRiZGQ0ZDViMjhfcmIwY1VJdm95VG1uS3lzSThlTDRVdU5vY1JPYmo5ZENfVG9rZW46QU5BbGI3Ym04bzljck94aU9iaGMxeWVhbnRoXzE2OTI3NTY4Mzk6MTY5Mjc2MDQzOV9WNA)

\*\*优点：\*\*在压力测试过程中，我们观察到系统在承受压力时表现出非常出色的稳定性。得益于 Go 语言底层生态对并发的原生支持，即使面对大量并发用户，系统仍能保持良好的稳定性，同时仅牺牲了较小的时间效率。这充分证明了我们的系统在应对高并发请求时具备强大能力和稳定表现。

\*\*缺点：\*\*经过分析，我们发现在高压力高并发的使用环境中，部分接口存在效率较低的问题。经过深入研究，我们发现这主要与系统整体架构设计和数据库表之间的耦合度较高有关。因此，我们计划在后续进行改良，以提高系统整体的效率和性能。

# 六、项目总结与反思

> 1.  （视频存储方面优化）在投稿过程中，针对视频上传部分，为了更有效地利用存储空间，我们可以考虑采用 TOS（Object Storage）存储方式进行优化。这种存储方式具有分布式、高可靠性、可扩展性强等特点，能够实现数据的高效存储和管理，从而提高整体存储空间的利用率。
> 2.  （整体架构方面优化）当前，整个系统采用的是单体架构设计，这种设计方式使得各个服务之间的耦合性相对较高，可能对系统的可扩展性和可维护性造成一定的影响。为了解决这一问题，我们可以在未来的系统优化中，考虑将单体架构转变为微服务架构，通过将数据进行分库分表的处理方式，进一步提高系统的性能和可维护性，从而更好地满足业务发展的需求。同时，我们还可以对系统进行更加细致的模块划分，降低模块间的耦合度，进一步优化整个系统的架构。
> 3.  由于团队整体开发经验不足，项目初期面临诸多挑战，尤其是在整体架构和数据库结构设计方面存在一定的不合理之处。然而，在后续的过程中，我们团队成员积极进行内部讨论，共同探讨并提出了一系列解决问题的有效方法，逐步优化了项目的开发和实施，也认识到了在真正的实际开发中并不能一人独挡一面，团队之间的协同合作也是一件非常重要的事情。

