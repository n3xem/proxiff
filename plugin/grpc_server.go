package plugin

import (
	"context"
	"net/http"

	"github.com/n3xem/proxiff/comparator"
	"github.com/n3xem/proxiff/plugin/proto"
)

// GRPCServer is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	proto.UnimplementedComparatorPluginServer
	// This is the real implementation
	Impl comparator.Comparator
}

func (m *GRPCServer) Compare(ctx context.Context, req *proto.CompareRequest) (*proto.CompareResponse, error) {
	// Convert proto.HTTPResponse to comparator.Response
	newer := protoToResponse(req.Newer)
	current := protoToResponse(req.Current)

	// Call the actual implementation
	result := m.Impl.Compare(newer, current)

	return &proto.CompareResponse{
		Match:      result.Match,
		Difference: result.Difference,
	}, nil
}

func protoToResponse(p *proto.HTTPResponse) *comparator.Response {
	headers := make(http.Header)
	for key, values := range p.Headers {
		headers[key] = values.Values
	}

	return &comparator.Response{
		StatusCode: int(p.StatusCode),
		Headers:    headers,
		Body:       p.Body,
	}
}
