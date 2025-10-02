package metrics

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"CalculatorWebService/internal/config"
)

type Metrics struct {
	reg        *prometheus.Registry
	Handler    *http.Handler
	Counters   map[string]*prometheus.CounterVec
	baseLabels prometheus.Labels
	mu         sync.RWMutex
}

func NewMetrics(initConfig config.MetricsConfig) *Metrics {
	reg := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	m := &Metrics{
		reg:        reg,
		Handler:    &handler,
		Counters:   make(map[string]*prometheus.CounterVec),
		baseLabels: make(prometheus.Labels),
	}
	m.SetupBaseLabels(initConfig)
	return m
}
func (m *Metrics) SetupBaseLabels(config config.MetricsConfig) {
	if config.ServiceName != "" {
		m.baseLabels["service"] = config.ServiceName
	}
	if config.ServiceVersion != "" {
		m.baseLabels["service_version"] = config.ServiceVersion
	}
	if config.ServerName != "" {
		m.baseLabels["server"] = config.ServerName
	}
}

// CountInc gives flex ability to have multiple metrics on the go wherever needed
// with different labels as needed. Base labels are added automatically.
// Metric is created on the fly if it does not exist yet.
// TODO: more use cases?
func (m *Metrics) CountInc(metricName string, labels prometheus.Labels) {
	metric := m.getCounter(metricName, labels)
	if metric == nil {
		return
	}
	for baseLabel, baseValue := range m.baseLabels {
		labels[baseLabel] = baseValue
	}
	metric.With(labels).Inc()
}
func (m *Metrics) CountAdd(metricName string, labels prometheus.Labels, addValue float64) {
	metric := m.getCounter(metricName, labels)
	if metric == nil {
		return
	}
	for baseLabel, baseValue := range m.baseLabels {
		labels[baseLabel] = baseValue
	}
	metric.With(labels).Add(addValue)
}

func (m *Metrics) newCounter(metricName string, labels prometheus.Labels) *prometheus.CounterVec {
	labelsNames := make([]string, 0, len(labels))
	for baseLabel, _ := range m.baseLabels {
		labelsNames = append(labelsNames, baseLabel)
	}
	for labelName, _ := range labels {
		labelsNames = append(labelsNames, labelName)
	}
	metric := promauto.With(m.reg).NewCounterVec(prometheus.CounterOpts{
		Name: metricName,
		Help: fmt.Sprintf("Counter for %s", metricName),
	}, labelsNames)
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exist := m.Counters[metricName]; exist {
		// someone else created it in the meantime
		return m.Counters[metricName]
	} else {
		m.Counters[metricName] = metric
	}
	return metric
}

func (m *Metrics) getCounter(metricName string, labels prometheus.Labels) *prometheus.CounterVec {
	if metricName == "" {
		return nil
	}
	m.mu.RLock()
	metric, exist := m.Counters[metricName]
	m.mu.RUnlock()
	if !exist {
		metric = m.newCounter(metricName, labels)
	}
	return metric
}

func (m *Metrics) PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		m.CountInc("http_requests_total", prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.FullPath(),
			"status": fmt.Sprintf("%d", status),
		})
	}
}
