package data

import (
	"context"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type sysJobRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysJobRepo(data *Data, logger log.Logger) biz.SysJobRepo {
	return &sysJobRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_job")),
	}
}

//列表
func (r *sysJobRepo) ListJob(ctx context.Context, in *model.SysJobSearchReq) ([]*model.SysJob, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysApis{})
	if in != nil {
		jobName := strings.TrimSpace(in.JobName)
		if jobName != "" {
			table = table.Where("job_name LIKE ?", fmt.Sprintf("%%%s%%", jobName))
		}
		jobGroup := strings.TrimSpace(in.JobGroup)
		if jobGroup != "" {
			table = table.Where("job_group = ?", jobGroup)
		}
		if in.Status != 0 {
			table = table.Where("status = ?", in.Status)
		}
	}

	var list []*model.SysJob
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
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

//添加
func (r *sysJobRepo) CreateJob(ctx context.Context, in *model.SysJob) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//修改
func (r *sysJobRepo) UpdateJob(ctx context.Context, id int64, in *model.SysJob) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysJob{}).
		Where("job_id = ?", id).Updates(in).Error
	return err
}

//获取
func (r *sysJobRepo) GetInfoById(ctx context.Context, id int64) (*model.SysJob, error) {
	var info *model.SysJob
	err := r.data.db.WithContext(ctx).Where("job_id = ?", id).First(&info).Error
	return info, err
}

//批量删除
func (r *sysJobRepo) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	err := r.data.db.WithContext(ctx).
		Where("job_id IN (?)", ids).Unscoped().Delete(&model.SysJob{}).Error
	return err
}

//获取已开启执行的任务
func (r *sysJobRepo) GetAllJobs(ctx context.Context) ([]*model.SysJob, int64, error) {
	var list []*model.SysJob
	var total int64
	table := r.data.db.WithContext(ctx)
	table = table.Where("status = ?", 2)
	err := table.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = table.Order("created_at DESC").Find(&list).Error
	return list, total, err
}
