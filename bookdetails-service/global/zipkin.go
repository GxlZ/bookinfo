package global

import (
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"context"
	"github.com/openzipkin/zipkin-go/reporter"
	"time"
)

func NewZipkinReporter() reporter.Reporter {
	return zipkinhttp.NewReporter(
		Conf.Zipkin.Url+"/api/v2/spans",
		zipkinhttp.Timeout(time.Duration(Conf.Zipkin.Reporter.Timeout)*time.Second),
		zipkinhttp.BatchSize(Conf.Zipkin.Reporter.BatchSize),
		zipkinhttp.BatchInterval(time.Duration(Conf.Zipkin.Reporter.BatchInterval)*time.Second),
	)
}

func NewZipkinSpanFromCtx(ctx context.Context, f zipkinTracerFunc) (span zipkingo.Span, newCtx context.Context, err error) {
	zipkinTracer, err := f()
	if err != nil {
		return
	}

	span, newCtx = zipkinTracer.StartSpanFromContext(ctx, "book-details")

	return
}

type zipkinTracerFunc func() (*zipkingo.Tracer, error)

func NewZipkinTracer(opts ...zipkingo.TracerOption) (*zipkingo.Tracer, error) {
	zEP, _ := zipkingo.NewEndpoint(
		Conf.Zipkin.ServiceName,
		"localhost"+Conf.GrpcServer.Addr,
	)

	tracerOptions := []zipkingo.TracerOption{
		zipkingo.WithLocalEndpoint(zEP),
	}

	zipkinTracer, err := zipkingo.NewTracer(
		zipkinReporter,
		append(tracerOptions, opts...)...,
	)

	return zipkinTracer, err
}