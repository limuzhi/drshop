package data

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/model"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
)

type sysOperLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysOperLogRepo(data *Data, logger log.Logger) biz.SysOperLogRepo {
	return &sysOperLogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_operlog")),
	}
}

//列表
func (r *sysOperLogRepo) ListOperLog(ctx context.Context, in *v1.OperLogListReq) ([]*model.SysOperLog, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysOperLog{})
	if in.OperName != "" {
		table = table.Where("oper_name = ?", in.OperName)
	}
	if in.Title != "" {
		table = table.Where("title LIKE ?", fmt.Sprintf("%%%s%%", in.Title))
	}
	if in.Status != "" {
		table = table.Where("status = ?", gconv.Int(in.Status))
	}
	if in.PageInfo.BeginTime != 0 {
		table = table.Where("oper_time >= ?", in.PageInfo.BeginTime)
	}
	if in.PageInfo.EndTime != 0 {
		table = table.Where("oper_time <= ?", in.PageInfo.EndTime)
	}
	var list []*model.SysOperLog
	var total int64
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	pageNum := int(in.PageInfo.GetPageNum())
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := int(in.PageInfo.GetPageSize())
	if pageSize <= 0 {
		pageSize = 10
	}
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("oper_time DESC").Find(&list).Error
	return list, total, err
}

//创建
func (r *sysOperLogRepo) BatchInsertOperLog(ctx context.Context, ins []*model.SysOperLog) error {
	return r.data.db.WithContext(ctx).CreateInBatches(ins, len(ins)).Error
}

//根据ID获取
func (r *sysOperLogRepo) GetInfoById(ctx context.Context, id int64) (*model.SysOperLog, error) {
	var info *model.SysOperLog
	err := r.data.db.WithContext(ctx).Where("oper_id = ?", id).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

//批量删除接口
func (r *sysOperLogRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	err := r.data.db.WithContext(ctx).Where("oper_id IN (?)", ids).
		Unscoped().Delete(&model.SysOperLog{}).Error
	return err
}

func (r *sysOperLogRepo) ClearOperLog(ctx context.Context) error {
	modelLog := &model.SysOperLog{}
	return r.data.db.WithContext(ctx).Exec("truncate " + modelLog.TableName()).Error
}

