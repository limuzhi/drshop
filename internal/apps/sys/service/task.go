package service

import (
	"context"
	v1 "drpshop/api/sys/v1"
)

//任务
func (s *SysService) TaskList(ctx context.Context, req *v1.TaskListReq) (*v1.TaskListRes, error) {
	return nil, nil
}

func (s *SysService) TaskDetail(ctx context.Context, req *v1.TaskDetailReq) (*v1.TaskData, error) {
	return nil, nil
}
func (s *SysService) TaskSave(ctx context.Context, re1 *v1.TaskSaveReq) (*v1.CommonRes, error) {
	return nil, nil
}
func (s *SysService) TaskChangeStatus(ctx context.Context, req *v1.TaskChangeStatusReq) (*v1.CommonRes, error) {
	return nil, nil
}
func (s *SysService) TaskRun(ctx context.Context, req *v1.TaskRunReq) (*v1.CommonRes, error) {
	return nil, nil
}
func (s *SysService) TaskDelete(ctx context.Context, req *v1.TaskDeleteReq) (*v1.CommonRes, error) {
	return nil, nil
}
