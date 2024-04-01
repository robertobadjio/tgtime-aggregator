package main

import (
	"fmt"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	aggregatorsvc "tgtime-aggregator/api/v1/pb/aggregator"
	"tgtime-aggregator/internal/aggregator"
	"tgtime-aggregator/internal/config"
	timeApp "tgtime-aggregator/pkg/time"
	"tgtime-aggregator/pkg/time/endpoints"
	"tgtime-aggregator/pkg/time/transport"
	"time"
)

func main() {
	cfg := config.New()
	go aggregate()

	var (
		//logger   log.Logger
		httpAddr = net.JoinHostPort("", cfg.HttpPort)
		grpcAddr = net.JoinHostPort("", cfg.GrpcPort)
	)

	var (
		s           = timeApp.NewService()
		eps         = endpoints.NewEndpointSet(s)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	// API Gateway
	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		g.Add(func() error {
			log.Printf("Serving http address %s", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(err error) {
			httpListener.Close()
		})
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatal(err)
		}
		g.Add(func() error {
			log.Printf("Serving grpc address %s", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			aggregatorsvc.RegisterAggregatorServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}

func aggregate() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 1, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	for {
		time.Sleep(d)
		d = 24 * time.Hour
		aggregator.AggregateTime()
	}
}
