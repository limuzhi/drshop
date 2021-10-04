package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type dictTypeRepo struct {
	data *Data
	log  *log.Helper
}

func NewDictTypeRepo(data *Data, logger log.Logger) biz.DictTypeRepo {
	return &dictTypeRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/dict_type")),
	}
}

//获取字典分类列表
func (r *dictTypeRepo) DictTypeList(ctx context.Context, in *vo.DictTypeListReq) (*v1.DictTypeListRes, error) {
	return r.data.sc.DictTypeList(ctx, &v1.DictTypeListReq{
		DictType: strings.TrimSpace(in.DictType),
		DictName: strings.TrimSpace(in.DictName),
		Status:   gconv.Int64(in.Status),
		PageInfo: &common.PageReq{
			PageNum:  int64(in.GetPageIndex()),
			PageSize: int64(in.GetPageSize()),
		},
	})
}

func (r *dictTypeRepo) DictTypeSelect(ctx context.Context) (*v1.DictTypeListRes, error) {
	return r.data.sc.DictTypeOptionSelect(ctx, &v1.CommonReq{})
}

func (r *dictTypeRepo) DictTypeDetail(ctx context.Context, id int64) (*v1.DictTypeListData, error) {
	return r.data.sc.DictTypeInfo(ctx, &v1.DictTypeInfoReq{DictId: id})
}

//添加字典分类数据
func (r *dictTypeRepo) DictTypeCreate(ctx context.Context, in *vo.DictTypeCreateReq) error {
	_, err := r.data.sc.DictTypeAdd(ctx, &v1.DictTypeAddReq{
		DictName: strings.TrimSpace(in.DictName),
		DictType: strings.TrimSpace(in.DictType),
		Remark:   strings.TrimSpace(in.Remark),
		Status:   gconv.Int64(in.Status),
		CreateBy: in.CreateBy,
	})
	return err
}

func (r *dictTypeRepo) DictTypeUpdate(ctx context.Context, in *vo.DictTypeUpdateReq) error {
	_, err := r.data.sc.DictTypeUpdate(ctx, &v1.DictTypeUpdateReq{
		DictId:   in.DictId,
		DictName: strings.TrimSpace(in.DictName),
		DictType: strings.TrimSpace(in.DictType),
		Remark:   strings.TrimSpace(in.Remark),
		Status:   gconv.Int64(in.Status),
		UpdateBy: in.UpdateBy,
	})
	return err
}

//删除字典分类数据
func (r *dictTypeRepo) DictTypeDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.DictTypeDelete(ctx, &v1.DictTypeDeleteReq{DictIds: ids})
	return err
}
