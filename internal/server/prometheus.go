package server

import "github.com/prometheus/client_golang/prometheus"

// 定义一个 Histogram 类型的指标
var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duratio(sec).",
		Buckets:   []float64{0.2, 0.5, 1, 2, 5, 10, 30}, // 根据场景需求配置 bucket 的范围
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

func init() {
	prometheus.MustRegister(
		_metricSeconds,
		_metricRequests,
	)
}
