FROM golang:latest
WORKDIR $GOPATH/src/github.com/SmartsYoung/service-center-demo/gin-blog
COPY . $GOPATH/src/github.com/SmartsYoung/service-center-demo/gin-blog
ENV GOPROXY https://goproxy.cn
RUN go build .
EXPOSE 8000
ENTRYPOINT ["./gin-blog"]
