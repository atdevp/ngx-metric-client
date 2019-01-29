package utils

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ngx_metric/conf"
)

func Influx() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: conf.InfluxAddr,
	})
	if err != nil {
		FileLogs.Critical("Error creating InfluxDB Client: %s", err.Error())
	}
	defer cli.Close()
	return cli
}
