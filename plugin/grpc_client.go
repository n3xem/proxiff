package plugin

import (
	"context"

	"github.com/n3xem/proxiff/comparator"
	"github.com/n3xem/proxiff/plugin/proto"
)

// GRPCClient is an implementation of comparator.Comparator that talks over RPC.
type GRPCClient struct {
	client proto.ComparatorPluginClient
}

func (m *GRPCClient) Compare(newer, current *comparator.Response) *comparator.Result {
	// Convert comparator.Response to proto.HTTPResponse
	newerProto := responseToProto(newer)
	currentProto := responseToProto(current)

	// Call the plugin
	resp, err := m.client.Compare(context.Background(), &proto.CompareRequest{
		Newer:   newerProto,
		Current: currentProto,
	})

	if err != nil {
		return &comparator.Result{
			Match:      false,
			Newer:      newer,
			Current:    current,
			Difference: "plugin error: " + err.Error(),
		}
	}

	return &comparator.Result{
		Match:      resp.Match,
		Newer:      newer,
		Current:    current,
		Difference: resp.Difference,
	}
}

func responseToProto(r *comparator.Response) *proto.HTTPResponse {
	headers := make(map[string]*proto.HeaderValues)
	for key, values := range r.Headers {
		headers[key] = &proto.HeaderValues{Values: values}
	}

	return &proto.HTTPResponse{
		StatusCode: int32(r.StatusCode),
		Headers:    headers,
		Body:       r.Body,
	}
}
