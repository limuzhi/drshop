package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
)

//任务节点
func (s *SysService) HostList(ctx context.Context, req *v1.HostListReq) (*v1.HostListRes, error) {
	return nil, nil
}
func (s *SysService) HostDetail(ctx context.Context, req *v1.HostDetailReq) (*v1.HostData, error) {
	return nil, nil
}

func (s *SysService) HostSave(ctx context.Context, req *v1.HostSaveReq) (*v1.CommonRes, error) {
	return nil, nil
}

func (s *SysService) HostDelete(ctx context.Context, req *v1.HostDeleteReq) (*v1.CommonRes, error) {
	return nil, nil
}
func (s *SysService) HostPing(ctx context.Context, req *v1.HostPingReq) (*v1.CommonRes, error) {
	return nil, nil
}
