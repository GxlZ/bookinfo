package handlers

import (
	"github.com/go-kit/kit/metrics"
	pb "bookinfo/pb/details"
	"time"
	"context"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func Instrumenting() SvcMiddleware {
	return func(next pb.BookDetailsServer) pb.BookDetailsServer {
		fieldKeys := []string{"method", "error"}
		requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "books_business",
			Subsystem: "book_details",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
		requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "books_business",
			Subsystem: "book_details",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys)
		return Instrmw{
			requestCount,
			requestLatency,
			next,
		}
	}
}

type Instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	pb.BookDetailsServer
}

func (this Instrmw) Detail(ctx context.Context, req *pb.DetailReq) (*pb.DetailResp, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "/v1/detail", "error", "false"}
		this.requestCount.With(lvs...).Add(1)
		this.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return this.BookDetailsServer.Detail(ctx, req)
}
