package biz

import (
	"context"
	"demo/prometheus-demo/data"
	"demo/prometheus-demo/domain"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
)

// ContainerDiskCollector 定义自定义 Collector
type ContainerDiskCollector struct {
	container string
	desc      *prometheus.Desc
}

// NewContainerDiskCollector 创建一个新的 ContainerDiskCollector
// 参数：容器名或容器id，请调用方自行确保存在
func NewContainerDiskCollector(container string) *ContainerDiskCollector {
	return &ContainerDiskCollector{
		container: container,
		desc: prometheus.NewDesc(
			"container_disk_usage",
			"Dynamic Container Disk Usage",
			[]string{"type", "mount"},
			nil,
		),
	}
}

// Describe 实现 prometheus.Collector 接口
func (c *ContainerDiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

// Collect 动态采集数据
func (c *ContainerDiskCollector) Collect(ch chan<- prometheus.Metric) {
	// 动态获取 ContainerDisk 数据
	container := data.NewContainer(c.container)
	usage, err := container.GetDiskUsage(context.Background())
	if err != nil {
		log.Fatal("获取容器磁盘使用状况失败", err)
		return
	}

	monitor := domain.NewContainerMonitor()
	disks, err := monitor.ParseOutput(usage)
	if err != nil {
		log.Fatal("解析失败:", err)
		return
	}

	// 将每个磁盘的数据发送到 Prometheus
	for _, disk := range disks {
		usedRate := strings.TrimSuffix(disk.UsedRate, "%")
		ch <- prometheus.MustNewConstMetric(
			c.desc,
			prometheus.GaugeValue,
			cast.ToFloat64(usedRate),
			disk.Type,
			disk.Mount,
		)
	}
}
