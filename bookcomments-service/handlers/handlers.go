package handlers

import (
	"context"

	pb "bookinfo/pb/comments"
	"bookinfo/bookcomments-service/global"
	"bookinfo/bookcomments-service/models"
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
	var resp pb.GetResp

	if in.Id == 0 {
		resp.Code = global.ERROR_PARAMS_ERROR.Code
		resp.Msg = global.ERROR_PARAMS_ERROR.Msg
		return &resp, nil
	}

	comments := []models.Comments{}

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

	return &resp, nil
}

type SvcMiddleware func(pb.BookCommentsServer) pb.BookCommentsServer
