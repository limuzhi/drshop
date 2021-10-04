package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/internal/apps/sys/global"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
)

type SysDictDataRepo interface {
	//列表
	ListDictData(ctx context.Context, in *v1.DictDataListReq) ([]*model.SysDictData, int64, error)
	//添加
	CreateDictData(ctx context.Context, in *model.SysDictData) error
	//修改
	UpdateDictData(ctx context.Context, id int64, in *model.SysDictData) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysDictData, error)
	//获取列表 dictType
	ListByDictType(ctx context.Context, dictType string) ([]*model.SysDictData, error)
	//批量删除接口
	BatchDeleteByIds(ctx context.Context, ids []int64) error
}

type SysDictDataUsecase struct {
	repo SysDictDataRepo
	log  *log.Helper
}

func NewSysDictDataUsecase(repo SysDictDataRepo, logger log.Logger) *SysDictDataUsecase {
	return &SysDictDataUsecase{repo: repo, log: log.NewHelper(log.With(logger,
		"module", "sys/biz/sys_dict_data"))}
}

//列表
func (uc *SysDictDataUsecase) ListDictData(ctx context.Context, in *v1.DictDataListReq) ([]*model.SysDictData, int64, error) {
	return uc.repo.ListDictData(ctx, in)
}

//添加
func (uc *SysDictDataUsecase) CreateDictData(ctx context.Context, in *model.SysDictData) error {
	return uc.repo.CreateDictData(ctx, in)
}

//修改
func (uc *SysDictDataUsecase) UpdateDictData(ctx context.Context, id int64, in *model.SysDictData) error {
	return uc.repo.UpdateDictData(ctx, id, in)
}

//获取
func (uc *SysDictDataUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysDictData, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//获取列表 dictType
func (uc *SysDictDataUsecase) ListByDictType(ctx context.Context, dictType string) ([]*model.SysDictData, error) {
	return uc.repo.ListByDictType(ctx, dictType)
}

//批量删除接口
func (uc *SysDictDataUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

///====--------
func (uc *SysDictDataUsecase) DtoOut(data *model.SysDictData) *v1.DictDataListData {
	out := &v1.DictDataListData{
		DictCode:  data.DictCode,
		DictSort:  gconv.Int64(data.DictSort),
		DictLabel: data.DictLabel,
		DictValue: data.DictValue,
		DictType:  data.DictType,
		CssClass:  data.CssClass,
		ListClass: data.ListClass,
		IsDefault: gconv.Int64(data.IsDefault),
		Status:    gconv.Int64(data.Status),
		CreateBy:  data.CreateBy,
		UpdateBy:  data.UpdateBy,
		Remark:    data.Remark,
		CreatedAt: global.GetDateByUnix(data.CreatedAt),
		UpdatedAt: global.GetDateByUnix(data.UpdatedAt),
	}
	return out
}
