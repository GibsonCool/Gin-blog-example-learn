# 使用最新版 golang 作为基础镜像
FROM golang:latest

RUN go version

# 设置工作目录，没有则自动新建
WORKDIR $GOPATH/src/Gin-blog-example
# 拷贝代码到当前
COPY . $GOPATH/src/Gin-blog-example

##将dep相关内容COPY到/go/src
#COPY ./vendor ./vendor
#COPY Gopkg.lock .
#COPY Gopkg.toml .

RUN go build .


EXPOSE 8000
ENTRYPOINT ["./Gin-blog-example"]