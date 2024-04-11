package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time"
)

type Set struct {
	CreateTimeEndpoint      endpoint.Endpoint
	GetTimeSummaryByDate    endpoint.Endpoint
	GetTimeSummaryAllByDate endpoint.Endpoint
	ServiceStatusEndpoint   endpoint.Endpoint
}

func NewEndpointSet(svc time.Service) Set {
	return Set{
		CreateTimeEndpoint:      MakeCreateTimeEndpoint(svc),
		GetTimeSummaryByDate:    MakeGetTimeSummaryByDateEndpoint(svc),
		GetTimeSummaryAllByDate: MakeGetTimeSummaryAllByDateEndpoint(svc),
		ServiceStatusEndpoint:   MakeServiceStatusEndpoint(svc),
	}
}

func MakeCreateTimeEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTimeRequest)

		t, err := svc.CreateTime(ctx, req.Time)
		if err != nil {
			return nil, err
		}
		return CreateTimeResponse{Time: t}, nil
	}
}

func MakeGetTimeSummaryByDateEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTimeSummaryByDateRequest)

		ts, err := svc.GetTimeSummaryByDate(ctx, req.MacAddress, req.Date)
		if err != nil {
			return GetTimeSummaryByDateResponse{TimeSummary: nil}, err
		}

		return GetTimeSummaryByDateResponse{TimeSummary: ts}, nil
	}
}

func MakeGetTimeSummaryAllByDateEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTimeSummaryAllByDateRequest)

		ts, err := svc.GetTimeSummaryAllByDate(ctx, req.Date)
		if err != nil {
			return GetTimeSummaryAllByDateResponse{Data: nil}, err
		}

		return GetTimeSummaryAllByDateResponse{Data: ts}, nil
	}
}

func MakeServiceStatusEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code := svc.ServiceStatus(ctx)
		return ServiceStatusResponse{Code: code}, nil
	}
}
