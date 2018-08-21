package main

import (
	"time"
	"os"
	"fmt"
	"bufio"
	"io"
)

type ReadFromFile struct{
	path string //日志文件的路径
}

func (r *ReadFromFile) Read (rc chan []byte){
	f, err := os.Open(r.path)
	if err != nil{
		panic(fmt.Sprintf("Open file error: %s", err.Error()))
	}

	f.Seek(0,2) // 字符指针移动到文件末尾，开始逐行读取
	rd := bufio.NewReader(f)

	for{
		line, err := rd.ReadBytes('\n')
		if err == io.EOF{
			time.Sleep(500 * time.Millisecond)
			continue
		}else if err != nil{
			panic(fmt.Sprintf("Read bytes error: %s", err.Error()))
		}

		fmt.Println(string(line[:len(line)-1]))
		fmt.Println(string(line[1:len(line)]))

		rc <- line[:len(line)-1]
	}
}