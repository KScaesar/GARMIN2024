package api

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewO11Y() *O11Y {
	o11y := &O11Y{
		ResponseMilliSeconds: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_response_milliseconds",
			Help:    "Histogram of response times for api in milliseconds",
			Buckets: []float64{10, 20, 50, 100, 300, 600, 1_000, 2_000, 5_000, 10_000, 30_000, 60_000}, // 10 ms ~ 60 s
		}, []string{"version", "method", "path", "status"}),

		RequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}, []string{"version", "method", "path", "status"}),
	}

	return o11y
}

type O11Y struct {
	ResponseMilliSeconds *prometheus.HistogramVec
	RequestsTotal        *prometheus.CounterVec
}

func (o11y *O11Y) Middleware(version string) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := CtxWithO11Y(c.Request.Context(), o11y)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()

		c.Next()

		method := c.Request.Method
		path := c.Request.URL.Path

		duration := time.Since(start).Seconds() * 1000 // convert ms
		status := strconv.Itoa(c.Writer.Status())

		o11y.ResponseMilliSeconds.WithLabelValues(version, method, path, status).Observe(duration)
		o11y.RequestsTotal.WithLabelValues(version, method, path, status).Inc()
	}
}

var o11yKey = "o11y"

func CtxWithO11Y(ctx context.Context, v *O11Y) context.Context {
	return context.WithValue(ctx, &o11yKey, v)
}

func CtxGetO11Y(ctx context.Context) (o11y *O11Y) {
	o11y, ok := ctx.Value(&o11yKey).(*O11Y)
	if !ok {
		return NewO11Y()
	}
	return o11y
}
