package handlers

import (
	"forum/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func (h *Handler) metricsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		metrics.RequestCounter.WithLabelValues(r.Method, r.URL.Path).Inc()

		next.ServeHTTP(w, r)
	}
}
