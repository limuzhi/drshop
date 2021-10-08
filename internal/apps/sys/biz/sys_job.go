package biz

import (
	"context"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysJobRepo interface {
	//列表
	ListJob(ctx context.Context, in *model.SysJobSearchReq) ([]*model.SysJob, int64, error)
	//获取已开启执行的任务
	GetAllJobs(ctx context.Context) ([]*model.SysJob, int64, error)
	//添加
	CreateJob(ctx context.Context, in *model.SysJob) error
	//修改
	UpdateJob(ctx context.Context, id int64, in *model.SysJob) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysJob, error)
	//批量删除
	BatchDeleteByIds(ctx context.Context, ids []int64) error
}

type SysJobUsecase struct {
	repo SysJobRepo
	log  *log.Helper
}

func NewSysJobUsecase(repo SysJobRepo, logger log.Logger) *SysJobUsecase {
	return &SysJobUsecase{repo: repo, log: log.NewHelper(
		log.With(logger, "module", "sys/biz/sys_job"))}
}

//列表
func (uc *SysJobUsecase) ListJob(ctx context.Context, in *model.SysJobSearchReq) ([]*model.SysJob, int64, error) {
	return uc.repo.ListJob(ctx, in)
}

//获取已开启执行的任务
func (uc *SysJobUsecase) GetAllJobs(ctx context.Context) ([]*model.SysJob, int64, error) {
	return uc.repo.GetAllJobs(ctx)
}

//添加
func (uc *SysJobUsecase) CreateJob(ctx context.Context, in *model.SysJob) error {
	return uc.repo.CreateJob(ctx, in)
}

//修改
func (uc *SysJobUsecase) UpdateJob(ctx context.Context, id int64, in *model.SysJob) error {
	return uc.repo.UpdateJob(ctx, id, in)
}

//获取
func (uc *SysJobUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysJob, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//批量删除
func (uc *SysJobUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//============================================
//TODO
//启动任务
func (uc *SysJobUsecase) JobStart(ctx context.Context, job *model.SysJob) error {
	return nil
}

//停止任务
func (uc *SysJobUsecase) JobStop(ctx context.Context, job *model.SysJob) error {
	return nil
}

//执行任务
func (uc *SysJobUsecase) JobRun(ctx context.Context, job *model.SysJob) error {
	return nil
}

//删除任务
func (uc *SysJobUsecase) DeleteJobByIds(ctx context.Context, ids []int64) error {
	return nil
}
