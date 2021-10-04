package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/internal/apps/sys/global"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
)

type SysDictTypeRepo interface {
	//列表
	ListDictType(ctx context.Context, in *v1.DictTypeListReq) ([]*model.SysDictType, int64, error)
	//添加
	CreateDictType(ctx context.Context, in *model.SysDictType) error
	//修改
	UpdateDictType(ctx context.Context, id int64, in *model.SysDictType) error
	//获取
	GetInfoById(ctx context.Context, id int64) (*model.SysDictType, error)
	//检查类型是否已经存在
	CheckByDictType(ctx context.Context, dictType string, ids ...int64) bool
	//批量删除接口
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//获取所有正常状态下的字典类型
	GetAllDictType(ctx context.Context) ([]*model.SysDictType, error)
}

type SysDictTypeUsecase struct {
	repo SysDictTypeRepo
	log  *log.Helper
}

func NewSysDictTypeUsecase(repo SysDictTypeRepo, logger log.Logger) *SysDictTypeUsecase {
	return &SysDictTypeUsecase{repo: repo, log: log.NewHelper(log.With(logger,
		"module", "sys/biz/sys_dict_type"))}
}

//列表
func (uc *SysDictTypeUsecase) ListDictType(ctx context.Context, in *v1.DictTypeListReq) ([]*model.SysDictType, int64, error) {
	return uc.repo.ListDictType(ctx, in)
}

//添加
func (uc *SysDictTypeUsecase) CreateDictType(ctx context.Context, in *model.SysDictType) error {
	return uc.repo.CreateDictType(ctx, in)
}

//修改
func (uc *SysDictTypeUsecase) UpdateDictType(ctx context.Context, id int64, in *model.SysDictType) error {
	return uc.repo.UpdateDictType(ctx, id, in)
}

//获取
func (uc *SysDictTypeUsecase) GetInfoById(ctx context.Context, id int64) (*model.SysDictType, error) {
	return uc.repo.GetInfoById(ctx, id)
}

//检查类型是否已经存在
func (uc *SysDictTypeUsecase) CheckByDictType(ctx context.Context, dictType string, ids ...int64) bool {
	return uc.repo.CheckByDictType(ctx, dictType, ids...)
}

//批量删除接口
func (uc *SysDictTypeUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//获取所有正常状态下的字典类型
func (uc *SysDictTypeUsecase) GetAllDictType(ctx context.Context) ([]*model.SysDictType, error) {
	return uc.repo.GetAllDictType(ctx)
}

///==========
func (uc *SysDictTypeUsecase) DtoOut(data *model.SysDictType) *v1.DictTypeListData {
	out := &v1.DictTypeListData{
		DictId:    data.DictId,
		DictName:  data.DictName,
		DictType:  data.DictType,
		Status:    gconv.Int64(data.Status),
		Remark:    data.Remark,
		CreateBy:  data.CreateBy,
		UpdateBy:  data.UpdateBy,
		CreatedAt: global.GetDateByUnix(data.CreatedAt),
		UpdatedAt: global.GetDateByUnix(data.UpdatedAt),
	}
	return out
}
