package token

import (
	"context"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"log"
	"strings"
	"time"
)

var (
	ErrAuthFail         = errors.New(401, "Authentication failed", "Missing token or token incorrect")
	ErrTokenExpiresFail = errors.New(401, "Authentication failed", "Missing token or token expires")
)

func AuthTokenGinMiddleware(tk *Token) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if cc, ok := kgin.FromGinContext(ctx); ok {
				var jwtToken string
				authKey := "Authorization"
				authHeader := cc.Request.Header.Get(authKey)
				if authHeader != "" {
					parts := strings.SplitN(authHeader, " ", 2)
					if len(parts) > 1 && parts[0] == "Bearer" && parts[1] != "" {
						jwtToken = parts[1]
					}
				}
				if jwtToken != "" {
					userInfo, _ := tk.ParseToken(jwtToken)
					if userInfo != nil && time.Now().Unix() < userInfo.ExpiresAt {
						ctx = WithLoginContext(cc.Request.Context(), userInfo)
						cc.Request = cc.Request.WithContext(ctx)
						reply, err = handler(cc.Request.Context(), cc.Request)
					}
				}
			}
			return
		}
	}
}

// 根据header读取jwt token
func AuthMiddleware(tk *Token) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var token string
			tr, ok := transport.FromServerContext(ctx)
			if ok {
				authorization := tr.RequestHeader().Get("authorization")
				if authorization != "" {
					token = authorization
				}
			}
			if tr.Kind() == transport.KindGRPC {
				token = MetadataToFromServerContext(ctx)
				log.Println("=======transport.KindGRPC=====22====", token)
			}
			log.Println("==token==Authorization===", token)
			//if token == "" {
			//	return nil, errors.BadRequest("token不得为空", "token不得为空")
			//}
			if token != "" {
				// 解析token
				claims, parseErr := tk.ParseToken(token)
				log.Println("========ParseToken======", claims, "======")
				if parseErr != nil {
					return nil, errors.BadRequest("token解析错误", "token解析错误")
				} else if time.Now().Unix() > claims.ExpiresAt {
					return nil, errors.BadRequest("token已过期", "token已过期")
				}
				if cc, ok := kgin.FromGinContext(ctx); ok {
					ctx = WithAuthorizationContext(cc.Request.Context(), token)
					cc.Request = cc.Request.WithContext(ctx)

					ctx = WithLoginContext(cc.Request.Context(), claims)
					cc.Request = cc.Request.WithContext(ctx)

					ctx = WithGlobalUidContext(cc.Request.Context(), claims.UserID)
					cc.Request = cc.Request.WithContext(ctx)
					return handler(cc.Request.Context(), cc.Request)
				} else {
					ctx = WithAuthorizationContext(ctx, token)
					ctx = WithLoginContext(ctx, claims)
					ctx = WithGlobalUidContext(ctx, claims.UserID)
				}
			}
			return handler(ctx, req)
		}
	}
}
