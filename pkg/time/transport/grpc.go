package transport

import (
	"context"
	"errors"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time_summary"
	"github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
)

type grpcServer struct {
	createTime     grpctransport.Handler
	getTimeSummary grpctransport.Handler
	time_v1.UnimplementedTimeV1Server
}

// NewGRPCServer ???
func NewGRPCServer(ep endpoints.Set) time_v1.TimeV1Server {
	return &grpcServer{
		createTime: grpctransport.NewServer(
			ep.CreateTimeEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
		),
		getTimeSummary: grpctransport.NewServer(
			ep.GetTimeSummary,
			decodeGRPCGetSummaryRequest,
			encodeGRPCGetSummaryResponse,
		),
	}
}

func (g *grpcServer) Create(
	ctx context.Context,
	r *time_v1.CreateRequest,
) (*time_v1.CreateResponse, error) {
	_, resp, err := g.createTime.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*time_v1.CreateResponse), nil
}

func (g *grpcServer) GetSummary(
	ctx context.Context,
	r *time_v1.GetSummaryRequest,
) (*time_v1.GetSummaryResponse, error) {
	_, resp, err := g.getTimeSummary.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*time_v1.GetSummaryResponse), nil
}

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*time_v1.CreateRequest)
	t := time.Time{MacAddress: req.MacAddress, Seconds: req.Seconds, RouterID: req.RouterId}

	return endpoints.CreateTimeRequest{Time: &t}, nil
}

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(endpoints.CreateTimeResponse)

	return &time_v1.CreateResponse{
		MacAddress: res.Time.MacAddress,
		Seconds:    res.Time.Seconds,
		RouterId:   res.Time.RouterID,
	}, nil
}

func decodeGRPCGetSummaryRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*time_v1.GetSummaryRequest)

	filters := make([]time_summary.Filter, 0, len(req.Filters))
	for _, filter := range req.Filters {
		filters = append(filters, time_summary.Filter{Key: filter.Key, Value: filter.Value})
	}

	return endpoints.GetTimeSummaryRequest{Filters: filters}, nil
}

func encodeGRPCGetSummaryResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(endpoints.GetTimeSummaryResponse)

	if !ok {
		return nil, errors.New("invalid response body")
	}

	tsList := make([]*time_v1.Summary, 0, len(res.TimeSummary))
	for _, ts := range res.TimeSummary {
		breaks := make([]*time_v1.Break, 0, len(ts.Breaks))
		for _, b := range ts.Breaks {
			breaks = append(breaks, &time_v1.Break{
				SecondsStart: b.SecondsStart,
				SecondsEnd:   b.SecondsEnd,
			})
		}

		tsList = append(tsList, &time_v1.Summary{
			MacAddress:   ts.MacAddress,
			Seconds:      ts.Seconds,
			Breaks:       breaks,
			Date:         ts.Date,
			SecondsStart: ts.SecondsStart,
			SecondsEnd:   ts.SecondsEnd,
		})
	}

	return &time_v1.GetSummaryResponse{Summary: tsList}, nil
}
