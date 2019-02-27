package main

import (
	"fmt"
	"time"
)

//协程demo
func main() {
	//开启读-协程个数
	chanNum := 10
	//放置需要处理数据的chan（读写）
	c := make(chan int, 100)
	//标识是否已处理完并退出主线程
	chanExit := make(chan int, 1)
	go write(c)
	time.Sleep(time.Second)
	for i := 1; i <= chanNum; i++ {
		fmt.Println("启动协程~：", i)
		go read(c, chanExit, i)
	}

	//主线程等待结束(chan阻塞)
	for {
		_, ok := <-chanExit
		if ok {
			break
		}
	}

	fmt.Println("main执行完成")
}

func write(c chan int) {
	for i := 0; i < 10000; i++ {
		fmt.Println("写入数据：", i)
		c <- i
	}

	fmt.Println("写入数据完成")
	close(c)
}

func read(c chan int, chanExit chan int, chanNum int) {
	for {
		v, ok := <-c
		fmt.Printf("chan 取值：%d，协程：%d \n", v, chanNum)
		if !ok {
			break
		}
	}

	fmt.Println("协程读取数据完成:", chanNum)
	//读取完成
	chanExit <- chanNum
}
