package data

import (
	"context"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type sysHostRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysHostRepo(data *Data, logger log.Logger) biz.SysHostRepo {
	return &sysHostRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_host")),
	}
}

//列表
func (r *sysHostRepo) ListHost(ctx context.Context, in *model.SysHostReq) ([]*model.SysHost, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysHost{})
	if in != nil {
		name := strings.TrimSpace(in.Name)
		if name != "" {
			table = table.Where("name LIKE ?", "%"+name+"%")
		}
	}

	var list []*model.SysHost
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if in.PageNum <= 0 {
		in.PageNum = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	pageNum := int(in.PageNum)
	pageSize := int(in.PageSize)
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

//创建
func (r *sysHostRepo) CreateHost(ctx context.Context, in *model.SysHost) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysHostRepo) UpdateHost(ctx context.Context, id int64, in *model.SysHost) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysHost{}).
		Where("host_id = ?", id).Updates(in).Error
	return err
}

//获取
func (r *sysHostRepo) GetInfoById(ctx context.Context, id int64) (*model.SysHost, error) {
	var info *model.SysHost
	err := r.data.db.WithContext(ctx).Where("host_id =", id).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

//删除
func (r *sysHostRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	err := r.data.db.WithContext(ctx).
		Where("host_id IN (?)", ids).
		Unscoped().Delete(&model.SysHost{}).Error
	return err
}