package main

import (
	"demo/prometheus-demo/biz"
	"demo/prometheus-demo/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Qps)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 注册 ContainerDiskCollector
	// 这里注册的容器名称请调用方自行确保存在
	diskCollector := biz.NewContainerDiskCollector("prometheus")
	prometheus.MustRegister(diskCollector)
	prometheus.MustRegister(collectors.NewGoCollector())

	// 对外提供 /metrics 接口，支持 prometheus 采集
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Fatal(r.Run())
}
