package biz

import (
	"context"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysUserOnlineRepo interface {
	//列表
	ListUserOnline(ctx context.Context, in *model.SysUserOnlineSearchReq) ([]*model.SysUserOnline, int64, error)
	//创建
	CreateUserOnline(ctx context.Context, in *model.SysUserOnline) error
	//根据ID获取
	GetListByIds(ctx context.Context, ids []int64) ([]*model.SysUserOnline, error)
	//保存
	SaveUserOnline(ctx context.Context, data *model.SysUserOnline) error
	//删除用户在线状态操作
	DeleteByToken(ctx context.Context, token string) error
	//删除用户在线状态操作
	BatchDeleteByIds(ctx context.Context, ids []int64) ([]string, error)
}

type SysUserOnlineUsecase struct {
	repo SysUserOnlineRepo
	log  *log.Helper
}

func NewSysUserOnlineUsecase(repo SysUserOnlineRepo, logger log.Logger) *SysUserOnlineUsecase {
	return &SysUserOnlineUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/sys_user_online")),
	}
}

//列表
func (uc *SysUserOnlineUsecase) ListUserOnline(ctx context.Context, in *model.SysUserOnlineSearchReq) ([]*model.SysUserOnline, int64, error) {
	return uc.repo.ListUserOnline(ctx, in)
}

//创建
func (uc *SysUserOnlineUsecase) CreateUserOnline(ctx context.Context, in *model.SysUserOnline) error {
	return uc.repo.CreateUserOnline(ctx, in)
}

//根据ID获取
func (uc *SysUserOnlineUsecase) GetListByIds(ctx context.Context, ids []int64) ([]*model.SysUserOnline, error) {
	return uc.repo.GetListByIds(ctx, ids)
}

//保存
func (uc *SysUserOnlineUsecase) SaveUserOnline(ctx context.Context, data *model.SysUserOnline) error {
	return uc.repo.SaveUserOnline(ctx, data)
}

//删除用户在线状态操作
func (uc *SysUserOnlineUsecase) DeleteByToken(ctx context.Context, token string) error {
	return uc.repo.DeleteByToken(ctx, token)
}

//删除用户在线状态操作
func (uc *SysUserOnlineUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) ([]string, error) {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}
