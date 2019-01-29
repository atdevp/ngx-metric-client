package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/httplib"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ngx_metric/conf"
	"github.com/ngx_metric/utils"
	cache "github.com/patrickmn/go-cache"
)

type BaseMetirc struct {
	RequestTime float64 `json:"req_time"`
	CodeOk      float64 `json:"code_ok"`
	Code5xx     float64 `json:"code_5xx"`
	Code4xx     float64 `json:"code_4xx"`
}

var c *cache.Cache
var cli client.Client

func init() {
	c = cache.New(3600*24*time.Second, 3600*time.Second)
	cli = utils.Influx()
}

func getMetrics(uri string) (item map[string]*BaseMetirc) {
	r := httplib.Get(uri)
	r.SetTimeout(1*time.Second, 3*time.Second)
	err := r.ToJSON(&item)
	if err != nil {
		utils.FileLogs.Error("get metrics failed. errmsg: ", err.Error())
	}
	return
}

func InsertDB(uri string) {
	metrics := getMetrics(uri)
	for k, v := range metrics {
		go insertDB(k, v)
	}
}

func getFields(key string, value *BaseMetirc) (fields map[string]interface{}, ok bool) {
	diffReqTime, ok := utils.Diff(key, "req_time", value.RequestTime, c)
	if !ok {
		return fields, false
	}
	diffCodeOk, ok := utils.Diff(key, "code_ok", value.CodeOk, c)
	if !ok {
		return fields, false
	}
	diffCode4xx, _ := utils.Diff(key, "code_4xx", value.Code4xx, c)
	diffCode5xx, _ := utils.Diff(key, "code_5xx", value.Code5xx, c)

	reqtime := diffReqTime / diffCodeOk * 1000
	qps := diffCodeOk / conf.Interval

	fields = map[string]interface{}{
		"req_time": utils.Decimal(reqtime),
		"qps":      qps,
		"code_2xx": diffCodeOk,
		"code_4xx": diffCode4xx,
		"code_5xx": diffCode5xx,
	}
	return fields, true
}

func insertDB(key string, value *BaseMetirc) {

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conf.InfluxDB,
		Precision: "s",
	})

	tags := map[string]string{
		"inter":    key,
		"domail":   conf.NgxDomain,
		"ngx_addr": conf.NgxAddr,
	}

	fields, ok := getFields(key, value)
	if !ok {
		return
	}

	pt, err := client.NewPoint(conf.InfluxTable, tags, fields, time.Now())
	if err != nil {
		errmsg := fmt.Sprintf("create newpoint failed, errmsg: %s, key: %s", err.Error(), tags["inter"])
		utils.FileLogs.Error(errmsg)
	}
	bp.AddPoint(pt)

	if err = cli.Write(bp); err != nil {
		utils.FileLogs.Error("write infludb failed, errmsg: %s", err.Error())
	}
}

func main() {
	duration := time.Duration(conf.Interval) * time.Second
	for {
		go InsertDB(conf.NgxAddr)
		time.Sleep(duration)
	}
	select {}
}
