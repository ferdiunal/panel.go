package observability

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

type requestMetrics struct {
	totalRequests  atomic.Uint64
	errorRequests  atomic.Uint64
	totalLatencyNs atomic.Uint64
	statusCounts   sync.Map // map[string]*atomic.Uint64
	routeCounts    sync.Map // map[string]*atomic.Uint64
}

var metricsStore = &requestMetrics{}

func RequestMetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		statusCode := c.Response().StatusCode()
		duration := time.Since(start)

		metricsStore.totalRequests.Add(1)
		metricsStore.totalLatencyNs.Add(uint64(duration.Nanoseconds()))
		if statusCode >= 400 {
			metricsStore.errorRequests.Add(1)
		}

		incrementMetricCounter(&metricsStore.statusCounts, fmt.Sprintf("%d", statusCode))
		routeKey := c.Method() + " " + c.Path()
		incrementMetricCounter(&metricsStore.routeCounts, routeKey)

		return err
	}
}

func MetricsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		total := metricsStore.totalRequests.Load()
		totalLatency := metricsStore.totalLatencyNs.Load()
		avgLatency := 0.0
		if total > 0 {
			avgLatency = float64(totalLatency) / float64(total) / float64(time.Second)
		}

		var b strings.Builder
		b.WriteString("# TYPE panel_requests_total counter\n")
		b.WriteString(fmt.Sprintf("panel_requests_total %d\n", total))
		b.WriteString("# TYPE panel_request_errors_total counter\n")
		b.WriteString(fmt.Sprintf("panel_request_errors_total %d\n", metricsStore.errorRequests.Load()))
		b.WriteString("# TYPE panel_request_duration_seconds_avg gauge\n")
		b.WriteString(fmt.Sprintf("panel_request_duration_seconds_avg %.6f\n", avgLatency))

		statusKeys := collectMetricKeys(&metricsStore.statusCounts)
		sort.Strings(statusKeys)
		for _, key := range statusKeys {
			counter := loadMetricCounter(&metricsStore.statusCounts, key)
			if counter == nil {
				continue
			}
			b.WriteString(fmt.Sprintf("panel_request_status_total{status=\"%s\"} %d\n", key, counter.Load()))
		}

		routeKeys := collectMetricKeys(&metricsStore.routeCounts)
		sort.Strings(routeKeys)
		for _, key := range routeKeys {
			counter := loadMetricCounter(&metricsStore.routeCounts, key)
			if counter == nil {
				continue
			}
			parts := strings.SplitN(key, " ", 2)
			method := parts[0]
			path := ""
			if len(parts) > 1 {
				path = parts[1]
			}
			b.WriteString(fmt.Sprintf("panel_request_route_total{method=\"%s\",path=\"%s\"} %d\n", method, path, counter.Load()))
		}

		c.Set(fiber.HeaderContentType, "text/plain; charset=utf-8")
		return c.SendString(b.String())
	}
}

func incrementMetricCounter(m *sync.Map, key string) {
	counter := loadMetricCounter(m, key)
	if counter == nil {
		newCounter := &atomic.Uint64{}
		actual, _ := m.LoadOrStore(key, newCounter)
		counter = actual.(*atomic.Uint64)
	}
	counter.Add(1)
}

func loadMetricCounter(m *sync.Map, key string) *atomic.Uint64 {
	value, ok := m.Load(key)
	if !ok {
		return nil
	}
	counter, ok := value.(*atomic.Uint64)
	if !ok {
		return nil
	}
	return counter
}

func collectMetricKeys(m *sync.Map) []string {
	keys := make([]string, 0)
	m.Range(func(key, _ interface{}) bool {
		k, ok := key.(string)
		if ok {
			keys = append(keys, k)
		}
		return true
	})
	return keys
}
