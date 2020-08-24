package main

import (
    "errors"
    "flag"
    "github.com/apache/servicecomb-service-center/pkg/util"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "strings"
    "time"
)


var (
    addr = flag.String("addr", "localhost:30100", "http service address")

    filename  string
    upgrader  = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
)



var closeCh = make(chan struct{})

type watcherConn struct {
}

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


func main() {




    conn, _, err := websocket.DefaultDialer.Dial(
        strings.Replace(s.Addr, "http://", "ws://", 1), nil)
    if err != nil {
        panic("Dial: " + err.Error())
    }


    EstablishWebSocketError(conn, errors.New("error"))
}

func EstablishWebSocketError(conn *websocket.Conn, err error) {
    remoteAddr := conn.RemoteAddr().String()
    log.Fatal(err, "establish[%s] websocket watch failed.", remoteAddr)
    if err := conn.WriteMessage(websocket.TextMessage, util.StringToBytesWithNoCopy(err.Error())); err != nil {
        log.Fatal(err, "establish[%s] websocket watch failed: write message failed.", remoteAddr)
    }
}