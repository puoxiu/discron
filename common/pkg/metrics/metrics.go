package metrics

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var (
	// CPU核心数及整体使用率
	cpuCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discron_node_cpu_cores_total",
			Help: "Total number of CPU cores",
		},
		[]string{"host"},
	)
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discron_node_cpu_usage_percent",
			Help: "Overall CPU usage percentage",
		},
		[]string{"host"},
	)

	// 内存核心指标
	memTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discron_node_mem_total_bytes",
			Help: "Total memory bytes",
		},
		[]string{"host"},
	)
	memUsedPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discron_node_mem_used_percent",
			Help: "Used memory percentage",
		},
		[]string{"host"},
	)

	// 磁盘使用率（根分区）
	diskUsedPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discron_node_disk_used_percent",
			Help: "Disk usage percentage of root partition",
		},
		[]string{"host"},
	)

	// 应用进程指标
	goRoutines = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discron_app_goroutines_total",
			Help: "Total number of goroutines",
		},
		[]string{"host"},
	)
)

func init() {
	// 仅注册核心指标
	prometheus.MustRegister(
		cpuCount, cpuUsage,
		memTotal, memUsedPercent,
		diskUsedPercent,
		goRoutines,
	)
}

// StartMetricsServer 启动简化版指标服务
func StartMetricsServer(port string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("[metrics] core metrics exposed at :%s/metrics", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	// 5秒采集一次核心指标（平衡实时性和性能）
	go func() {
		updateCoreMetrics() // 立即采集一次
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateCoreMetrics()
		}
	}()
}

// 采集核心指标
func updateCoreMetrics() {
	hostname, _ := os.Hostname()
	hostLabel := prometheus.Labels{"host": hostname}

	// CPU指标
	if count, err := cpu.Counts(true); err == nil {
		cpuCount.With(hostLabel).Set(float64(count))
	}
	if usage, err := cpu.Percent(1*time.Second, false); err == nil && len(usage) > 0 {
		cpuUsage.With(hostLabel).Set(usage[0])
	}

	// 内存指标
	if memInfo, err := mem.VirtualMemory(); err == nil {
		memTotal.With(hostLabel).Set(float64(memInfo.Total))
		memUsedPercent.With(hostLabel).Set(memInfo.UsedPercent)
	}

	// 磁盘指标（仅监控根分区）
	if diskInfo, err := disk.Usage("/"); err == nil {
		diskUsedPercent.With(hostLabel).Set(diskInfo.UsedPercent)
	}

	// 应用goroutine数量
	goRoutines.With(hostLabel).Set(float64(runtime.NumGoroutine()))
}