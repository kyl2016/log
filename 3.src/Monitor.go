package main

import (
	"time"
	"encoding/json"
	"net/http"
	"io"
)

type SystemInfo struct{
	HandleLine 		int		`json:"handle_line"`		// 总处理日志行数
	Tps				float64	`json:"tps"`			// 系统吞吐量
	ReadChanLen		int		`json:"readChanLen"`	// read channel 长度
	WriteChanLen	int		`json:"writeChanLen"`	// write channel 长度
	RunTime			string	`json:"runTime"`		// 运行总时间
	ErrNum			int		`json:"errNum"`			// 错误数
}

type Monitor struct{
	startTime 	time.Time
	data		SystemInfo
	tpsSli		[]int
}

const(
	TypeHandleLine 	= 0
	TypeErrNum		= 1
)

var TypeMonitorChan = make(chan int, 200)

func (m *Monitor) start(lp *LogProcess){
	go func(){
		for n := range TypeMonitorChan {
			switch n {
				case TypeErrNum: m.data.ErrNum += 1
				case TypeHandleLine: m.data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Second)

	go func(){
		for{
			<- ticker.C
			m.tpsSli = append(m.tpsSli, m.data.HandleLine)
			if len(m.tpsSli) > 2{
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request){
		m.data.RunTime = time.Now().Sub(m.startTime).String()
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriteChanLen = len(lp.wc)

		if len(m.tpsSli) > 2 {
			m.data.Tps = float64(m.tpsSli[1] - m.tpsSli[0]) / 5
		}

		ret, _ := json.MarshalIndent(m.data, "", "\t")
		io.WriteString(w, string(ret))
	})

	http.ListenAndServe(":9193", nil)
}