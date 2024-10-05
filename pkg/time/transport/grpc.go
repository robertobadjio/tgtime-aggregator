package transport

import (
	"context"
	"errors"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	aggregatorsvc "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
)

type grpcServer struct {
	createTime     grpctransport.Handler
	getTimeSummary grpctransport.Handler
	aggregatorsvc.UnimplementedAggregatorServer
}

func NewGRPCServer(ep endpoints.Set) aggregatorsvc.AggregatorServer {
	return &grpcServer{
		createTime: grpctransport.NewServer(
			ep.CreateTimeEndpoint,
			decodeGRPCCreateTimeRequest,
			encodeGRPCCreateTimeResponse,
		),
		getTimeSummary: grpctransport.NewServer(
			ep.GetTimeSummary,
			decodeGRPCGetTimeSummaryRequest,
			encodeGRPCGetTimeSummaryResponse,
		),
	}
}

func (g *grpcServer) CreateTime(
	ctx context.Context,
	r *aggregatorsvc.CreateTimeRequest,
) (*aggregatorsvc.CreateTimeResponse, error) {
	_, resp, err := g.createTime.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*aggregatorsvc.CreateTimeResponse), nil
}

func (g *grpcServer) GetTimeSummary(
	ctx context.Context,
	r *aggregatorsvc.GetTimeSummaryRequest,
) (*aggregatorsvc.GetTimeSummaryResponse, error) {
	_, resp, err := g.getTimeSummary.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*aggregatorsvc.GetTimeSummaryResponse), nil
}

func decodeGRPCCreateTimeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*aggregatorsvc.CreateTimeRequest)
	t := time.TimeUser{MacAddress: req.MacAddress, Seconds: req.Seconds, RouterId: int8(req.RouterId)}

	return endpoints.CreateTimeRequest{Time: &t}, nil
}

func encodeGRPCCreateTimeResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(endpoints.CreateTimeResponse)

	return &aggregatorsvc.CreateTimeResponse{
		MacAddress: res.Time.MacAddress,
		Seconds:    res.Time.Seconds,
		RouterId:   int64(res.Time.RouterId),
	}, nil
}

func decodeGRPCGetTimeSummaryRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*aggregatorsvc.GetTimeSummaryRequest)

	filters := make([]*time_summary.Filter, 0, len(req.Filters)) // TODO: Неэффективная инициализация среза #21
	for _, filter := range req.Filters {
		filters = append(filters, &time_summary.Filter{Key: filter.Key, Value: filter.Value})
	}

	return endpoints.GetTimeSummaryRequest{Filters: filters}, nil
}

func encodeGRPCGetTimeSummaryResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(endpoints.GetTimeSummaryResponse)

	if !ok {
		return nil, errors.New("invalid response body")
	}

	tsList := make([]*aggregatorsvc.TimeSummary, 0, len(res.TimeSummary)) // TODO: Неэффективная инициализация среза #21
	for _, ts := range res.TimeSummary {
		tsList = append(tsList, &aggregatorsvc.TimeSummary{
			MacAddress:   ts.MacAddress,
			Seconds:      ts.Seconds,
			BreaksJson:   string(ts.BreaksJson),
			Date:         ts.Date,
			SecondsStart: ts.SecondsStart,
			SecondsEnd:   ts.SecondsEnd,
		})
	}

	return &aggregatorsvc.GetTimeSummaryResponse{TimeSummary: tsList}, nil
}
