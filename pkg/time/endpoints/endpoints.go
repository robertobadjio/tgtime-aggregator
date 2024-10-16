package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time"
)

// Set ???
type Set struct {
	CreateTimeEndpoint    endpoint.Endpoint
	GetTimeSummary        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

// NewEndpointSet ???
func NewEndpointSet(svc time.Service) Set {
	return Set{
		CreateTimeEndpoint:    MakeCreateTimeEndpoint(svc),
		GetTimeSummary:        MakeGetTimeSummaryEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

// MakeCreateTimeEndpoint ???
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

// MakeGetTimeSummaryEndpoint ???
func MakeGetTimeSummaryEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTimeSummaryRequest)

		ts, err := svc.GetTimeSummary(ctx, req.Filters)
		if err != nil {
			return GetTimeSummaryResponse{TimeSummary: nil}, err
		}

		return GetTimeSummaryResponse{TimeSummary: ts}, nil
	}
}

// MakeServiceStatusEndpoint ???
func MakeServiceStatusEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code := svc.ServiceStatus(ctx)
		return ServiceStatusResponse{Code: code}, nil
	}
}
