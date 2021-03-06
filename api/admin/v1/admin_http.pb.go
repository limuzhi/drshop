// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.0.5

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type AdminServiceHTTPServer interface {
	CaptchaImg(context.Context, *CaptchaImgReq) (*CaptchaImgReply, error)
}

func RegisterAdminServiceHTTPServer(s *http.Server, srv AdminServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/admin/v1/captcha", _AdminService_CaptchaImg0_HTTP_Handler(srv))
}

func _AdminService_CaptchaImg0_HTTP_Handler(srv AdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CaptchaImgReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.admin.v1.AdminService/CaptchaImg")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CaptchaImg(ctx, req.(*CaptchaImgReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CaptchaImgReply)
		return ctx.Result(200, reply)
	}
}

type AdminServiceHTTPClient interface {
	CaptchaImg(ctx context.Context, req *CaptchaImgReq, opts ...http.CallOption) (rsp *CaptchaImgReply, err error)
}

type AdminServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewAdminServiceHTTPClient(client *http.Client) AdminServiceHTTPClient {
	return &AdminServiceHTTPClientImpl{client}
}

func (c *AdminServiceHTTPClientImpl) CaptchaImg(ctx context.Context, in *CaptchaImgReq, opts ...http.CallOption) (*CaptchaImgReply, error) {
	var out CaptchaImgReply
	pattern := "/admin/v1/captcha"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.admin.v1.AdminService/CaptchaImg"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
