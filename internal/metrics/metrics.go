package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Registry            *prometheus.Registry
	TasksStartedTotal   prometheus.Counter
	TasksCompletedTotal prometheus.Counter
	TasksFailedTotal    prometheus.Counter
	TasksInProgress     prometheus.Gauge
	HTTPRequestsTotal   *prometheus.CounterVec
}

var (
	once     sync.Once
	instance *Metrics
)

func NewMetrics() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			Registry: prometheus.NewRegistry(),
			TasksStartedTotal: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "tasks_started_total",
				},
			),
			TasksCompletedTotal: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "tasks_completed_total",
				},
			),
			TasksFailedTotal: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "tasks_failed_total",
				},
			),
			TasksInProgress: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: "tasks_in_progress",
				},
			),
			HTTPRequestsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
				},
				[]string{"endpoint", "status"},
			),
		}
		instance.Registry.MustRegister(
			instance.TasksStartedTotal,
			instance.TasksCompletedTotal,
			instance.TasksFailedTotal,
			instance.TasksInProgress,
			instance.HTTPRequestsTotal,
		)
	})
	return instance
}
