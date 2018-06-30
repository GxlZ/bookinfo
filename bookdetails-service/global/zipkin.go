package global

import (
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"context"
	"github.com/openzipkin/zipkin-go/reporter"
	"time"
	zipkingomodel "github.com/openzipkin/zipkin-go/model"
)

const (
	ZIPKIN_OPTION_TAG      = 1 + iota
	ZIPKIN_OPTION_ANNOTATE
)

type zipkinOption struct {
	OptionType int
	zipkinTag  ZipkinTag
	Annotate   zipkingomodel.Annotation
}
type ZipkinTag struct {
	K string
	V string
}

func NewZipkinReporter() reporter.Reporter {
	return zipkinhttp.NewReporter(
		Conf.Zipkin.Url+"/api/v2/spans",
		zipkinhttp.Timeout(time.Duration(Conf.Zipkin.Reporter.Timeout)*time.Second),
		zipkinhttp.BatchSize(Conf.Zipkin.Reporter.BatchSize),
		zipkinhttp.BatchInterval(time.Duration(Conf.Zipkin.Reporter.BatchInterval)*time.Second),
		zipkinhttp.MaxBacklog(Conf.Zipkin.Reporter.MaxBacklog),
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


func zipkin(svcName string ,ctx context.Context, opts ...zipkinOption) (zipkingo.Span, context.Context, error) {
	span, newCtx, err := NewZipkinSpanFromCtx(ctx, func() (*zipkingo.Tracer, error) {
		return NewZipkinTracer()
	})

	if err != nil {
		return span, newCtx, err
	}

	span.Annotate(time.Now(), "in "+svcName)

	for _, option := range opts {
		switch option.OptionType {
		case ZIPKIN_OPTION_TAG:
			span.Tag(option.zipkinTag.K, option.zipkinTag.V)
		case ZIPKIN_OPTION_ANNOTATE:
			span.Annotate(option.Annotate.Timestamp, option.Annotate.Value)
		}
	}

	span.SetName(svcName+" execute")

	return span, newCtx, err
}
