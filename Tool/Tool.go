package Tool

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Tool struct{}

func NewTool() *Tool {
	return &Tool{}
}

//http 请求
func (t *Tool) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求失败： %s \n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("获取结果失败： %s \n", err)
	}

	return body, err
}

//ts流文件下载
func (t *Tool) DownTs(url, name string) {
	fmt.Println("下载中：", name)
	//发出请求
	str, err := t.Get(url)
	if err != nil {
		fmt.Println("下载失败：", name)
	} else {
		if ioutil.WriteFile(name, str, 0644) == nil {
			fmt.Println("下载成功:", name)
		} else {
			fmt.Println("写入失败：", name)
		}
	}
}

//把所有ts合并为mp4
func (t *Tool) Merge(name string) {
	fileObj, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		os.Exit(2)
	}

	defer fileObj.Close()

	for i := 0; i < 3; i++ {
		byte, _ := t.ReadTsFile(i)
		if _, err := fileObj.Write(byte); err == nil {
			fmt.Println("合并成功：", i)
		}
	}

}

//读取本地ts文件
func (t *Tool) ReadTsFile(i int) ([]byte, error) {
	filePth := "./download/ts/" + fmt.Sprintf("%d", i) + ".ts"
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	defer func() {
		f.Close()
		os.Remove(filePth)
	}()

	return ioutil.ReadAll(f)
}
