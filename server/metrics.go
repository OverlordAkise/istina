package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
    _ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
    "github.com/dlmiddlecote/sqlstats"
)

func RegisterMetrics(r *gin.Engine, db *sqlx.DB) {

	//Counter for only uphill data, Vec for variable labels
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "istina",
			Name:      "http_requests_total",
			Help:      "Number of received http requests",
		},
		[]string{"code", "method", "host", "url"},
	)
	prometheus.MustRegister(requestCounter)

	// Histogram for buckets, Vec for variable labels
	durationCounter := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "istina",
			Name:      "http_duration_seconds",
			Help:      "HTTP request latencies in seconds",
			Buckets:   []float64{0.0001, 0.001, 0.01, 0.1, 0.3, 0.5, 1, 2},
		},
		[]string{"code", "method", "url"},
	)
	prometheus.MustRegister(durationCounter)

	//Summary for variable data that maybe doesnt fit into buckets
	sizeCounter := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "istina",
			Name:      "http_request_size_bytes",
			Help:      "HTTP request size in bytes",
		},
	)
	prometheus.MustRegister(sizeCounter)

	//Middleware
	r.Use(func(c *gin.Context) {
		t := time.Now()
		//bodylength, _ := io.Copy(ioutil.Discard, c.Request.Body)
		bodylength := computeApproximateRequestSize(c.Request)

		c.Next()

		elapsed := float64(time.Since(t)) / float64(time.Second)
		status := strconv.Itoa(c.Writer.Status())

		if c.Request.URL.Path != "/metrics" && c.Request.URL.Path != "/favicon.ico" {
			requestCounter.WithLabelValues(status, c.Request.Method, c.Request.Host, c.Request.URL.Path).Inc()
			durationCounter.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(elapsed)
			sizeCounter.Observe(float64(bodylength))
		}
	})
    
    collector := sqlstats.NewStatsCollector("istina", db)
    prometheus.MustRegister(collector)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// From https://github.com/zsais/go-gin-prometheus/blob/2199a42d96c1d40f249909ed2f27d42449c7fc94/middleware.go#L397
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
