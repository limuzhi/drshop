package taskrpc

import (
	v1 "drpshop/api/tasknode/v1"
	"drpshop/internal/apps/sys/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

const (
	backOffMaxDelay = 3 * time.Second
	dialTimeout     = 2 * time.Second
)

var (
	Pool = &GRPCPool{
		conns: make(map[string]*Client),
	}

	keepAliveParams = keepalive.ClientParameters{
		Time:                20 * time.Second,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}
)

type Client struct {
	conn       *grpc.ClientConn
	taskClient v1.TaskClient
}

type GRPCPool struct {
	// map key格式 ip:port
	conns map[string]*Client
	mu    sync.RWMutex
}

func (p *GRPCPool) Get(addr string) (v1.TaskClient, error) {
	p.mu.RLock()
	client, ok := p.conns[addr]
	p.mu.RUnlock()
	if ok {
		return client.taskClient, nil
	}

	client, err := p.factory(addr)
	if err != nil {
		return nil, err
	}

	return client.taskClient, nil
}

// 释放连接
func (p *GRPCPool) Release(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	client, ok := p.conns[addr]
	if !ok {
		return
	}
	delete(p.conns, addr)
	client.conn.Close()
}

// 创建连接
func (p *GRPCPool) factory(addr string) (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	client, ok := p.conns[addr]
	if ok {
		return client, nil
	}
	discovery := data.NewDiscovery(addr, "http")

	conn, taskClient, err := data.NewTaskServiceClient(discovery)
	if err != nil {
		return nil, err
	}
	client = &Client{
		conn:       conn,
		taskClient: taskClient,
	}

	p.conns[addr] = client

	return client, nil
}
