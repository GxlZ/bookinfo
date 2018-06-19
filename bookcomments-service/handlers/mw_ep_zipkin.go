package handlers

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"bookinfo/bookcomments-service/global"
	zipkingo "github.com/openzipkin/zipkin-go"
	"time"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"github.com/openzipkin/zipkin-go/model"
	"google.golang.org/grpc/metadata"
)

func ZipkinEndpointMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			zipkinTracer, err := global.NewZipkinTracer()
			if err != nil {
				global.Logger.Error("zipkinTracer create failed,", err)
				return next(ctx, request)
			}

			var sc model.SpanContext
			md, _ := metadata.FromIncomingContext(ctx);
			sc = zipkinTracer.Extract(b3.ExtractGRPC(&md))

			span := zipkinTracer.StartSpan("in service book-comments", zipkingo.Parent(sc))
			span.Annotate(time.Now(), "in endpoint")

			defer func() {
				span.Annotate(time.Now(), "out endpoint")
				span.Finish()
			}()

			ctx = zipkingo.NewContext(ctx, span)
			return next(ctx, request)
		}
	}
}
