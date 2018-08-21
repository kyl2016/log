package main

import (
	"log"
	"github.com/influxdata/influxdb/client/v2"
	"strings"
)

type WriteToInfluxDB struct{
	influxDBDsn string
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	infSli := strings.Split(w.influxDBDsn, "@")

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: infSli[0],
		Username: infSli[1],
		Password: infSli[2],
	})

	if err != nil{
		log.Fatal(err)
	}

	for v := range wc{
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database: infSli[3],
			Precision: infSli[4],
		})

		if err != nil{
			log.Fatal(err)
		}

		tags := map[string]string{"Path":v.Path, "Method":v.Method, "Scheme":v.Scheme, "Status":v.Status}
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime": v.RequestTime,
			"BytesSent": v.BytesSent,
		}

		pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil{
			log.Fatal(err)
		}

		bp.AddPoint(pt)

		if err := c.Write(bp); err != nil{
			log.Fatal(err)
		}

		log.Println("write success!")
	}
}