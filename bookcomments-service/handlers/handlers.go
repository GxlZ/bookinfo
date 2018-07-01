package handlers

import (
	"context"

	pb "bookinfo/pb/comments"
	"bookinfo/bookcomments-service/global"
	"github.com/pquerna/ffjson/ffjson"
	"fmt"
	"bookinfo/bookcomments-service/models"
	"time"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.BookCommentsServer {
	return bookcommentsService{}
}

type bookcommentsService struct{}

// Post implements Service.
func (s bookcommentsService) Post(ctx context.Context, in *pb.PostReq) (*pb.PostResp, error) {
	var resp pb.PostResp
	resp = pb.PostResp{
		// Code:
		// Msg:
	}
	return &resp, nil
}

// Get implements Service.
func (s bookcommentsService) Get(ctx context.Context, in *pb.GetReq) (*pb.GetResp, error) {
	var m map[int]int
	m[1]=1

	var resp pb.GetResp

	if in.Id == 0 {
		resp.Code = global.ERROR_PARAMS_ERROR.Code
		resp.Msg = global.ERROR_PARAMS_ERROR.Msg
		return &resp, nil
	}

	var redisKey = fmt.Sprintf("comments_%d", in.Id)

	// get from cache
	{
		cacheBytes := global.Redis.WarpGet(ctx, redisKey).Val()

		var comments []*pb.Comment

		if len(cacheBytes) > 0 {
			if err := ffjson.Unmarshal([]byte(cacheBytes), &comments); err != nil {
				global.Logger.Warnln("redis get error:", err)
			} else {
				resp.Code = global.SUCCESS.Code
				resp.Msg = global.SUCCESS.Msg
				resp.Data = comments
				return &resp, nil
			}
		}
	}

	// get from db
	var comments []models.Comments
	global.BOOK_DB.WarpRawScan(ctx, &comments, "select * from comments where book_id = ?", in.Id)

	data := []*pb.Comment{}
	for _, comment := range comments {
		data = append(data, &pb.Comment{
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp.Code = global.SUCCESS.Code
	resp.Msg = global.SUCCESS.Msg
	resp.Data = data

	go func(ctx context.Context, comments []models.Comments) {
		if err := global.Redis.WarpSet(ctx, redisKey, comments, 3600*time.Second).Err(); err != nil {
			global.Logger.Warnln("redis set error:", err)
		}
	}(ctx, comments)
	return &resp, nil
}

type SvcMiddleware func(pb.BookCommentsServer) pb.BookCommentsServer
