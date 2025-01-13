package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var counter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "The total number of HTTP request",
	},
	[]string{"uri"},
)

// 注册 Counter 指标采集器
func init() {
	prometheus.MustRegister(counter)
}

// Qps 采集中间件
// 可以根据 label 筛选指标，如下所示
// http_request_total{uri="/api/xxx/detail"}
// 可以使用 irate 函数查询 QPS
// irate(http_request_total{uri="/api/xxx/detail"}[10m]) 十分钟内这个接口的平均 QPS
func Qps(c *gin.Context) {
	counter.WithLabelValues(c.Request.RequestURI).Inc()
}
