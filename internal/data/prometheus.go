package data

import "github.com/prometheus/client_golang/prometheus"

var (
	// SQL语句执行的耗时（直方图）
	_sqlHistogramTracing = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sql_tracing_h_total",
			Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"run_sql"},
	)

	// Redis执行的耗时（直方图）
	_redisHistogramTracing = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_tracing_h_total",
			Buckets: []float64{0.01, 0.03, 0.05, 0.1, 0.3, 0.5, 1},
		},
		[]string{"command_name"},
	)
)

func init() {
	prometheus.MustRegister(
		_sqlHistogramTracing,
		_redisHistogramTracing,
	)
}
