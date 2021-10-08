package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/global"
	"drpshop/internal/apps/sys/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysLoginLogRepo interface {
	//列表
	ListLoginLog(ctx context.Context, in *v1.LoginlogListReq) ([]*model.SysLoginLog, int64, error)
	//创建
	CreateLoginLog(ctx context.Context, ins []*model.SysLoginLog) error
	//删除
	BatchDeleteByIds(ctx context.Context, ins []int64) error
	//获取
	GetInfo(ctx context.Context, id int64) (*model.SysLoginLog, error)
	//清空日志
	ClearLoginLog(ctx context.Context) error
}

type SysLoginLogUsecase struct {
	repo SysLoginLogRepo
	log  *log.Helper
}

func NewSysLoginLogUsecase(repo SysLoginLogRepo, logger log.Logger) *SysLoginLogUsecase {
	return &SysLoginLogUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/sys_user"))}
}

//列表
func (uc *SysLoginLogUsecase) ListLoginLog(ctx context.Context, in *v1.LoginlogListReq) ([]*model.SysLoginLog, int64, error) {
	return uc.repo.ListLoginLog(ctx, in)
}

//创建
func (uc *SysLoginLogUsecase) CreateLoginLog(ctx context.Context, ins []*model.SysLoginLog) error {
	return uc.repo.CreateLoginLog(ctx, ins)
}

//删除
func (uc *SysLoginLogUsecase) BatchDeleteByIds(ctx context.Context, ins []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ins)
}

//删除
func (uc *SysLoginLogUsecase) ClearLoginLog(ctx context.Context) error {
	return uc.repo.ClearLoginLog(ctx)
}

//获取
func (uc *SysLoginLogUsecase) GetInfo(ctx context.Context, id int64) (*model.SysLoginLog, error) {
	return uc.repo.GetInfo(ctx, id)
}

//=============

func (uc *SysLoginLogUsecase) SaveLoginlog(ctx context.Context, data *model.SysLoginLog) error {
	ins := make([]*model.SysLoginLog, 0)
	ins = append(ins, data)
	return uc.repo.CreateLoginLog(ctx, ins)
}

func (uc *SysLoginLogUsecase) DtoOut(data []*model.SysLoginLog) []*v1.LoginlogListData {
	out := make([]*v1.LoginlogListData, len(data))
	for k, v := range data {
		info := &v1.LoginlogListData{
			LoginName:     v.LoginName,
			LoginUid:      v.LoginUid,
			Ipaddr:        v.Ipaddr,
			LoginLocation: v.LoginLocation,
			Browser:       v.Browser,
			Os:            v.Os,
			Status:        v.Status,
			Msg:           v.Msg,
			LoginTime:     global.GetDateByUnix(v.LoginTime),
			LoginId:       v.LoginId,
		}
		out[k] = info
	}
	return out
}
