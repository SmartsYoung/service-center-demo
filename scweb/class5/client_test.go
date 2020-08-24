package main

import (
    "github.com/gorilla/websocket"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestWatcherConn_ServeHTTP(t *testing.T) {
    s := httptest.NewServer(&watcherConn{})

    conn, _, _ := websocket.DefaultDialer.Dial(
        strings.Replace(s.URL, "http://", "ws://", 1), nil)
    t.Log("conn", conn)
}
