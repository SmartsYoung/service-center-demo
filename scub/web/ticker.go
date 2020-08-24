package main

import (
	"fmt"
	"time"
)

func main() {

	// 初始化通道
	ch11 := make(chan int, 1000)
	sign := make(chan byte, 1)

	// 给 ch11 通道写入数据
	go func() {
		for i := 0; i < 100; i++ {
			// go func ()  {
			//     time.Sleep(1e9)
			//         ch11 <- i
			// }()
			time.Sleep(1e9)
			ch11 <- i
		}
	}()
	// 单独起一个 Goroutine 执行 select
	//go func(){
	var e int
	ok := true
	// 首先声明一个 * time.Timer 类型的值，然后在相关 case 之后声明的匿名函数中尽可能的复用它
	var timer *time.Timer
	timer = time.NewTimer(10 * time.Second)
	//expire := <- timer.C
	//fmt.Printf("Expiration time: %v.\n", expire)
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("Timeout.")
				ok = false
				break
			case e = <-ch11:
				fmt.Printf("ch11 -> %d\n", e)
			}
			// 终止 for 循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()

	//}()

	// 惯用手法，读取 sign 通道数据，为了等待 select 的 Goroutine 执行。
	flag := <-sign
	fmt.Println(flag)
}
