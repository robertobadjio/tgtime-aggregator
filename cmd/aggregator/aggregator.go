package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	aggregatorsvc "tgtime-aggregator/api/v1/pb/aggregator"
	"tgtime-aggregator/internal/aggregator"
	"tgtime-aggregator/internal/config"
	implementationT "tgtime-aggregator/internal/domain/time/implementation"
	tPgRepo "tgtime-aggregator/internal/domain/time/pg_db"
	implementationTs "tgtime-aggregator/internal/domain/time_summary/implementation"
	tsPgRepo "tgtime-aggregator/internal/domain/time_summary/pg_db"
	"tgtime-aggregator/internal/tgtime_api_client"
	timeApp "tgtime-aggregator/pkg/time"
	"tgtime-aggregator/pkg/time/endpoints"
	"tgtime-aggregator/pkg/time/transport"
	"time"
)

var Db *sql.DB

func main() {
	cfg := config.New()

	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("", cfg.HttpPort)
		grpcAddr = net.JoinHostPort("", cfg.GrpcPort)
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	go aggregate(logger)

	var (
		s           = timeApp.NewService()
		eps         = endpoints.NewEndpointSet(s)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	// API Gateway
	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(err error) {
			httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			aggregatorsvc.RegisterAggregatorServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
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
	logger.Log("exit", g.Run())
}

func aggregate(logger log.Logger) {
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

		tRepo := tPgRepo.NewPgRepository(Db)
		tService := implementationT.NewTimeService(tRepo, logger)

		tsRepo := tsPgRepo.NewPgRepository(Db)
		tsService := implementationTs.NewTimeSummaryService(tsRepo, logger)

		agr := aggregator.NewAggregator(getDate("Europe/Moscow"), tService)
		ctx := context.TODO()
		users := getUsers()
		for _, user := range users.Users {
			timeSummary, err := agr.AggregateTime(ctx, user)
			if err != nil {
				logger.Log("msg", err.Error())
				continue
			}

			err = tsService.CreateTimeSummary(ctx, timeSummary)
			if err != nil {
				logger.Log("msg", err.Error())
			}
		}
	}
}

func getDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}

func getUsers() *tgtime_api_client.Users {
	apiClient := tgtime_api_client.NewTgTimeClient()
	users, _ := apiClient.GetAllUsers()

	return users
}
