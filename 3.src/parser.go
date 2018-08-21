package main

import (
	"strconv"
	"log"
	"net/url"
	"strings"
	"time"
	"regexp"
)

type LogProcess struct{
	rc chan []byte
	wc chan *Message
	read Reader
	write Writer
}

func (l *LogProcess) Process() {
	// 172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	location, _ := time.LoadLocation("Asia/Shanghai")

	for v := range l.rc{
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14{
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		msg := &Message{}

		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], location)
		if err != nil{
			log.Println("ParseInLocation faile:", err.Error(), ret[4])
			continue
		}
		msg.TimeLocal = t

		bytesSent, _ := strconv.Atoi(ret[8])
		msg.BytesSent = bytesSent

		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3{
			log.Println("string.Split fail:", ret[6])
			continue
		}

		u, err := url.Parse(reqSli[1])
		if err != nil{
			log.Println("url Parse fail:", err)
			continue
		}
		msg.Path = u.Path

		msg.Scheme = ret[5]
		msg.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		msg.UpstreamTime = upstreamTime

		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		msg.RequestTime = requestTime

		l.wc <- msg
	}
}