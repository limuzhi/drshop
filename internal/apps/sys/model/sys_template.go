package model

type SysTemplate struct {
	Id    int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code  string `gorm:"column:code;type:varchar(64);comment:'参数名称'" json:"code"`
	Key   string `gorm:"column:key;type:varchar(64);comment:'参数键名'" json:"key"`
	Value string `gorm:"column:value;type:text;comment:'参数键值'" json:"value"`
}

func (e *SysTemplate) TableName() string {
	return "sys_template"
}

const SlackTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
备注: {{.Remark}}
`
const EmailTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
备注: {{.Remark}}
`
const WebhookTemplate = `
{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}",
  "remark": "{{.Remark}}"
}
`

const (
	SlackCode        = "slack"
	SlackUrlKey      = "url"
	SlackTemplateKey = "template"
	SlackChannelKey  = "channel"
)

const (
	MailCode        = "mail"
	MailTemplateKey = "template"
	MailServerKey   = "server"
	MailUserKey     = "user"
)

const (
	WebhookCode        = "webhook"
	WebhookTemplateKey = "template"
	WebhookUrlKey      = "url"
)

// region slack配置

type Slack struct {
	Url      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Template string    `json:"template"`
}

type Channel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Mail struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	MailUsers []MailUser `json:"mail_users"`
	Template  string     `json:"template"`
}

type MailUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type WebHook struct {
	Url      string `json:"url"`
	Template string `json:"template"`
}
