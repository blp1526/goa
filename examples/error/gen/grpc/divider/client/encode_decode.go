// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// divider gRPC client encoders and decoders
//
// Command:
// $ goa gen goa.design/goa/examples/error/design -o
// $(GOPATH)/src/goa.design/goa/examples/error

package client

import (
	"context"

	dividersvc "goa.design/goa/examples/error/gen/divider"
	dividerpb "goa.design/goa/examples/error/gen/grpc/divider"
)

// EncodeIntegerDivideRequest encodes requests sent to divider integer_divide
// endpoint.
func EncodeIntegerDivideRequest(ctx context.Context, p *dividersvc.IntOperands) (context.Context, *dividerpb.IntegerDivideRequest) {
	req := NewIntegerDivideRequest(p)
	return ctx, req
}

// DecodeIntegerDivideResponse decodes responses from the divider
// integer_divide endpoint.
func DecodeIntegerDivideResponse(ctx context.Context, resp *dividerpb.IntegerDivideResponse) (int, error) {
	res := NewIntegerDivideResponse(resp)
	return res, nil
}

// EncodeDivideRequest encodes requests sent to divider divide endpoint.
func EncodeDivideRequest(ctx context.Context, p *dividersvc.FloatOperands) (context.Context, *dividerpb.DivideRequest) {
	req := NewDivideRequest(p)
	return ctx, req
}

// DecodeDivideResponse decodes responses from the divider divide endpoint.
func DecodeDivideResponse(ctx context.Context, resp *dividerpb.DivideResponse) (float64, error) {
	res := NewDivideResponse(resp)
	return res, nil
}