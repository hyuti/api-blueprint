package grpc

import (
	"context"
	"github.com/hyuti/API-Golang-Template/internal/example/proto"
	"github.com/hyuti/API-Golang-Template/internal/example/usecase"
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
		PageSize: request.PageSize,
		Page:     request.Page,
		Search:   request.Search,
	})
	if err != nil {
		return nil, usecaseGrpcMapper(err)
	}
	paginatedResp := resp.PaginatedResponse
	return &proto.ExampleListResponse{
		Next:     paginatedResp.Next,
		Previous: paginatedResp.Prev,
		PageSize: paginatedResp.PageSize,
		Count:    paginatedResp.Count,
		Data:     paginatedResp.Data,
	}, nil
}
