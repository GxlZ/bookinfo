package handlers

import (
	"testing"
	pb "bookinfo/pb/details"
	"bookinfo/pb/comments"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/davecgh/go-spew/spew"
)

type MockBooksDetailsServer struct {
	mock.Mock
}

func (this MockBooksDetailsServer) Detail(ctx context.Context,req *pb.DetailReq) (*pb.DetailResp, error) {
	args := this.Mock.Called(ctx, req)
	spew.Dump(args)
	return args.Get(0).(*pb.DetailResp), args.Error(1)
}

var booksDetailsServer = NewService()

func TestBooksDetailsV1Detail(t *testing.T) {

	mockBooksDetailsServer := MockBooksDetailsServer{}
	mockBooksDetailsServer.On(
		"Detail",
		nil,
		&pb.DetailReq{Id: 1},
	).Return(pb.DetailResp{
		Code: 200,
		Msg:  "success",
		Data: &pb.DetailRespData{
			Id:    1,
			Name:  "西游记",
			Intro: "《西游记》是中国古代第一部浪漫主义章回体长篇神魔小说。现存明刊百回本《西游记》均无作者署名。清代学者吴玉搢等首先提出《西游记》作者是明代吴承恩 [1]  。这部小说以“唐僧取经”这一历史事件为蓝本，通过作者的艺术加工，深刻地描绘了当时的社会现实。全书主要描写了孙悟空出世及大闹天宫后，遇见了唐僧、猪八戒和沙僧三人，西行取经，一路降妖伏魔，经历了九九八十一难，终于到达西天见到如来佛祖，最终五圣成真的故事。",
			Comments: []*comments.Comment{
				{Content: "测试评论1", CreatedAt: "2018-06-20 06:59:40"},
				{Content: "测试评论2", CreatedAt: "2018-06-20 06:59:40"},
				{Content: "测试评论3", CreatedAt: "2018-06-20 06:59:40"},
			},
		},
	}, nil)

	testData := struct {
		in   int64
		want pb.DetailResp
	}{
		in: 1,
		want: pb.DetailResp{
			Code: 200,
			Msg:  "success",
			Data: &pb.DetailRespData{
				Id:    1,
				Name:  "西游记",
				Intro: "《西游记》是中国古代第一部浪漫主义章回体长篇神魔小说。现存明刊百回本《西游记》均无作者署名。清代学者吴玉搢等首先提出《西游记》作者是明代吴承恩 [1]  。这部小说以“唐僧取经”这一历史事件为蓝本，通过作者的艺术加工，深刻地描绘了当时的社会现实。全书主要描写了孙悟空出世及大闹天宫后，遇见了唐僧、猪八戒和沙僧三人，西行取经，一路降妖伏魔，经历了九九八十一难，终于到达西天见到如来佛祖，最终五圣成真的故事。",
				Comments: []*comments.Comment{
					{Content: "测试评论1", CreatedAt: "2018-06-20 06:59:40"},
					{Content: "测试评论2", CreatedAt: "2018-06-20 06:59:40"},
					{Content: "测试评论3", CreatedAt: "2018-06-20 06:59:40"},
				},
			},
		},
	}


	req := &pb.DetailReq{Id: testData.in}
	resp, err := mockBooksDetailsServer.Detail(context.Background(), req)


	if err != nil {
		t.Error("books-details /v1/detail", err)
	}

	if resp.Code != testData.want.Code ||
		resp.Msg != testData.want.Msg ||
		resp.Data.Id != testData.want.Data.Id ||
		resp.Data.Name != testData.want.Data.Name ||
		resp.Data.Intro != testData.want.Data.Intro {
		t.Error(
			"books-details /v1/detail,",
			"req [", req, " ],",
			"resp [", resp, " ],",
			"want resp [", testData.want, " ]",
		)
	}
}

func BenchmarkBooksDetailsV1Detail(b *testing.B) {
	b.ResetTimer()
	req := &pb.DetailReq{Id: 1}
	for i := 0; i < b.N; i++ {
		booksDetailsServer.Detail(context.Background(), req)
	}
}
