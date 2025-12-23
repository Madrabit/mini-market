package web

import (
	"github.com/go-chi/chi/v5"
	mv "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var registerOnce sync.Once

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The total number of HTTP requests made.",
		},
		[]string{"method", "route", "status"},
	)

	httpInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Current number of in-flight HTTP requests.",
		},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "route", "status"},
	)

	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes.",
			Buckets: []float64{200, 500, 1_000, 5_000, 10_000, 50_000, 100_000, 500_000, 1_000_000, 5_000_000},
		},
		[]string{"method", "route", "status"},
	)

	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes (from Content-Length when available).",
			Buckets: []float64{200, 500, 1_000, 5_000, 10_000, 50_000, 100_000, 500_000, 1_000_000, 5_000_000},
		},
		[]string{"method", "route"},
	)
)

func Init(reg prometheus.Registerer) {
	registerOnce.Do(func() {
		reg.MustRegister(
			httpRequestTotal,
			httpInFlight,
			httpRequestDuration,
			httpResponseSize,
			httpRequestSize)
	})
}

func Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		ww := mv.NewWrapResponseWriter(w, r.ProtoMajor)
		httpInFlight.Inc()
		defer httpInFlight.Dec()
		next.ServeHTTP(ww, r)
		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			route = "unknown"
		}
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		httpRequestTotal.WithLabelValues(r.Method, route, strconv.Itoa(ww.Status())).Inc()
		httpRequestDuration.WithLabelValues(r.Method, route, strconv.Itoa(ww.Status())).Observe(time.Since(start).Seconds())
		httpResponseSize.WithLabelValues(r.Method, route, strconv.Itoa(ww.Status())).Observe(float64(ww.BytesWritten()))
		// Request size: только если Content-Length известен
		if r.ContentLength >= 0 {
			httpRequestSize.WithLabelValues(r.Method, route).Observe(float64(r.ContentLength))
		}
	})
}
