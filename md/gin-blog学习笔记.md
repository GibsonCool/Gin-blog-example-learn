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

