package data

import (
	"context"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/model"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type sysTaskRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysTaskRepo(data *Data, logger log.Logger) biz.SysTaskRepo {
	return &sysTaskRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/task")),
	}
}

//列表
func (r *sysTaskRepo) ListTask(ctx context.Context, in *model.SysTaskReq) ([]*model.SysTask, int64, error) {
	table := r.data.db.WithContext(ctx)
	table = table.Model(&model.SysTask{})
	if in.TaskId > 0 {
		table = table.Where("task_id = ?", in.TaskId)
	}
	if in.Name != "" {
		table = table.Where("name LIKE ?", "%"+in.Name+"%")
	}
	if in.Protocol > 0 {
		table = table.Where("protocol = ?", in.Protocol)
	}
	if in.Status != 0 {
		table = table.Where("status = ?", in.Status)
	}
	if in.Tag != "" {
		table = table.Where("tag = ?", in.Tag)
	}

	var list []*model.SysTask
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
	err = table.Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Preload("Hosts").Find(&list).Error
	return list, total, err
}

//详情
func (r *sysTaskRepo) GetInfoById(ctx context.Context, id int64) (*model.SysTask, error) {
	var info *model.SysTask
	err := r.data.db.Where("task_id = ?", id).Preload("hosts").First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, err
}

//创建
func (r *sysTaskRepo) CreateTask(ctx context.Context, in *model.SysTask) error {
	return r.data.db.WithContext(ctx).Create(in).Error
}

//更新
func (r *sysTaskRepo) UpdateTask(ctx context.Context, id int64, in *model.SysTask) error {
	tx := r.data.db.WithContext(ctx)
	err := tx.Model(in).Where("task_id = ?", id).Updates(in).Error
	if err != nil {
		return err
	}
	err = tx.Model(in).Association("Hosts").Replace(in.Hosts)
	if err != nil {
		return err
	}
	return nil
}

//删除
func (r *sysTaskRepo) DeleteById(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Where("task_id = ?", id).
		Delete(&model.SysTask{}).Error
}

//获取所有激活任务
func (r *sysTaskRepo) ActiveList(ctx context.Context, pageNum int, pageSize int) ([]*model.SysTask, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 2000
	}
	var list []*model.SysTask
	table := r.data.db.WithContext(ctx)
	err := table.Where("status = ? AND level = ?", model.Enabled, model.TaskLevelParent).
		Preload("Hosts").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

//获取某个主机下的所有激活任务
func (r *sysTaskRepo) ActiveListByHostId(ctx context.Context, hostId int64) ([]*model.SysTask, error) {
	var taskHostList []*model.SysTaskHost
	tx := r.data.db.WithContext(ctx)
	err := tx.Model(&model.SysTaskHost{}).Where("host_id = ?", hostId).Find(&taskHostList).Error
	if err != nil {
		return nil, err
	}
	taskIds := make([]int64, len(taskHostList))
	for i, value := range taskHostList {
		taskIds[i] = value.TaskId
	}
	if len(taskIds) == 0 {
		return nil, errors.New("未找到任务ID")
	}
	var list []*model.SysTask
	table := tx.Model(&model.SysTask{}).Where("status = ? AND level =?", model.Enabled, model.TaskLevelParent)
	table = table.Where("task_id IN(?)", taskIds).Preload("Hosts").Find(&list)
	if err := table.Error; err != nil {
		return nil, err
	}
	return list, nil
}

//获取依赖任务列表
func (r *sysTaskRepo) GetDependencyTaskList(ctx context.Context, inIds string) ([]*model.SysTask, error) {
	if inIds == "" {
		return nil, errors.New("taskIds不能为空")
	}
	idList := strings.Split(inIds, ",")
	taskIds := make([]interface{}, len(idList))
	for i, v := range idList {
		taskIds[i] = v
	}
	var list []*model.SysTask
	err := r.data.db.WithContext(ctx).
		Where("level = ?", model.TaskLevelChild).
		Preload("Hosts").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

//更新任务状态
func (r *sysTaskRepo) ChangeStatus(ctx context.Context, taskId int64, status model.Status) error {
	err := r.data.db.WithContext(ctx).Model(&model.SysTask{}).Where("task_id = ?", taskId).
		Updates(map[string]interface{}{"status": status}).Error
	return err
}
