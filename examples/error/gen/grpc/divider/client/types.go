// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// divider gRPC client types
//
// Command:
// $ goa gen goa.design/goa/examples/error/design -o
// $(GOPATH)/src/goa.design/goa/examples/error

package client

import (
	dividersvc "goa.design/goa/examples/error/gen/divider"
	dividerpb "goa.design/goa/examples/error/gen/grpc/divider"
)

// NewIntOperands builds the gRPC request type from the payload of the
// "integer_divide" endpoint of the "divider" service.
func NewIntOperands(p *dividersvc.IntOperands) *dividerpb.IntOperands {
	v := &dividerpb.IntOperands{
		A: int32(p.A),
		B: int32(p.B),
	}
	return v
}

// NewFloatOperands builds the gRPC request type from the payload of the
// "divide" endpoint of the "divider" service.
func NewFloatOperands(p *dividersvc.FloatOperands) *dividerpb.FloatOperands {
	v := &dividerpb.FloatOperands{
		A: p.A,
		B: p.B,
	}
	return v
}