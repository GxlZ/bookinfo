package server

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	// 3d Party
	"google.golang.org/grpc"

	// This Service
	"bookinfo/bookdetails-service/handlers"
	"bookinfo/bookdetails-service/svc"
	pb "bookinfo/pb/details"
	"bookinfo/bookdetails-service/global"

	_ "github.com/mkevac/debugcharts"
	_ "expvar"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"expvar"
	"runtime"
	"time"
)

var startTime = time.Now()

func NewEndpoints() svc.Endpoints {
	// Business domain.
	var service pb.BookDetailsServer
	{
		service = handlers.NewService()
		// Wrap Service with middlewares. See handlers/middlewares.go
		service = handlers.WrapService(service)
	}

	// Endpoint domain.
	var (
		detailEndpoint = svc.MakeDetailEndpoint(service)
	)

	endpoints := svc.Endpoints{
		DetailEndpoint: detailEndpoint,
	}

	// Wrap selected Endpoints with middlewares. See handlers/middlewares.go
	endpoints = handlers.WrapEndpoints(endpoints)

	return endpoints
}

// Run starts a new http server, gRPC server, and a debug server with the
// passed config and logger
func Run() {
	endpoints := NewEndpoints()

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go handlers.InterruptHandler(errc)

	// Debug listener.
	go func() {
		expvar.Publish("Goroutines", expvar.Func(func() interface{}{
			return runtime.NumGoroutine()
		}))
		expvar.Publish("Uptime", expvar.Func(func() interface{}{
			uptime := time.Since(startTime)
			return int64(uptime)
		}))

		log.Println("transport", "debug", "addr", global.Conf.DebugServer.Addr)

		errc <- http.ListenAndServe(global.Conf.DebugServer.Addr, nil)
	}()

	// HTTP transport.
	go func() {
		log.Println("transport", "HTTP", "addr", global.Conf.HttpServer.Addr)
		h := svc.MakeHTTPHandler(endpoints)
		errc <- http.ListenAndServe(global.Conf.HttpServer.Addr, h)
	}()

	//prometheus mertics
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/metrics",
			promhttp.Handler().ServeHTTP,
		)

		log.Println("transport", "instrumenting", "adapter", "prometheus", "addr", global.Conf.MetricsServer.Addr)

		errc <- http.ListenAndServe(global.Conf.MetricsServer.Addr, mux)
	}()

	// gRPC transport.
	go func() {
		log.Println("transport", "gRPC", "addr", global.Conf.GrpcServer.Addr)
		ln, err := net.Listen("tcp", global.Conf.GrpcServer.Addr)
		if err != nil {
			errc <- err
			return
		}

		srv := svc.MakeGRPCServer(endpoints)
		s := grpc.NewServer()
		pb.RegisterBookDetailsServer(s, srv)

		errc <- s.Serve(ln)
	}()

	// Run!
	log.Println("exit", <-errc)
}
