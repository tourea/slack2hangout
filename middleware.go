package main

import (
"context"
"strconv"
"time"

"github.com/go-kit/kit/log"
"github.com/go-kit/kit/log/level"
"github.com/go-kit/kit/metrics"
kitprometheus "github.com/go-kit/kit/metrics/prometheus"
stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// Middleware wrapper for instrumentation and logging
type Middleware func(Service) Service

type loggingMiddleware struct {
    next           Service
    logger         log.Logger
    requestCount   metrics.Counter
    requestLatency metrics.Histogram
}

// LoggingMiddleware instrumentation and logging
func LoggingMiddleware(logger log.Logger) Middleware {
    fieldKeys := []string{"method", "error"}
    var requestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
        Namespace: "iris",
        Subsystem: "slack2hangout",
        Name:      "request_total",
        Help:      "Number of requests received.",
    }, fieldKeys)

    requestLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
        Namespace: "iris",
        Subsystem: "slack2hangout",
        Name:      "request_latency_seconds",
        Help:      "Total duration of requests in seconds.",
    }, fieldKeys)

    return func(next Service) Service {
        return &loggingMiddleware{
            next,
            logger,
            requestCount,
            requestLatency,
        }
    }
}

func (mw loggingMiddleware) Forward(ctx context.Context, space string, key string, token string, slackNotification SlackNotification) (err error) {
    defer func(begin time.Time) {
        lvs := []string{"method", "Forward", "error", strconv.FormatBool(err != nil)}
        mw.requestCount.With(lvs...).Add(1)
        mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())

        if err == nil {
            _ = level.Debug(mw.logger).Log("method", "Forward", "space", space, "took", time.Since(begin))
        } else {
            _ = level.Debug(mw.logger).Log("method", "Forward", "space", space, "took", time.Since(begin), "err", err)
        }
    }(time.Now())

    return mw.next.Forward(ctx, space, key, token, slackNotification)
}
