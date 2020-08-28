FROM golang:latest
WORKDIR $GOPATH/src/github.com/SmartsYoung/service-center-demo/gin-blog
COPY . $GOPATH/src/github.com/SmartsYoung/service-center-demo/gin-blog
RUN go build .
EXPOSE 8000
ENTRYPOINT ["./gin-blog"]