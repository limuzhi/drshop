package model

import "time"

type SysTask struct {
	TaskId           int64                `gorm:"column:task_id;primaryKey;autoIncrement" json:"taskId"`
	Name             string               `gorm:"column:name;type:varchar(32);comment:'任务名称'" json:"name"`
	Level            TaskLevel            `gorm:"column:level;type:varchar(128);comment:'任务等级 1: 主任务 2: 依赖任务'" json:"level"`
	DependencyTaskId string               `gorm:"column:dependency_task_id;type:varchar(128);default:'';comment:'依赖任务ID,多个ID逗号分隔'" json:"dependencyTaskId"`
	DependencyStatus TaskDependencyStatus `gorm:"column:dependency_status;type:varchar(128);default:1;comment:'依赖关系 1:强依赖 主任务执行成功, 依赖任务才会被执行 2:弱依赖'" json:"dependency_status"`
	Spec             string               `gorm:"column:spec;type:varchar(64);comment:'crontab'" json:"spec"`
	Protocol         TaskProtocol         `gorm:"column:protocol;type:tinyint;comment:'协议 1:http 2:系统命令-RPC方式执行命令'" json:"protocol"`
	Command          string               `gorm:"column:command;type:varchar(255);comment:'URL地址或shell命令'" json:"command"`
	HttpMethod       TaskHTTPMethod       `gorm:"column:http_method;type:tinyint;default:1;comment:'http请求方法 1 GET,2POST'" json:"httpMethod"`
	Timeout          int                  `gorm:"column:timeout;type:mediumint;default:0;comment:'任务执行超时时间(单位秒),0不限制'" json:"timeout"`
	Multi            int8                 `gorm:"column:multi;type:tinyint;default:1;comment:'是否允许多实例运行" json:"multi"`
	RetryTimes       int8                 `gorm:"column:retry_times;type:tinyint;default:0;comment:'重试次数" json:"retryTimes"`
	RetryInterval    int16                `gorm:"column:retry_interval;type:smallint;default:0;comment:'重试间隔时间" json:"retryInterval"`
	NotifyStatus     int8                 `gorm:"column:notify_status;type:tinyint;default:1;comment:'任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知" json:"notifyStatus"`
	NotifyType       int8                 `gorm:"column:notify_type;type:tinyint;default:0;comment:'通知类型 1: 邮件 2: slack 3: webhook" json:"notifyType"`
	NotifyReceiverId string               `gorm:"column:notify_receiver_id;type:varchar(256);default:'';comment:'通知接受者ID, sys_task_setting表主键ID，多个ID逗号分隔 json:"notifyReceiverId"`
	NotifyKeyword    string               `gorm:"column:notify_keyword;type:varchar(128);default:'';comment:'关键字'" json:"notifyKeyword"`
	Tag              string               `gorm:"column:tag;type:varchar(32);default:'';comment:'标签'" json:"tag"`
	Remark           string               `gorm:"column:remark;type:varchar(255);default:'';comment:'备注'" json:"remark"`
	Status           Status               `gorm:"column:status;type:tinyint;default:2;comment:'状态（2正常 1停用）'" json:"status"`

	BaseModelTime
	ControlBy
	Hosts       []*SysHost `gorm:"many2many:sys_task_host;foreignKey:task_id;joinForeignKey:task_id;References:host_id;JoinReferences:host_id;" json:"hosts"`
	NextRunTime time.Time  `gorm:"-" json:"nextRunTime"`
}

func (e *SysTask) TableName() string {
	return "sys_task"
}
