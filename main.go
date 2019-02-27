package main

import (
	"m3u8-golang/Tool"
	"fmt"
	"os"
	"strings"
	"time"
	"bufio"
	"io"
)

var tool = Tool.NewTool()

type tsStruct struct {
	index int
	url   string
}

func main() {
	err := os.MkdirAll("./download/ts", os.ModePerm)
	if err != nil {
		fmt.Println("创建目录失败：", err)
		return
	}

	f, err := os.Open("./list.txt")
	if err != nil {
		fmt.Println("读取list.txt文件失败 ：", err)
		return
	}

	defer f.Close()

	videoIndex := 1
	exitChanLen := 100
	rd := bufio.NewReader(f)
	for {
		url, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}

		//TODO 校验url格式
		//获取真实m3u8地址
		m3u8 := strings.Replace(url, "index.m3u8", "1000k/hls/index.m3u8", -1)
		//去除换行符
		m3u81 := strings.Replace(m3u8, "\n", "", -1)
		tsFileUrls := strings.Replace(m3u81, "\r", "", -1)

		//发出请求
		str, err := tool.Get(tsFileUrls)
		if err != nil {
			fmt.Println("获取tsFileUrls失败：", err)
			return
		}

		tsFiles := strings.Split(string(str), "\n")

		var maxIndex int
		tsChan := make(chan tsStruct, 100)
		exitChan := make(chan bool, exitChanLen)
		go writeTsChan(tsFileUrls, tsFiles, tsChan, &maxIndex)
		//开启读的协程数
		time.Sleep(time.Millisecond * 100)
		for i := 0; i < exitChanLen; i++ {
			go readTsChan(tsChan, exitChan)
		}

		//主线程等待
		//确保所有任务都已完成
		for i := 0; i < exitChanLen; i++ {
			<-exitChan
		}

		fmt.Println("下载完成，总共下载文件数：", maxIndex)
		//下载完成，合并ts
		tool.Merge("./download/"+fmt.Sprintf("%d", videoIndex)+".mp4", maxIndex)
		videoIndex++
	}
}

func writeTsChan(tsFileUrls string, tsFiles []string, tsChan chan tsStruct, maxIndex *int) {
	i := 0
	for _, v := range tsFiles {
		if strings.HasSuffix(v, ".ts") {
			tmp := tsStruct{i, strings.Replace(tsFileUrls, "index.m3u8", v, -1)}
			tsChan <- tmp
			fmt.Println("添加任务：", i)
			i++
		}
	}

	close(tsChan)
	*maxIndex = i
	fmt.Println("所有任务已添加完：", *maxIndex)
}

func readTsChan(tsChan chan tsStruct, exitChan chan bool) {
	for {
		v, ok := <-tsChan
		if !ok {
			break
		}

		fmt.Println("下载任务：", v.index)
		tsFileName := "./download/ts/" + fmt.Sprintf("%d", v.index) + ".ts"
		tool.DownTs(v.url, tsFileName)
	}

	//读取完成
	exitChan <- true
}
