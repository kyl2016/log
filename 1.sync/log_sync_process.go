package main

import (
	"time"
	"strings"
	"fmt"
)

type LogProcess struct{
	path 				string // 读取文件的路径
	influxDBDataSource 	string // influxDB data source
	rc 					chan string // 读取数据存放的channel
	wc 					chan string // 从解析模块到存储模块传递数据
}

// 读取模块
func (l *LogProcess) ReadFromFile() {
	Line := "message"
	l.rc <- Line
}

// 解析模块
func (l *LogProcess) Parse(){
	var Line string
	Line = <- l.rc

	l.wc <- strings.ToUpper(Line)
}

// 写入模块
func (l *LogProcess) WriteToInfluxDB(){
	fmt.Println(<- l.wc)
}

func main(){
	lp := &LogProcess{
		rc: make(chan string),
		wc: make(chan string),
		path: "/tmp/access.log",
		influxDBDataSource: "username&password...",
	}

	go lp.ReadFromFile() // 等价于: go (* lp).ReadFromFile()
	go lp.Parse()
	go lp.WriteToInfluxDB()

	time.Sleep(10 * time.Second)
}