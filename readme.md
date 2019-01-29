ngx-metric-client
===

Regularly request the ngx status interface to get data stored in influxdb.


## Configuration

- Interval: Request ngx status interval time. 
- Expirtime: The maximum time the local cache retains the key.
- InfluxAddr: Influxdb address. For example: http://127.0.0.1:8086
- InfluxDB: Database.
- NgxAddr: Ngx status. For example: http://127.0.0.1/status
- NgxDomain: Ngx domain. For example: api.test.com


## Installation

It is a golang project

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/
cd $GOPATH/src/github.com/
git clone https://github.com/atdevp/ngx-metric-client.git
cd ngx-metric-client
go get
./control build
./control start

# goto influxdb to select data
```
