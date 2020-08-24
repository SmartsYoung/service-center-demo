package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strconv"
    "time"
)

var mc *MessageCenter

type Message struct {
    Uid int
    Message string
}

type MessageCenter struct {
    // 测试 没有加读写锁
    messageList []*Message
    userList map[int]chan string
}

func NewMessageCenter() *MessageCenter {
    mc := new(MessageCenter)
    mc.messageList = make([]*Message, 0, 100)
    mc.userList = make(map[int]chan string)
    return mc
}

func (mc *MessageCenter) GetMessage(uid int) []string {
    messages := make([]string, 0, 10)
    for i, msg := range mc.messageList {
        if msg == nil {
            continue
        }
        if msg.Uid == uid {
            messages = append(messages, msg.Message)
            // 临时方案 只是测试用 应更换为list
            mc.messageList[i] = nil
        }
    }
    return messages
}

func (mc *MessageCenter) GetMessageChan(uid int) <- chan string {
    messageChan := make(chan string)
    mc.userList[uid] = messageChan
    return messageChan
}

func (mc *MessageCenter) SendMessage(uid int, message string) {
    messageChan, exist := mc.userList[uid]
    if exist {
        messageChan <- message
        return
    }
    // 未考虑同一账号多登陆情况
    mc.messageList = append(mc.messageList, &Message{uid, message})
}

func (mc *MessageCenter) RemoveUser(uid int) {
    _, exist := mc.userList[uid]
    if exist {
        delete(mc.userList, uid)
    }
}

func IndexServer(w http.ResponseWriter, req *http.Request) {
    http.ServeFile(w, req, "longpoll.html")
}

func SendMessageServer(w http.ResponseWriter, req *http.Request) {
    uid, _ := strconv.Atoi(req.FormValue("uid"))
    message := req.FormValue("message")

    mc.SendMessage(uid, message)

    io.WriteString(w, `{}`)
}

func PollMessageServer(w http.ResponseWriter, req *http.Request) {
    uid, _ := strconv.Atoi(req.FormValue("uid"))

    messages := mc.GetMessage(uid)

    if len(messages) > 0 {
        jsonData, _ := json.Marshal(map[string]interface{}{"status":0, "messages":messages})
        w.Write(jsonData)
        return
    }

    messageChan := mc.GetMessageChan(uid)

    select {
    case message := <- messageChan:
        jsonData, _ := json.Marshal(map[string]interface{}{"status":0, "messages":[]string{message}})
        w.Write(jsonData)
    case <- time.After(10 * time.Second):
        mc.RemoveUser(uid)
        jsonData, _ := json.Marshal(map[string]interface{}{"status":1, "messages":nil})
        n, err := w.Write(jsonData)
        fmt.Println(n, err)
    }
}

func main() {
    fmt.Println("http://127.0.0.1:8080/")

    mc = NewMessageCenter()

    http.HandleFunc("/", IndexServer)
    http.HandleFunc("/sendmessage", SendMessageServer)
    http.HandleFunc("/pollmessage", PollMessageServer)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
