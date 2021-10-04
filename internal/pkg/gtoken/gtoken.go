package gtoken

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

type GToken struct {
	// server name
	ServerName string
	// 缓存key
	CacheKey string
	// 令牌有效的持续时间。可选，默认为一小时。
	Timeout time.Duration
	//此字段允许客户端刷新其令牌，直到通过MaxRefresh。
	//请注意，客户端可以在MaxRefresh的最后一刻刷新其令牌。
	//这意味着令牌的最大有效时间跨度为TokenTime+MaxRefresh。
	//可选，默认为0表示不可刷新。
	MaxRefresh time.Duration
	// Callback function that should perform the authentication of the user based on login info.
	// Must return user data as user identifier, it will be stored in Claim Array. Required.
	// Check error (e) to determine the appropriate error message.
	Authenticator func(ctx context.Context) (interface{}, error)

	Authorizator func(data interface{}, c *gin.Context) bool
}
