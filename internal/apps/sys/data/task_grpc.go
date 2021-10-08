package data

import (
	"context"
	v1 "drpshop/api/task/v1"
	"drpshop/internal/apps/sys/biz"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type Client struct {
	conn       *grpc.ClientConn
	taskClient v1.TaskServiceClient
}

type taskGrpcPoolRepo struct {
	data *Data
	log  *log.Helper

	taskMap sync.Map
	// map key格式 ip:port
	conns map[string]*Client
	mu    sync.RWMutex
}

func NewTaskGrpcPoolRepo(data *Data, logger log.Logger) biz.TaskGrpcPoolRepo {
	return &taskGrpcPoolRepo{
		data:    data,
		conns:   make(map[string]*Client),
		taskMap: sync.Map{},
		log:     log.NewHelper(log.With(logger, "module", "sys/data/task_grpc")),
	}
}

func (p *taskGrpcPoolRepo) Get(addr string) (v1.TaskServiceClient, error) {
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
func (p *taskGrpcPoolRepo) Release(addr string) {
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
func (p *taskGrpcPoolRepo) factory(addr string) (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	client, ok := p.conns[addr]
	if ok {
		return client, nil
	}
	conn, taskClient, err := NewTaskServiceClient(addr)
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

func (p *taskGrpcPoolRepo) generateTaskUniqueKey(ip string, port int, id int64) string {
	return fmt.Sprintf("%s:%d:%d", ip, port, id)
}

func (p *taskGrpcPoolRepo) getWithEndpoint(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func (p *taskGrpcPoolRepo) Stop(ip string, port int, id int64) {
	key := p.generateTaskUniqueKey(ip, port, id)
	cancel, ok := p.taskMap.Load(key)
	if !ok {
		return
	}
	cancel.(context.CancelFunc)()
}

func (p *taskGrpcPoolRepo) Exec(ip string, port int, taskReq *v1.TaskReq) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			p.log.Error("panic#data/task_grpc.go:Exec#", err)
		}
	}()
	addr := fmt.Sprintf("%s:%d", ip, port)
	c, err := p.Get(addr)
	if err != nil {
		return "", err
	}
	if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
		taskReq.Timeout = 86400
	}
	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	taskUniqueKey := p.generateTaskUniqueKey(ip, port, taskReq.Id)
	p.taskMap.Store(taskUniqueKey, cancel)
	defer p.taskMap.Delete(taskUniqueKey)
	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		return parseGRPCError(err)
	}
	if resp.Error == "" {
		return resp.Output, nil
	}
	return resp.Output, errors.New(resp.Error)
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errors.New("无法连接远程服务器")
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时, 强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}
