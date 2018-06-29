package handlers

import (
	"context"
	"bookinfo/bookdetails-service/global"
	"bookinfo/bookdetails-service/models"
	pb "bookinfo/pb/details"
	"time"
	zipkingo "github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	commentspb "bookinfo/pb/comments"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.BookDetailsServer {
	return bookdetailsService{}
}

type bookdetailsService struct{}

// Detail implements Service.
func (s bookdetailsService) Detail(ctx context.Context, in *pb.DetailReq) (*pb.DetailResp, error) {
	//zipkin
	var newCtx context.Context
	var zipkinSpan zipkingo.Span
	{
		span, ctx, err := global.NewZipkinSpanFromCtx(ctx, func() (*zipkingo.Tracer, error) {
			return global.NewZipkinTracer()
		})

		zipkinSpan = span
		newCtx = ctx

		if err != nil {
			global.Logger.Error("zipkin span create failed,", err)
		} else {
			span.SetName("get book info")
			span.Annotate(time.Now(), "in svc")

			defer func() {
				span.Annotate(time.Now(), "out svc")
				span.Finish()
			}()
		}
	}

	var resp pb.DetailResp

	var redisKey = fmt.Sprintf("book_detail_%d", in.Id)

	var book models.Books

	//read from cache
	{
		cacheBytes := global.Redis.Get(redisKey).Val()

		if len(cacheBytes) > 0 {
			if err := ffjson.Unmarshal([]byte(cacheBytes), &book); err != nil {
				global.Logger.Warnln("redis get error:", err)
			} else {
				resp.Code = global.SUCCESS.Code
				resp.Msg = global.SUCCESS.Msg
				resp.Data = &pb.DetailRespData{
					Id:    int32(book.ID),
					Name:  book.Name,
					Intro: book.Intro,
				}
				return &resp, nil
			}
		}
	}

	//base info from db
	{

		if in.Id == 0 {
			resp.Code = global.ERROR_PARAMS_ERROR.Code
			resp.Msg = global.ERROR_PARAMS_ERROR.Msg
			return &resp, nil
		}

		global.BOOK_DB.WarpRawScan(newCtx, &book, "select * from books where id = ?", in.Id)

		if book.ID == 0 {
			resp.Code = global.ERROR_RESOURCE_NOT_FOUND.Code
			resp.Msg = global.ERROR_RESOURCE_NOT_FOUND.Msg
			return &resp, nil
		}

		resp.Code = global.SUCCESS.Code
		resp.Msg = global.SUCCESS.Msg
		resp.Data = &pb.DetailRespData{
			Id:    int32(book.ID),
			Name:  book.Name,
			Intro: book.Intro,
		}
	}

	//comments from grpc
	{
		c, _ := global.NewGrpcClient(
			newCtx,
			zipkinSpan,
			global.Conf.Servers.BookComments.Grpc,
			func(ctx context.Context, conn *grpc.ClientConn) (resp interface{}, err error) {
				c := commentspb.NewBookCommentsClient(conn)

				resp, err = c.Get(ctx, &commentspb.GetReq{Id: 1})

				return
			},
			grpc.WithInsecure(),
			grpc.WithTimeout(10*time.Second),
		)
		res, err := c.Go()

		if err != nil {
			return &resp, nil
		}

		commentsResp := res.(*commentspb.GetResp)
		if commentsResp.Code != global.SUCCESS.Code {
			return &resp, nil
		}
		resp.Data.Comments = commentsResp.Data
	}

	go func(book models.Books) {
		bytes, err := ffjson.Marshal(book)
		if err != nil {
			global.Logger.Warnln("json marshal error:", err)
		}
		if err := global.Redis.Set(redisKey, bytes, 3600*time.Second).Err(); err != nil {
			global.Logger.Warnln("redis set error:", err)
		}
	}(book)

	return &resp, nil
}

//定义中间件接口
type SvcMiddleware func(pb.BookDetailsServer) pb.BookDetailsServer
