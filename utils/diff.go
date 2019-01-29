package utils

import (
	"fmt"
	"time"

	"github.com/ngx_metric/conf"
	cache "github.com/patrickmn/go-cache"
)

func Diff(inter string, metric string, newv float64, c *cache.Cache) (diff float64, ok bool) {
	longkey := fmt.Sprintf("%s:%s", inter, metric)
	oldv, ok := c.Get(longkey)

	if ok {
		num, ok := oldv.(float64)
		if !ok {
			FileLogs.Error("%s assertion float64 error", oldv)
			return 0, false
		}
		c.Set(longkey, newv, conf.Expirtime*time.Second)

		diff := Decimal(newv - num)
		if diff < 0 {
			return 0, false
		}

		return diff, true
	} else {
		c.Set(longkey, newv, conf.Expirtime*time.Second)
		FileLogs.Error("no exist %s", longkey)
		return Decimal(newv), true
	}
}
