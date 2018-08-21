package main

import (
	"time"
)

type Reader interface{
	Read(rc chan []byte)
}

type Writer interface{
	Write(wc chan *Message)
}

type Message struct {
	TimeLocal						time.Time			// 时间戳
	BytesSent						int					// 流量
	Path, Method, Scheme, Status 	string 				// 路径，方法，协议，状态
	UpstreamTime, RequestTime 		float64 			// 延迟，响应时间
}