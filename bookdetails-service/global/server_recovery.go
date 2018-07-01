package global

import (
	"net/http"
	"runtime/debug"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

type HttpServerPanicHandler struct {
	http.Handler
}

func (h HttpServerPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			Logger.Errorln(e, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Keep calm! Try again after prayer."))
		}
	}()
	h.Handler.ServeHTTP(w, r)
}

func GrpcServerPanicHandlerOptions() []grpc.ServerOption{
	var recoveryOpts = []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			Logger.Errorln(p, string(debug.Stack()))
			return nil
		}),
	}

	var serverOpts = []grpc.ServerOption{

		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	}

	return serverOpts
}
