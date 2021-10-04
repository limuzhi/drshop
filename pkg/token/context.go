package token

import (
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/gogf/gf/util/gconv"
)

type LoginInfoKey struct{}

type XmdGlobalUid struct{}

var Authorization = "x-md-global-token"

func WithLoginContext(ctx context.Context, loginInfo *UserClaims) context.Context {
	return context.WithValue(ctx, LoginInfoKey{}, loginInfo)
}

func FormLoginContext(ctx context.Context) *UserClaims {
	data := ctx.Value(LoginInfoKey{})
	if data != nil {
		return data.(*UserClaims)
	}
	return nil
}

func WithGlobalUidContext(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, XmdGlobalUid{}, userId)
}

func FormGlobalUidContext(ctx context.Context) int64 {
	data := ctx.Value(XmdGlobalUid{})
	if data != nil {
		return gconv.Int64(data)
	}
	return 0
}


func WithAuthorizationContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, Authorization, token)
}

func FormAuthorizationContext(ctx context.Context) string {
	data := ctx.Value(Authorization)
	if data != nil {
		return gconv.String(data)
	}
	return ""
}

func NewContextMetadataClientToken(ctx context.Context) context.Context {
	token := FormAuthorizationContext(ctx)
	outCtx := context.Background()
	if len(token) > 0 {
		outCtx = metadata.AppendToClientContext(outCtx, Authorization, token)
	}
	return outCtx
}

func MetadataToFromServerContext(ctx context.Context) string {
	var token string
	md, ok := metadata.FromServerContext(ctx)
	if ok && len(md.Get(Authorization)) > 0 {
		token = md.Get(Authorization)
	}
	return token
}
