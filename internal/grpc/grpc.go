package grpc

import (
	"context"
	"github.com/hyuti/api-blueprint/internal/proto"
	"github.com/hyuti/api-blueprint/internal/usecase"
	"github.com/hyuti/api-blueprint/pkg/collection"
)

type apiGolangTemplateServiceServer struct {
	proto.UnimplementedApiGolangTemplateServer
	uc1 usecase.ExampleUseCase
}

var _ proto.ApiGolangTemplateServer = (*apiGolangTemplateServiceServer)(nil)

func New(
	uc1 usecase.ExampleUseCase,
) proto.ApiGolangTemplateServer {
	return &apiGolangTemplateServiceServer{
		uc1: uc1,
	}
}

func (p *apiGolangTemplateServiceServer) ListExample(ctx context.Context, request *proto.ExampleListRequest) (*proto.ExampleListResponse, error) {
	resp, err := p.uc1.List(ctx, &usecase.ExampleReq{
		PaginatedRequest: usecase.PaginatedRequest{
			PageSize: request.PageSize,
			Page:     request.Page,
		},
	})
	if err != nil {
		return nil, usecaseGrpcMapper(err)
	}
	paginatedResp := resp
	return &proto.ExampleListResponse{
		Next:     paginatedResp.Next,
		Previous: paginatedResp.Prev,
		PageSize: paginatedResp.PageSize,
		Count:    paginatedResp.Count,
		Data: collection.Map(resp.Data, func(item usecase.Example, _ int) *proto.Example {
			return item.Example
		}),
	}, nil
}
