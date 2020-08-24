package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/websocket"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "reflect"

    "log"
    "time"
)

type watcherConn struct {
}

var closeCh = make(chan struct{})

func (h *watcherConn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var upgrader = websocket.Upgrader{}
    conn, _ := upgrader.Upgrade(w, r, nil)
    for {
        conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
        conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
        _, _, err := conn.ReadMessage()
        if err != nil {
            return
        }
        <-closeCh
        conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
        conn.Close()
        return
    }
}

var listwatch = "ws://127.0.0.1:30100/v4/default/registry/microservices/7062417bf9ebd4c646bb23059003cea42180894a/listwatcher"
var gowatcher = "/v4/default/registry/microservices/090c0124ce3b3f29bdec717fdb3adf34b782e25d/watcher"
var testwatcher = "/v4/default/registry/microservices/a4c50161eae3e1c7a40201f675a24401ff835af7/watcher"
var (
    addr = flag.String("addr", "localhost:30100", "http service address")
    path = flag.String("path", "/v4/default/registry/microservices/090c0124ce3b3f29bdec717fdb3adf34b782e25d/watcher", "websocket path")
)

func main() {
    flag.Parse()
    log.SetFlags(0)

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
    log.Printf("connecting to %s", u.String())
    //向服务器发送连接请求，websocket 统一使用 ws://，
    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if nil != err {
        log.Println(err)
        return
    }
    fmt.Println(reflect.TypeOf(conn))

    defer conn.Close()

    for {
        log.Printf("waiter")
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            return
        }
        log.Printf("recv: %s", message)
    }
}


// event[{"action":"CREATE","key":{"tenant":"default/default","appId":"default","serviceName":"myserver1","version":"0.0.1"},"instance":{"App":"","ServiceName":"","instanceId":"ef56a1b8e0a411eab38cfa163e00c6b9",
// "serviceId":"6f794dc34b1ae96274337baa2e408c3222282448","endpoints":["localhost:808"],"hostName":"insdfsdsdfsdff2233","status":"UP","properties":{"Name":"12"},"healthCheck":{"mode":"push","interval":30,"times":3},
// "timestamp":"1597680958","modTimestamp":"1597680958","version":"0.0.1"}}]

/*
{
  "action": "CREATE",
  "key": {
    "appId": "default",
    "serviceName": "hi-ser",
    "version": "0.0.3"
  },
  "instance": {
    "instanceId": "string",
    "serviceId": "string",
    "version": "0.0.3",
    "hostName": "yc",
    "endpoints": [
      "localhost:9290"
    ],
    "status": "UP",
    "properties": {
      "Age": "12"
    },
    "healthCheck": {
      "mode": "push",
      "interval": 30,
      "times": 3
    },
    "timestamp": "string",
    "modTimestamp": "string"
  }
}
*/