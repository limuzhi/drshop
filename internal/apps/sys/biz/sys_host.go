package biz

import (
	"context"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysHostRepo interface {
	//列表
	ListHost(ctx context.Context, in *model.SysHostReq) ([]*model.SysHost, int64, error)
	//创建
	CreateHost(ctx context.Context, in *model.SysHost) error
	//更新
	UpdateHost(ctx context.Context, id int64, in *model.SysHost) error
	//删除
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysHost, error)
}

type SysHostUsecase struct {
	repo SysHostRepo
	log  *log.Helper
}

func NewSysHostUsecase(repo SysHostRepo, logger log.Logger) *SysHostUsecase {
	return &SysHostUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/sys_host")),
	}
}

//列表
func (uc *SysHostUsecase) ListHost(ctx context.Context, in *model.SysHostReq) ([]*model.SysHost, int64, error) {
	return uc.repo.ListHost(ctx, in)
}

//创建
func (uc *SysHostUsecase) CreateHost(ctx context.Context, in *model.SysHost) error {
	return uc.repo.CreateHost(ctx, in)
}

//更新
func (uc *SysHostUsecase) UpdateHost(ctx context.Context, id int64, in *model.SysHost) error {
	return uc.repo.UpdateHost(ctx, id, in)
}

//删除
func (uc *SysHostUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//获取
func (uc *SysHostUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysHost, error) {
	return uc.repo.GetInfoById(ctx, id)
}
