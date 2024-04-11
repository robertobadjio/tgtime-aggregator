package transport

import (
	"context"
	"errors"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	aggregatorsvc "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
)

type grpcServer struct {
	createTime              grpctransport.Handler
	getTimeSummaryByDate    grpctransport.Handler
	getTimeSummaryAllByDate grpctransport.Handler
	aggregatorsvc.UnimplementedAggregatorServer
}

func NewGRPCServer(ep endpoints.Set) aggregatorsvc.AggregatorServer {
	return &grpcServer{
		createTime: grpctransport.NewServer(
			ep.CreateTimeEndpoint,
			decodeGRPCCreateTimeRequest,
			encodeGRPCCreateTimeResponse,
		),
		getTimeSummaryByDate: grpctransport.NewServer(
			ep.GetTimeSummaryByDate,
			decodeGRPCGetTimeSummaryByDateRequest,
			encodeGRPCGetTimeSummaryByDateResponse,
		),
		getTimeSummaryAllByDate: grpctransport.NewServer(
			ep.GetTimeSummaryAllByDate,
			decodeGRPCGetTimeSummaryAllByDateRequest,
			encodeGRPCGetTimeSummaryAllByDateResponse,
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

func (g *grpcServer) GetTimeSummaryByDate(
	ctx context.Context,
	r *aggregatorsvc.GetTimeSummaryByDateRequest,
) (*aggregatorsvc.GetTimeSummaryByDateResponse, error) {
	_, resp, err := g.getTimeSummaryByDate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*aggregatorsvc.GetTimeSummaryByDateResponse), nil
}

func (g *grpcServer) GetTimeSummaryAllByDate(
	ctx context.Context,
	r *aggregatorsvc.GetTimeSummaryAllByDateRequest,
) (*aggregatorsvc.GetTimeSummaryAllByDateResponse, error) {
	_, resp, err := g.getTimeSummaryAllByDate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*aggregatorsvc.GetTimeSummaryAllByDateResponse), nil
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

func decodeGRPCGetTimeSummaryByDateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*aggregatorsvc.GetTimeSummaryByDateRequest)

	return endpoints.GetTimeSummaryByDateRequest{MacAddress: req.MacAddress, Date: req.Date}, nil
}

func encodeGRPCGetTimeSummaryByDateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(endpoints.GetTimeSummaryByDateResponse)

	if !ok {
		return nil, errors.New("invalid response body")
	}

	ts := aggregatorsvc.TimeSummary{
		MacAddress:   res.TimeSummary.MacAddress,
		Seconds:      res.TimeSummary.Seconds,
		BreaksJson:   string(res.TimeSummary.BreaksJson),
		Date:         res.TimeSummary.Date,
		SecondsStart: res.TimeSummary.SecondsStart,
		SecondsEnd:   res.TimeSummary.SecondsEnd,
	}

	return &aggregatorsvc.GetTimeSummaryByDateResponse{TimeSummary: &ts}, nil
}

func decodeGRPCGetTimeSummaryAllByDateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*aggregatorsvc.GetTimeSummaryAllByDateRequest)

	return endpoints.GetTimeSummaryAllByDateRequest{Date: req.Date}, nil
}

func encodeGRPCGetTimeSummaryAllByDateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(endpoints.GetTimeSummaryAllByDateResponse)

	if !ok {
		return nil, errors.New("invalid response body")
	}

	tsList := make([]*aggregatorsvc.TimeSummary, 0, len(res.Data)) // TODO: Неэффективная инициализация среза #21
	for _, ts := range res.Data {
		tsList = append(tsList, &aggregatorsvc.TimeSummary{
			MacAddress:   ts.MacAddress,
			Seconds:      ts.Seconds,
			BreaksJson:   string(ts.BreaksJson),
			Date:         ts.Date,
			SecondsStart: ts.SecondsStart,
			SecondsEnd:   ts.SecondsEnd,
		})
	}

	return &aggregatorsvc.GetTimeSummaryAllByDateResponse{TimeSummary: tsList}, nil
}
