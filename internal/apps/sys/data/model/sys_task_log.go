package model

// 任务执行日志
type SysTaskLog struct {
	Id         int64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskId     int64         `gorm:"column:task_id;type:bigint;comment:'任务id'" json:"taskId"`
	Name       string        `gorm:"column:name;type:varchar(32);comment:'任务名称'" json:"name"`
	Spec       string        `gorm:"column:spec;type:varchar(64);comment:'crontab'" json:"spec"`
	Protocol   TaskProtocol  `gorm:"column:protocol;type:tinyint;comment:'协议 1:http 2:系统命令-RPC方式执行命令'" json:"protocol"`
	Command    string        `gorm:"column:command;type:varchar(255);comment:'URL地址或shell命令'" json:"command"`
	Timeout    int           `gorm:"column:timeout;type:mediumint;default:0;comment:'任务执行超时时间(单位秒),0不限制'" json:"timeout"`
	RetryTimes int8          `gorm:"column:retry_times;type:tinyint;default:0;comment:'重试次数" json:"retryTimes"`
	Hostname   string        `gorm:"column:hostname;type:varchar(128);default:'';comment:'RPC主机名，逗号分隔'" json:"hostname"`
	StartTime  int64         `gorm:"column:start_time;type:int;default:0;comment:'开始执行时间'" json:"startTime"`
	EndTime    int64         `gorm:"column:end_time;type:int;default:0;comment:'执行完成（失败）时间'" json:"endTime"`
	Result     string        `gorm:"column:result;type:mediumtext;comment:'执行结果'" json:"result"`
	Status     Status `gorm:"column:status;type:tinyint;default:1;comment:'状态 1:执行失败 2:执行中  3:执行完毕 4:任务取消(上次任务未执行完成) 5:异步执行'" json:"status"`

	BaseModelTime
	ControlBy
	TotalTime int `gorm:"-" json:"totalTime"`
}

func (e *SysTaskLog) TableName() string {
	return "sys_task_log"
}
