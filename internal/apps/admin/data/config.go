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

type configRepo struct {
	data *Data
	log  *log.Helper
}

func NewConfigRepo(data *Data, logger log.Logger) biz.ConfigRepo {
	return &configRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/config")),
	}
}

func (r *configRepo) GetConfigByKey(ctx context.Context, key string) (*v1.ConfigInfoByKeyRes, error) {
	info, err := r.data.sc.ConfigInfoByKey(ctx, &v1.ConfigInfoByKeyReq{ConfigKey: key})
	if err != nil {
		return nil, errors.Unauthorized("GetConfigByKeyError", err.Error())
	}
	return info, nil
}

//配置
func (r *configRepo) ConfigDetail(ctx context.Context, id int64) (*v1.ConfigDetailRes, error) {
	return r.data.sc.ConfigDetail(ctx, &v1.ConfigDetailReq{ConfigId: id})
}

func (r *configRepo) ConfigList(ctx context.Context, in *vo.ConfigListReq) (*v1.ConfigListRes, error) {
	return r.data.sc.ConfigList(ctx, &v1.ConfigListReq{
		ConfigKey:  strings.TrimSpace(in.ConfigKey),
		ConfigName: strings.TrimSpace(in.ConfigName),
		ConfigType: strings.TrimSpace(in.ConfigType),
		PageInfo: &common.PageReq{
			PageNum:  int64(in.GetPageIndex()),
			PageSize: int64(in.GetPageSize()),
		},
	})
}
func (r *configRepo) ConfigCreate(ctx context.Context, in *vo.ConfigCreateReq) error {
	_, err := r.data.sc.ConfigCreate(ctx, &v1.ConfigCreateUpdateReq{
		ConfigType:  gconv.Int64(in.ConfigType),
		ConfigName:  strings.TrimSpace(in.ConfigName),
		ConfigKey:   strings.TrimSpace(in.ConfigKey),
		ConfigValue: strings.TrimSpace(in.ConfigValue),
		Remark:      strings.TrimSpace(in.Remark),
		IsFrontend:  gconv.Int64(in.ConfigType),
		CreateBy:    in.CreateBy,
	})
	return err
}
func (r *configRepo) ConfigUpdate(ctx context.Context, in *vo.ConfigUpdateReq) error {
	_, err := r.data.sc.ConfigUpdate(ctx, &v1.ConfigCreateUpdateReq{
		ConfigId:    in.ConfigId,
		ConfigType:  gconv.Int64(in.ConfigType),
		ConfigName:  strings.TrimSpace(in.ConfigName),
		ConfigKey:   strings.TrimSpace(in.ConfigKey),
		ConfigValue: strings.TrimSpace(in.ConfigValue),
		Remark:      strings.TrimSpace(in.Remark),
		IsFrontend:  gconv.Int64(in.IsFrontend),
		UpdateBy:    in.UpdateBy,
	})
	return err
}
func (r *configRepo) ConfigDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.ConfigDelete(ctx, &v1.ConfigDeleteReq{ConfigIds: ids})
	return err
}
