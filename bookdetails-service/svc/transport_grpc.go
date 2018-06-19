package svc

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "bookinfo/pb/details"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC BookDetailsServer.
func MakeGRPCServer(endpoints Endpoints) pb.BookDetailsServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(metadataToContext),
	}
	return &grpcServer{
		// bookdetails

		detail: grpctransport.NewServer(
			endpoints.DetailEndpoint,
			DecodeGRPCDetailRequest,
			EncodeGRPCDetailResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the BookDetailsServer interface
type grpcServer struct {
	detail grpctransport.Handler
}

// Methods for grpcServer to implement BookDetailsServer interface

func (s *grpcServer) Detail(ctx context.Context, req *pb.DetailReq) (*pb.DetailResp, error) {
	_, rep, err := s.detail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DetailResp), nil
}

// Server Decode

// DecodeGRPCDetailRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC detail request to a user-domain detail request. Primarily useful in a server.
func DecodeGRPCDetailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DetailReq)
	return req, nil
}

// Server Encode

// EncodeGRPCDetailResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain detail response to a gRPC detail reply. Primarily useful in a server.
func EncodeGRPCDetailResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.DetailResp)
	return resp, nil
}

// Helpers

func metadataToContext(ctx context.Context, md metadata.MD) context.Context {
	for k, v := range md {
		if v != nil {
			// The key is added both in metadata format (k) which is all lower
			// and the http.CanonicalHeaderKey of the key so that it can be
			// accessed in either format
			ctx = context.WithValue(ctx, k, v[0])
			ctx = context.WithValue(ctx, http.CanonicalHeaderKey(k), v[0])
		}
	}

	return ctx
}
