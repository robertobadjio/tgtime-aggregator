package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/oklog/oklog/pkg/group"
	"github.com/robertobadjio/tgtime-aggregator/internal/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/config"
	"github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	timeApp "github.com/robertobadjio/tgtime-aggregator/pkg/time"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/transport"
)

var logger log.Logger

func main() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		_ = logger.Log("config", "grpc", "error", err.Error())
		os.Exit(1)
	}

	httpConfig, err := config.NewHTTPConfig()
	if err != nil {
		_ = logger.Log("config", "http", "error", err.Error())
		os.Exit(1)
	}

	go aggregator.Aggregate()
	//go sendPreviousDayInfo(logger)

	var (
		s           = timeApp.NewService()
		eps         = endpoints.NewEndpointSet(s)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	// API Gateway
	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpConfig.Address())
		if err != nil {
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "HTTP", "addr", httpConfig.Address())
			srv := &http.Server{
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				Handler:      httpHandler,
			}
			return srv.Serve(httpListener)
		}, func(err error) {
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			_ = httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", grpcConfig.Address())
		if err != nil {
			_ = logger.Log("transport", "GRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "GRPC", "addr", grpcConfig.Address())
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			reflection.Register(baseServer)
			time_v1.RegisterTimeV1Server(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			_ = grpcListener.Close()
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
	_ = logger.Log("exit", g.Run())
}

/*
func sendPreviousDayInfo(logger log.Logger) {
	r := tsPgRepo.NewPgRepository(db.GetDB())
	tsService := implementationTs.NewTimeSummaryService(r, logger)

	cfg := config.New()
	ctx := context.Background()

	previousDate := time.Now().AddDate(0, 0, -1)
	filters := make([]*time_summary.Filter, 0, 1)
	filters = append(filters, &time_summary.Filter{Key: "date", Value: previousDate.Format("2006-01-02")})
	ts, _ := tsService.GetTimeSummary(ctx, filters) // TODO: Handle error

	k := kafka.NewKafka(cfg.KafkaHost + ":" + cfg.KafkaPort)
	for _, t := range ts {
		breaks := make([]*kafka.Break, 0, len(t.Breaks))
		for _, b := range t.Breaks {
			breaks = append(breaks, &kafka.Break{
				SecondsStart: b.SecondsStart,
				SecondsEnd:   b.SecondsEnd,
			})
		}

		m := kafka.PreviousDayInfoMessage{
			MacAddress:   t.MacAddress,
			Seconds:      t.Seconds,
			Breaks:       breaks,
			Date:         t.Date,
			SecondsStart: t.SecondsStart,
			SecondsEnd:   t.SecondsEnd,
		}
		_ = k.Produce(ctx, m, kafka.PreviousDayInfoTopic) // TODO: Handle error
	}
}
*/
