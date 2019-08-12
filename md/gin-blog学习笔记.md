项目知识点
==========
## 目录结构
```
Gin-blog-example/
├── conf                    //用于存储配置文件
├── middleware              //应用中间件
├── models                  //引用数据库模型
├── pkg                     //第三方包
├── routers                 //路由处理逻辑
├── md                      //笔记文档
└── runtime                 //应用运行时数据

```

## 本项目所引用的第三方库
```
    github.com/go-ini/in            //读取 .ini 配置文件
    github.com/Unknwon/com          //Unknwon的工具库，包含了常用的一些封装
    github.com/go-sql-driver/mysql  // MySQL 驱动包
    github.com/jinzhu/gorm          // go 中实现数据库访问ORM（对象关系映射）方便利用面向对象的方法对数据库进行CRUD
    github.com/astaxie/beego/validation     //beego的表单验证库  
    github.com/fvbock/endless       //实现 Golang HTTP/HTTPS 服务重新启动的零停机

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
*值得注意的是JWT默认是不加密的，任何人都可以读到，所以不要把秘密信息放到这个不封，最后json也要用BAse64URL转成字符串
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
> 1、http://127.0.0.1:8000/swagger/index.html 访问如果出现“404 page not found”。需要在routers.go中加下路由配置，这个连载文章中没有提到  
> 2、路由配置重启后，刷新能访问，但是界面却出现 ”Failed to load spec.“ 是因为没有将swag init初始化生成的文件夹进行导包。所以加载失败 import 中加入下 _ "you_project_name/docs"