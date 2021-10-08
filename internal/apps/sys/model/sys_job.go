package model

type SysJob struct {
	JobId          int64  `gorm:"column:job_id;primaryKey;autoIncrement" json:"jobId"`
	JobName        string `gorm:"column:job_name;type:varchar(128);comment:'任务名称'" json:"jobName"`
	JobParams      string `gorm:"column:job_params;type:varchar(255);comment:'参数'" json:"jobParams"`
	JobGroup       string `gorm:"column:job_group;type:varchar(64);default:'DEFAULT';comment:'任务组名'" json:"jobGroup"`
	InvokeTarget   string `gorm:"column:invoke_target;type:varchar(500);comment:'调用目标字符串'" json:"invokeTarget"`
	CronExpression string `gorm:"column:cron_expression;type:varchar(255);comment:'cron执行表达式'" json:"cronExpression"`
	MisfirePolicy  int    `gorm:"column:misfire_policy;type:tinyint;default:1;comment:'计划执行策略（1多次执行 2执行一次）'" json:"misfirePolicy"`
	Concurrent     int    `gorm:"column:concurrent;type:tinyint;default:1;comment:'是否并发执行（2允许 1禁止）'" json:"concurrent"`
	Status         int    `gorm:"column:status;type:tinyint(1);cdefault:2;omment:'状态（2正常 1停用'" json:"status"`
	Remark         string `gorm:"column:remark;type:varchar(255);comment:'备注'" json:"remark"`
	BaseModelTime
	ControlBy
}

func (e *SysJob) TableName() string {
	return "sys_job"
}
