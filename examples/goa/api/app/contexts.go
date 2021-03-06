// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "cellar": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/dikhan/terraform-provider-openapi/examples/goa/api/design
// --out=$(GOPATH)/src/github.com/dikhan/terraform-provider-openapi/examples/goa/api
// --version=v1.3.1

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// CreateBottleContext provides the bottle create action context.
type CreateBottleContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *BottlePayload
}

// NewCreateBottleContext parses the incoming request URL and body, performs validations and creates the
// context used by the bottle controller create action.
func NewCreateBottleContext(ctx context.Context, r *http.Request, service *goa.Service) (*CreateBottleContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := CreateBottleContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// Created sends a HTTP response with status code 201.
func (ctx *CreateBottleContext) Created(r *Bottle) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.gophercon.goa.bottle")
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 201, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *CreateBottleContext) BadRequest(r error) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *CreateBottleContext) InternalServerError() error {
	ctx.ResponseData.WriteHeader(500)
	return nil
}

// ShowBottleContext provides the bottle show action context.
type ShowBottleContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	ID string
}

// NewShowBottleContext parses the incoming request URL and body, performs validations and creates the
// context used by the bottle controller show action.
func NewShowBottleContext(ctx context.Context, r *http.Request, service *goa.Service) (*ShowBottleContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := ShowBottleContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramID := req.Params["id"]
	if len(paramID) > 0 {
		rawID := paramID[0]
		rctx.ID = rawID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowBottleContext) OK(r *Bottle) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.gophercon.goa.bottle")
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowBottleContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}
