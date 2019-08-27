来自于对 [《跟煎鱼学GO》]（https://eddycjy.gitbook.io/golang/di-3-ke-gin/install）整个系列教程的学习，并简单记录知识点以及学习过程中遇到的问题
感谢作者的技术输出

项目知识点
==========
## 目录结构
```
Gin-blog-example/
├── conf                    //用于存储配置文件
├── data                    //MySQL数据挂载卷目录 
├── docs                    //swagger 生成的 api 文档相关信息
├── logs                    //生成日志
├── md                      //笔记markdown    
├── middleware              //应用中间件
├── models                  //引用数据库模型
├── pkg                     //第三方包
├── routers                 //路由处理逻辑
├── runtime                 //应用运行时数据
└── service                 //抽离封装 api 业务逻辑，进行缓存，DB操作包装

```

## 本项目所引用的第三方库
```
    github.com/golang/dep/cmd/dep   //官方包依赖管理工具
    github.com/go-ini/in            //读取 .ini 配置文件
    github.com/Unknwon/com          //Unknwon的工具库，包含了常用的一些封装
    github.com/go-sql-driver/mysql  // MySQL 驱动包
    github.com/jinzhu/gorm          // go 中实现数据库访问ORM（对象关系映射）方便利用面向对象的方法对数据库进行CRUD
    github.com/astaxie/beego/validation     //beego的表单验证库  
    github.com/fvbock/endless       //实现 Golang HTTP/HTTPS 服务重新启动的零停机
    github.com/robfig/cron          //定时任务
    github.com/gomodule/redigo/redis    //redis go
    github.com/360EntSecGroup-Skylar/excelize/v2    // xlsx 操作库
    github.com/boombuler/barcode    // 二维码
```


## cURL
> <u>curl</u> (Command Line URL viewer) 是一种命令行工具，作用是发出网络请求，然后得到和提取数据，显示在"标准输出"（stdout）上面。
> 它支持多种协议  
> [curl命令详解](https://www.jianshu.com/p/07c4dddae43a)


## app.ini 和 go-ini
> 将项目的一些配置从代码中提取出来，放到 .ini 文件中然后去读取  
> 拉取go-ini/ini的依赖包  
> ``go get -u github.com/go-ini/ini``


## GORM
> [GORM 中文文档](http://gorm.book.jasperxu.com/)  
> Golang 写的，对开发人员友好的ORM库，
> GORM 可以通过自定义 callbacks 自定义回调方法注册进钩子方法中，避免多个文件去进行重复回调的书写

## beego/validation
> [表单验证](https://beego.me/docs/mvc/controller/validation.md)  
> beego 也是一个快速开发 go 应用的 HTTP 框架。这里使用了该框架的表单验证模块
用于数据验证和错误收集的模块。

## JWT（Json Web token）
> JWT 是一个很长的字符串，中间用 '.' 分隔成三个部分。这三个部分依次如下
```
 * Header(头部）
 * Payload(负载)
 * Signature(签名)
 
 写成一行就是这个样子
 Header.Payload.Signature
 
 再看个实际的例子
   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```
> Header部分是一个JSON对象，描述JWT的元数据，通常如下
>> ```
>> {
>>  "alg": "HS256",
>>  "typ": "JWT"
>> }
>>```  
> alg属性表示签名的算法（algorithm），默认是 HMAC SHA256（写成 HS256）；  
> typ属性表示这个令牌（token）的类型（type），JWT 令牌统一写为JWT。      
> 最后，将上面的 JSON 对象使用 Base64URL 算法（详见后文）转成字符串。  

---
> Payload 部分也是一个JSON对象，用来存放实际需要传递的数据，JWT规定了7个官方字段，提供选用  
>> ```
>> iss (issuer)：签发人
>> exp (expiration time)：过期时间
>> sub (subject)：主题
>> aud (audience)：受众
>> nbf (Not Before)：生效时间
>> iat (Issued At)：签发时间
>> jti (JWT ID)：编号
>> ```
> 除了官方字段，还可以在这个部分定义私有字段，例如
>> ```
>> {
>>     "sub":"1111",
>>     "admin":true
>> }
>> ```
*值得注意的是JWT默认是不加密的，任何人都可以读到，所以不要把秘密信息放到这个部分，最后 json 也要用 Base64URL 转成字符串
---
> Signature 部分是对前两个部分的签名，防止数据被篡改  
> 首先，需要指定一个密钥（secret）。这个密钥只有服务器才知道，不能泄露给用户。然后，使用 Header 里面指定的签名算法（默认是 HMAC SHA256），按照下面的公式产生签名。
>> ```
>> HMACSHA256(
>>   base64UrlEncode(header) + "." +
>>   base64UrlEncode(payload),
>>   secret)
>> ```
> 算出签名以后，把 Header、Payload、Signature 三个部分拼成一个字符串，每个部分之间用"点"（.）分隔，就可以返回给用户。
>
> 下面的文章讲的不错  
>> *[JSON Web Token 在 web 应用间安全的传递信息](http://blog.didispace.com/json-web-token-web-security/)*  
>> *[八幅漫画理解使用 JWT设计的单点登录系统](http://blog.didispace.com/user-authentication-with-jwt/)*   
>> *[理解JWT的使用场景和优劣](http://blog.didispace.com/learn-how-to-use-jwt-xjf/)*

## endless
> [如何优雅地重启go程序--endless篇](https://blog.csdn.net/tomatomas/article/details/94839857)
> endless 通过监听信号量，完成对服务器管控一系列操作，达到服务重新启动的零停机效果 

## swagger
> [使用swaggo自动生成Restful API文档](https://ieevee.com/tech/2018/04/19/go-swag.html)  
> 1、http://xxxx:xxxx/swagger/index.html 访问如果出现“404 page not found”。需要在routers.go中加下路由配置，这个连载文章中没有提到  
> 2、路由配置重启后，刷新能访问，但是界面却出现 ”Failed to load spec.“ 是因为没有将swag init初始化生成的文件夹进行导包。所以加载失败 import 中加入下 _ "you_project_name/docs"

## docker 
> 创建 Dockerfile 文件定义 Docker 镜像生成流程。文件内容是一条一条指令。每一条指令构建一层。
> 这些指令应用于基础镜像并最终创建一个新的镜像  

> 我们的 go 项目需要特别注意的一点：要使用 [github.com/golang/dep/cmd/dep] 来管理项目的第三方库依赖
> 否则的话在进行 [docker build -t xxx .]  的时候会出现如下类似错误
>> ```
>> pkg/util/pagination.go:5:2: cannot find package "github.com/Unknwon/com" in any of:
>>         /usr/local/go/src/github.com/Unknwon/com (from $GOROOT)
>>         /go/src/github.com/Unknwon/com (from $GOPATH)
>> ```
> 这个原因很简单，我们使用 docker 构建的项目基于网络拉取的 golang 镜像的 `$GOROOT` 和  
> `$GOPATH`中并没有我们实际本机中`$GOPATH`里面下载的第三方库。并且`Dockerfile`仅仅能  
> 识别出当前文件根目录下的文件和文件夹。这个 `dep` 的作用主要就是将自己依赖 `$GOPATH`   
> 的第三方库 copy  到自己项目的根目录新建文件夹 'vendor' 下。仅仅供自己使用，如此一来  
> 构建镜像的时候所有依赖的完整代码都有了就可以构建通过了
> ```
> dep init -gopath -v
> ```
> 该命令会先从`$GOPATH`查找既有的依赖包，若不存在则从对应网络资源处下载
>> [Go依赖管理工具dep](https://eddycjy.gitbook.io/golang/di-2-ke-bao-guan-li/dep)  
>> [docker构建golang分布式带依赖库项目镜像](https://blog.csdn.net/u012740992/article/details/91841021)
>    
> ---   
> #### 拉取 mysql 容器
>> ```
>> docker pull mysql
>> ```
> #### MySQL挂在数据卷
>> 首先创建一个目录用于存放数据卷
>> 这里我选在在当前项目下创建
>> ```
>> mkdir data
>> cd data
>> mkdir docker-mysql
>> cd docker-mysql
>> pwd  //获取全路径
>> //启动时候将常见好的挂载目录绑定
>> docker run --name blog-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -v /Users/coulson/go/src/Gin-blog-example/data/docker-mysql:/var/lib/mysql -d mysql:5.7
> 创建成功后，观察当前项目下的 /data/docker-mysql，多了不少数据库文件
>  
> #### 修改 conf/app.ini 中的 [database] 下的 HOST 为上面 docker 启动所配置的名称`blog-mysql`和密码`123456`
>> ```
>> ...
>> 
>> [database]
>> TYPE = mysql
>> USER = root
>> PASSWORD = 123456
>> HOST = blog-mysql:3306
>> NAME = blog
>> TABLE_PREFIX = blog_
>> ...
>> 
> #### 编译可执行文件
>> ```
>> CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o Gin-blog-example .
> #### 配置完成构建我们的镜像
>> ```
>> docker build -t gin-blog-docker-scratch .  
>> docker images //查看我们创建好的项目镜像
>> #gin-blog-docker-scratch   latest              16d0501e7a3f        3 minutes ago       398MB
> #### golang+mysql 将我们的项目和MySQL关联起来
>> ```
>> // 通过 --link  可以在单机容器内直接使用其关联的容器别名进行访问，而不通过ip。
>> docker run --link blog-mysql:mysql -p 8000:8000 gin-blog-docker-scratch


## cron
> [cron表达式详解](https://www.cnblogs.com/linjiqin/p/3178452.html)  
> cron 可以简单理解遵循一串定义的字符规则，来描述定时任务

## redis
> 是一种使用内存存储(in-memory)的非关系数据库  
> [Redis实战](http://redisinaction.com/index.html)  
> [Redis 命令参考](http://doc.redisfans.com/index.html#)  
> [Redis 设计与实现](http://redisbook.com/)