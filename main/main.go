package main

import (
	"Tool"
	"fmt"
	"os"
	"strings"
)

func main() {

	err := os.MkdirAll("./download/ts", os.ModePerm)
	if err != nil {
		fmt.Println("创建目录失败：", err)
		return
	}

	url := "https://ifeng.com-l-ifeng.com/20190222/27878_b3b9ee0d/index.m3u8"
	//获取真实m3u8地址
	tsFileUrls := strings.Replace(url, "index.m3u8", "1000k/hls/index.m3u8", -1)
	//发出请求
	tool := Tool.NewTool()
	str, err := tool.Get(tsFileUrls)
	if err != nil {
		fmt.Println("获取tsFileUrls失败：", err)
		return
	}

	tsFiles := strings.Split(string(str), "\n")
	i := 0
	for _, v := range tsFiles {
		if strings.HasSuffix(v, ".ts") {
			tsUrl := strings.Replace(tsFileUrls, "index.m3u8", v, -1)
			tsFileName := "./download/ts/" + fmt.Sprintf("%d", i) + ".ts"

			//TODO goroutine 开启多协程
			tool.DownTs(tsUrl, tsFileName)
		}
	}

	tool.Merge("./download/0.mp4")
}
