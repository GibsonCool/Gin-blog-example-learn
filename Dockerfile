# 使用最新版官方 golang 作为基础镜像
# FROM golang:latest

# Scratch镜像 ，简介小巧基本是个空镜像，我们不在Golang容器中现场编译，只需要一个能够运行执行文件的环境即可
FROM scratch

# 设置工作目录，没有则自动新建
WORKDIR $GOPATH/src/Gin-blog-example
# 拷贝代码到当前
COPY . $GOPATH/src/Gin-blog-example

##将dep相关内容COPY到/go/src
#COPY ./vendor ./vendor
#COPY Gopkg.lock .
#COPY Gopkg.toml .

# 使用 scratch 镜像后这里就进行编译了，而是在构建钱先编译好
# RUN go build .


EXPOSE 8000
ENTRYPOINT ["./Gin-blog-example"]