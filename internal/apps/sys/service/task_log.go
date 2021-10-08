package service

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
)

//任务日志
func (s *SysService) TaskLogList(ctx context.Context, req *v1.TaskLogListReq) (*v1.TaskLogListRes, error) {
	return nil, nil
}
func (s *SysService) TaskLogClear(ctx context.Context, req *common.NullReq) (*common.NullRes, error) {
	return nil, nil
}
func (s *SysService) TaskLogStop(ctx context.Context, req *v1.TaskLogStopReq) (*common.NullRes, error) {
	return nil, nil
}
