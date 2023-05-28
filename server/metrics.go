package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

func RegisterMetrics(r *gin.Engine) {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "istina",
			Name:      "http_requests_total",
			Help:      "Number of received http requests",
		},
		[]string{"code", "method", "handler", "host", "url"},
	)
	if err := prometheus.Register(requestCounter); err != nil {
		panic(err)
	}
	durationCounter := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "istina",
			Name:      "http_duration_seconds",
			Help:      "HTTP request latencies in seconds",
		},
		[]string{"code", "method", "url"},
	)
	if err := prometheus.Register(durationCounter); err != nil {
		panic(err)
	}
	sizeCounter := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "istina",
			Name:      "http_request_size_bytes",
			Help:      "HTTP request size in bytes",
		},
	)
	if err := prometheus.Register(sizeCounter); err != nil {
		panic(err)
	}

	//Middleware
	r.Use(func(c *gin.Context) {
		t := time.Now()
		bodylength, _ := io.Copy(ioutil.Discard, c.Request.Body)

		c.Next()

		elapsed := float64(time.Since(t)) / float64(time.Second)
		status := strconv.Itoa(c.Writer.Status())

		if c.Request.URL.Path != "/metrics" && c.Request.URL.Path != "/favicon.ico" {
			requestCounter.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, c.Request.URL.Path).Inc()
			durationCounter.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(elapsed)
			sizeCounter.Observe(float64(bodylength))
		}
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
