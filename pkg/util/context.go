package util

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/metadata"
	"strings"
)

func GetClientIp(ctx context.Context) string {
	clientIP := ""
	if tr, ok := transport.FromServerContext(ctx); ok {
		switch tr.Kind() {
		case transport.KindGRPC:
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				rips := md.Get("x-real-ip")
				if len(rips) > 0 {
					clientIP = rips[0]
				}
			}
		case transport.KindHTTP:
			clientIP = tr.RequestHeader().Get("X-Forwarded-For")
			clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
			if clientIP == "" {
				clientIP = strings.TrimSpace(tr.RequestHeader().Get("X-Real-Ip"))
			}
			if clientIP == "" {
				if addr := tr.RequestHeader().Get("X-Appengine-Remote-Addr"); addr != "" {
					clientIP = addr
				}
			}
			if clientIP == "" {
				if addr := tr.RequestHeader().Get("X-Appengine-Remote-Addr"); addr != "" {
					clientIP = addr
				}
			}
		}
	}
	return ""
}
