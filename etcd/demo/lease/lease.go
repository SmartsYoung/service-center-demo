package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.etcd.io/etcd/clientv3"
)

func main() {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: time.Second * 5,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("connect to etcd success.")
    defer cli.Close()

    // 创建一个5秒的租约
    resp, err := cli.Grant(context.TODO(), 5)
    if err != nil {
        log.Fatal(err)
    }

    // 5秒钟之后, /nazha/ 这个key就会被移除
    _, err = cli.Put(context.TODO(), "/nazha/", "dsb", clientv3.WithLease(resp.ID))
    if err != nil {
        log.Fatal(err)
    }
}
