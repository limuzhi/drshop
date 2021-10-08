package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/model"
	"drpshop/internal/apps/sys/global"
	"github.com/go-kratos/kratos/v2/log"
)

type SysConfigRepo interface {
	//列表
	ListConfig(ctx context.Context, in *v1.ConfigListReq) ([]*model.SysConfig, int64, error)
	//删除
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//创建
	CreateConfig(ctx context.Context, in *model.SysConfig) error
	//修改
	UpdateConfig(ctx context.Context, id int64, in *model.SysConfig) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysConfig, error)
	//
	GetInfoByConfigKey(ctx context.Context, configKey string) (*model.SysConfig, error)
}

type SysConfigUsecase struct {
	repo SysConfigRepo
	log  *log.Helper
}

func NewSysConfigUsecase(repo SysConfigRepo, logger log.Logger) *SysConfigUsecase {
	return &SysConfigUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "sys/biz/sys_config"))}
}

//列表
func (e *SysConfigUsecase) ListConfig(ctx context.Context, in *v1.ConfigListReq) ([]*model.SysConfig,
	int64, error) {
	if in.PageInfo.PageNum <= 0 {
		in.PageInfo.PageNum = 1
	}
	if in.PageInfo.PageSize <= 0 {
		in.PageInfo.PageSize = 10
	}
	return e.repo.ListConfig(ctx, in)
}

//删除
func (e *SysConfigUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return e.repo.BatchDeleteByIds(ctx, ids)
}

//创建
func (e *SysConfigUsecase) CreateConfig(ctx context.Context, in *model.SysConfig) error {
	return e.repo.CreateConfig(ctx, in)
}

//修改
func (e *SysConfigUsecase) UpdateConfig(ctx context.Context, id int64, in *model.SysConfig) error {
	return e.repo.UpdateConfig(ctx, id, in)
}

//获取
func (e *SysConfigUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysConfig, error) {
	return e.repo.GetInfoById(ctx, id)
}

func (e *SysConfigUsecase) GetInfoByConfigKey(ctx context.Context, configKey string) (*model.SysConfig, error) {
	return e.repo.GetInfoByConfigKey(ctx, configKey)
}

///=============

func (e *SysConfigUsecase) DtoOut(data *model.SysConfig) *v1.ConfigDetailRes {
	info := &v1.ConfigDetailRes{
		ConfigId:    data.ConfigId,
		ConfigName:  data.ConfigName,
		ConfigKey:   data.ConfigKey,
		ConfigValue: data.ConfigValue,
		ConfigType:  int64(data.ConfigType),
		IsFrontend:  int64(data.IsFrontend),
		Remark:      data.Remark,
		CreateBy:    data.CreateBy,
		UpdateBy:    data.UpdateBy,
		CreatedAt:   global.GetDateByUnix(data.CreatedAt),
		UpdatedAt:   global.GetDateByUnix(data.UpdatedAt),
	}
	return info
}
