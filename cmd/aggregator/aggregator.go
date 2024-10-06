package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/oklog/oklog/pkg/group"
	aggregatorsvc "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/config"
	"github.com/robertobadjio/tgtime-aggregator/internal/db"
	implementationT "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/implementation"
	tPgRepo "github.com/robertobadjio/tgtime-aggregator/internal/domain/time/pg_db"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	implementationTs "github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary/implementation"
	tsPgRepo "github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary/pg_db"
	"github.com/robertobadjio/tgtime-aggregator/internal/kafka"
	timeApp "github.com/robertobadjio/tgtime-aggregator/pkg/time"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/transport"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Db *sql.DB

func main() {
	cfg := config.New()

	var err error
	dbConn := db.GetDB()
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}

	err = m.Run()
	if err != nil {
		panic(err)
	}

	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("", cfg.HttpPort)
		grpcAddr = net.JoinHostPort("", cfg.GrpcPort)
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	go aggregate(logger)
	go sendPreviousDayInfo(logger)

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
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(err error) {
			httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			_ = logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "gRPC", "addr", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			aggregatorsvc.RegisterAggregatorServer(baseServer, grpcServer)
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
		m := kafka.PreviousDayInfoMessage{
			MacAddress:   t.MacAddress,
			Seconds:      t.Seconds,
			BreaksJson:   t.BreaksJson,
			Date:         t.Date,
			SecondsStart: t.SecondsStart,
			SecondsEnd:   t.SecondsEnd,
		}
		_ = k.Produce(ctx, m, kafka.PreviousDayInfoTopic) // TODO: Handle error
	}
}

// TODO: Горизонтальное масштабирование
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
		macAddresses, _ := tService.GetMacAddresses(ctx, getDate("Europe/Moscow")) // TODO: Handle error
		for _, macAddress := range macAddresses {
			timeSummary, err := agr.AggregateTime(ctx, macAddress)
			if err != nil {
				_ = logger.Log("msg", err.Error())
				continue
			}

			err = tsService.CreateTimeSummary(ctx, timeSummary)
			if err != nil {
				_ = logger.Log("msg", err.Error())
			}
		}
	}
}

func getDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}
