package model

import (
	"gorm.io/plugin/soft_delete"
)

type BaseModelTime struct {
	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime;comment:'创建时间'" json:"createdAt"`
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime;comment:'更新时间'" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:flag;column:deleted_at;comment:'删除时间'" json:"-"`
}

type ControlBy struct {
	CreateBy int64 `gorm:"column:create_by;comment:创建者" json:"createBy"`
	UpdateBy int64 `gorm:"column:update_by;comment:更新者" json:"updateBy"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy int64) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int64) {
	e.UpdateBy = updateBy
}

type Status int

const (
	Disabled Status = 1 // 禁用
	Failure  Status = 1 // 失败
	Enabled  Status = 2 // 启用
	Running  Status = 2 // 运行中
	Finish   Status = 3 // 完成
	Cancel   Status = 4 // 取消
)
