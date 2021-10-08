package model

type SysTaskHost struct {
	TaskId                     int64                `gorm:"column:task_id;" json:"taskId"`
	HostId                     int64                `gorm:"column:host_id;" json:"hostId"`
}

func (e *SysTaskHost) TableName() string {
	return "sys_task_host"
}
