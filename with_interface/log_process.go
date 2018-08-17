package main

import (
	"time"
	"strings"
	"fmt"
)

type Reader interface{
	Read(rc chan string)
}

type Writer interface{
	Write(wc chan string)
}

type LogProcess struct{
	r Reader
	w Writer
	rc chan string
	wc chan string
}

type ReadFromFile struct{
	path string // 读取文件的路径
}

func (r *ReadFromFile) Read(rc chan string){
	Line := "message extend"
	rc <- Line
}

type WriteToInfluxDB struct{
	influxDBDataSource string // influxDB data source
}

func (w *WriteToInfluxDB) Write(wc chan string){
	fmt.Println(<-wc)
}

func (l *LogProcess) Parse(){
	var Line string
	Line = <- l.rc
	l.wc <- strings.ToUpper(Line)
}

func main(){
	reader := &ReadFromFile{
		path : "",
	}

	writer := &WriteToInfluxDB{
		influxDBDataSource : "",
	}

	l := &LogProcess{
		r : reader,
		w : writer, 
		rc : make(chan string), 
		wc : make(chan string),
   }

   go l.r.Read(l.rc)
   go l.Parse()
   go l.w.Write(l.wc)

	time.Sleep(100000)

	// for{}
}