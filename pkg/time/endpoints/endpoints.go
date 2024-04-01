package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"tgtime-aggregator/pkg/time"
)

type Set struct {
	CreateTimeEndpoint    endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc time.Service) Set {
	return Set{
		CreateTimeEndpoint:    MakeCreateTimeEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
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

func MakeServiceStatusEndpoint(svc time.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code := svc.ServiceStatus(ctx)
		return ServiceStatusResponse{Code: code}, nil
	}
}
