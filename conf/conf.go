package conf

const (
	Interval  = 5
	Expirtime = 3600 * 24
)

var (
	InfluxAddr  = "http://127.0.0.1:8086"
	InfluxDB    = "ngx"
	InfluxTable = "ngx_metric"
	NgxAddr     = "http://127.0.0.1/status"
	NgxDomain   = "api.test.com"
)
