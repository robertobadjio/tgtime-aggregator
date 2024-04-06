package transport

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	aggregatorsvc "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"github.com/robertobadjio/tgtime-aggregator/internal/domain/time"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
)

type grpcServer struct {
	createTime grpctransport.Handler
	aggregatorsvc.UnimplementedAggregatorServer
}

func NewGRPCServer(ep endpoints.Set) aggregatorsvc.AggregatorServer {
	return &grpcServer{
		createTime: grpctransport.NewServer(
			ep.CreateTimeEndpoint,
			decodeGRPCCreateTimeRequest,
			encodeGRPCCreateTimeResponse,
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

func decodeGRPCCreateTimeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*aggregatorsvc.CreateTimeRequest)
	u := time.TimeUser{MacAddress: req.MacAddress, Seconds: req.Seconds, RouterId: int8(req.RouterId)}

	return endpoints.CreateTimeRequest{Time: &u}, nil
}

func encodeGRPCCreateTimeResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(endpoints.CreateTimeResponse)

	return &aggregatorsvc.CreateTimeResponse{MacAddress: res.Time.MacAddress, Seconds: res.Time.Seconds, RouterId: int64(res.Time.RouterId)}, nil
}
