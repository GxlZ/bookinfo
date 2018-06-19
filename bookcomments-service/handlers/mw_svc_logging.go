package handlers

import (
	"time"
	"context"
	"bookinfo/bookcomments-service/global"
	pb "bookinfo/pb/comments"
)

func LoggingSvcMiddleware() SvcMiddleware {
	return func(next pb.BookCommentsServer) pb.BookCommentsServer {
		return loggingSvcMiddleware{next}
	}
}

type loggingSvcMiddleware struct {
	Next pb.BookCommentsServer
}

func (this loggingSvcMiddleware) Post(ctx context.Context, in *pb.PostReq) (resp *pb.PostResp, err error) {
	defer func(begin time.Time) {
		global.Logger.InfoWithFields(func() *global.LogFields {
			return global.NewLogFields().
				Append("method", "post").
				Append("input", in).
				Append("err", err).
				Append("duration", time.Since(begin))
		}, "out svc")
	}(time.Now())

	return this.Next.Post(ctx, in)
}

func (this loggingSvcMiddleware) Get(ctx context.Context, in *pb.GetReq) (resp *pb.GetResp, err error) {
	defer func(begin time.Time) {
		global.Logger.InfoWithFields(func() *global.LogFields {
			return global.NewLogFields().
				Append("method", "get").
				Append("input", in).
				Append("err", err).
				Append("duration", time.Since(begin))
		}, "out svc")
	}(time.Now())

	return this.Next.Get(ctx, in)
}
