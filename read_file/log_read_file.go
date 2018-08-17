package main

import (
	"io"
	"bufio"
	"os"
	"time"
	"strings"
	"fmt"
)

type Reader interface{
	Read(rc chan []byte)
}

type Writer interface{
	Write(wc chan string)
}

type LogProcess struct{
	r Reader
	w Writer
	rc chan []byte
	wc chan string
}

type ReadFromFile struct{
	path string // 读取文件的路径
}

func (r *ReadFromFile) Read(rc chan []byte){

	// 打开文件
	f,err := os.Open(r.path)
	if err != nil{
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}

	// 从文件末尾开始逐行读取内容
	f.Seek(0,2)
	rd := bufio.NewReader(f)

	for{
		line, err := rd.ReadBytes('\n')
		if err == io.EOF{
			fmt.Println("EOF")
			time.Sleep(500 * time.Millisecond)
			continue
		} 
		
		if err != nil{
			panic(fmt.Sprintf("read error:%s", err.Error()))
		}

		rc <- line[:len(line)-1]
	}
}

type WriteToInfluxDB struct{
	influxDBDataSource string // influxDB data source
}

func (w *WriteToInfluxDB) Write(wc chan string){
	for v := range wc{
		fmt.Println(v)
	}
}

func (l *LogProcess) Parse(){
	for v := range l.rc{
		l.wc <- strings.ToUpper(string(v))
	}
}

func main(){
	reader := &ReadFromFile{
		path : "D:/Temp/ClientError.txt",
	}

	writer := &WriteToInfluxDB{
		influxDBDataSource : "",
	}

	l := &LogProcess{
		r : reader,
		w : writer, 
		rc : make(chan []byte), 
		wc : make(chan string),
   }

   go l.r.Read(l.rc)
   go l.Parse()
   go l.w.Write(l.wc)

	time.Sleep(30 * time.Second)

	// for{}
}