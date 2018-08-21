package main

import (
	"flag"
)

func main(){
	var path, influxDBDsn string
	flag.StringVar(&path, "path", "D:\\07.Demo\\Go\\src\\github.com\\kyl2016\\log\\src\\access.log", "read file path")	
	flag.StringVar(&influxDBDsn, "influxDBDsn", "http://localhost:8086@test@test@mydb@s", "influx data source")
	flag.Parse()

	r := &ReadFromFile{
		path: path,
	}

	w := &WriteToInfluxDB{
		influxDBDsn: influxDBDsn,
	}

	lp := &LogProcess{
		rc: make(chan []byte, 200),
		wc: make(chan *Message, 200),
		read: r, 
		write: w,
	}

	go lp.read.Read(lp.rc)

	for i:=0; i<2; i++{
		go lp.Process()
	}

	for i:=0; i<4; i++{
		go lp.write.Write(lp.wc)
	}

	for{

	}
}