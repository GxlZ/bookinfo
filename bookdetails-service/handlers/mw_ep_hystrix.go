package handlers

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/afex/hystrix-go/hystrix"
	"context"
	pb "bookinfo/pb/details"
	"bookinfo/bookdetails-service/global"
)

func Hystrix(commandName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			err = hystrix.Do(commandName, func() (err error) {
				resp, err = next(ctx, request)
				return err
			}, nil)

			if err != nil {
				return pb.DetailResp{
					Code: global.ERROR_TOO_MANY_CONNECTIONS.Code,
					Msg:  global.ERROR_TOO_MANY_CONNECTIONS.Msg,
				}, nil
			}

			return resp, nil
		}
	}
}
