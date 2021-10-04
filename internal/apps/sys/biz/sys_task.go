package biz

import (
	"context"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysTaskRepo interface {
	//列表
	ListTask(ctx context.Context, in *model.SysTaskReq) ([]*model.SysTask, int64, error)
	//详情
	GetInfoById(ctx context.Context, id int64) (*model.SysTask, error)
	//创建
	CreateTask(ctx context.Context, in *model.SysTask) error
	//更新
	UpdateTask(ctx context.Context, id int64, in *model.SysTask) error
	//删除
	DeleteById(ctx context.Context, id int64) error
	//获取所有激活任务
	ActiveList(ctx context.Context, pageNum int, pageSize int) ([]*model.SysTask, error)
	//获取某个主机下的所有激活任务
	ActiveListByHostId(ctx context.Context, hostId int64) ([]*model.SysTask, error)
	//获取依赖任务列表
	GetDependencyTaskList(ctx context.Context, taskIds string) ([]*model.SysTask, error)
	//更新任务状态
	ChangeStatus(ctx context.Context, taskId int64, status model.Status) error
}

type SysTaskUsecase struct {
	repo SysTaskRepo
	log  *log.Helper
}

func NewSysTaskUsecase(repo SysTaskRepo, logger log.Logger) *SysTaskUsecase {
	return &SysTaskUsecase{repo: repo, log: log.NewHelper(
		log.With(logger, "module", "sys/biz/sys_task"))}
}

//列表
func (uc *SysTaskUsecase) ListTask(ctx context.Context, in *model.SysTaskReq) ([]*model.SysTask, int64, error) {
	return uc.repo.ListTask(ctx, in)
}

//详情
func (uc *SysTaskUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysTask, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//创建
func (uc *SysTaskUsecase) CreateTask(ctx context.Context, in *model.SysTask) error {
	return uc.repo.CreateTask(ctx, in)
}

//更新
func (uc *SysTaskUsecase) UpdateTask(ctx context.Context, id int64, in *model.SysTask) error {
	return uc.repo.UpdateTask(ctx, id, in)
}

//删除
func (uc *SysTaskUsecase) DeleteById(ctx context.Context, id int64) error {
	return uc.repo.DeleteById(ctx, id)
}

//获取所有激活任务
func (uc *SysTaskUsecase) ActiveList(ctx context.Context, pageNum int, pageSize int) ([]*model.SysTask, error) {
	return uc.repo.ActiveList(ctx, pageNum, pageSize)
}

//获取某个主机下的所有激活任务
func (uc *SysTaskUsecase) ActiveListByHostId(ctx context.Context, hostId int64) ([]*model.SysTask, error) {
	return uc.repo.ActiveListByHostId(ctx, hostId)
}

//获取依赖任务列表
func (uc *SysTaskUsecase) GetDependencyTaskList(ctx context.Context, taskIds string) ([]*model.SysTask, error) {
	return uc.repo.GetDependencyTaskList(ctx, taskIds)
}

//更新任务状态
func (uc *SysTaskUsecase) ChangeStatus(ctx context.Context, taskId int64, status model.Status) error {
	return uc.repo.ChangeStatus(ctx, taskId, status)
}

//=====
