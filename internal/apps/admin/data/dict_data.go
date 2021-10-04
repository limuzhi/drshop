package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type dictDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewDictDataRepo(data *Data, logger log.Logger) biz.DictDataRepo {
	return &dictDataRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/dict_data")),
	}
}

func (r *dictDataRepo) GetDictDataListSelect(ctx context.Context, dictType string) (*v1.DictDataListRes, error) {
	res, err := r.data.sc.DictDataListByDictType(ctx, &v1.DictDataListByDictTypeReq{DictType: dictType})
	if err != nil {
		return nil, errors.Unauthorized("DictListByDictTypeError", err.Error())
	}
	return res, nil
}

func (r *dictDataRepo) DictDataList(ctx context.Context, in *vo.DictDataListReq) (*v1.DictDataListRes, error) {
	return r.data.sc.DictDataList(ctx, &v1.DictDataListReq{
		DictType:  strings.TrimSpace(in.DictType),
		DictLabel: strings.TrimSpace(in.DictLabel),
		DictValue: strings.TrimSpace(in.DictValue),
		Status:    gconv.Int64(in.Status),
		PageInfo: &common.PageReq{
			PageNum:  int64(in.GetPageIndex()),
			PageSize: int64(in.GetPageSize()),
		},
	})
}

func (r *dictDataRepo) DictDataDetail(ctx context.Context, id int64) (*v1.DictDataListData, error) {
	return r.data.sc.DictDataInfoByDictCode(ctx, &v1.DictDataInfoByDictCodeReq{DictCode: id})
}

func (r *dictDataRepo) DictDataCreate(ctx context.Context, in *vo.DictDataCreateReq) error {
	_, err := r.data.sc.DictDataAdd(ctx, &v1.DictDataAddReq{
		DictSort:  gconv.Int64(in.DictSort),
		DictLabel: in.DictLabel,
		DictValue: in.DictValue,
		DictType:  in.DictType,
		CssClass:  in.CssClass,
		ListClass: in.ListClass,
		IsDefault: gconv.Int64(in.IsDefault),
		Status:    gconv.Int64(in.Status),
		CreateBy:  in.CreateBy,
		Remark:    in.Remark,
	})
	return err
}

func (r *dictDataRepo) DictDataUpdate(ctx context.Context, in *vo.DictDataUpdateReq) error {
	_, err := r.data.sc.DictDataUpdate(ctx, &v1.DictDataUpdateReq{
		DictCode:  in.DictCode,
		DictSort:  gconv.Int64(in.DictSort),
		DictLabel: in.DictLabel,
		DictValue: in.DictValue,
		DictType:  in.DictType,
		CssClass:  in.CssClass,
		ListClass: in.ListClass,
		IsDefault: gconv.Int64(in.IsDefault),
		Status:    gconv.Int64(in.Status),
		UpdateBy:  in.UpdateBy,
		Remark:    in.Remark,
	})
	return err
}

func (r *dictDataRepo) DictDataDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.DictDataDelete(ctx, &v1.DictDataDeleteReq{DictCodes: ids})
	return err
}
