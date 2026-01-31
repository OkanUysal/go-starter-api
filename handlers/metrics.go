package handlers

import "github.com/OkanUysal/go-metrics"

var Metrics *metrics.Metrics

func SetMetrics(m *metrics.Metrics) {
	Metrics = m
}
