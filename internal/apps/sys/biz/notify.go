package biz

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"html/template"
	"time"
)

type Message map[string]interface{}

type Notifiable interface {
	Send(msg Message)
}

type NotifyUsecase struct {
	log        *log.Helper
	mailSrv    *MailUsecase
	slackSrv   *SlackUsecase
	webhookSrv *WebhookUsecase
	repo       SysTemplateRepo
}

func NewNotifyUsecase(repo SysTemplateRepo, logger log.Logger) *NotifyUsecase {
	return &NotifyUsecase{
		repo:       repo,
		mailSrv:    NewMailUsecase(repo, logger),
		slackSrv:   NewSlackUsecase(repo, logger),
		webhookSrv: NewWebhookUsecase(repo, logger),
		log:        log.NewHelper(log.With(logger, "module", "sys/notify")),
	}
}

var queue = make(chan Message, 100)

// 把消息推入队列
func Push(msg Message) {
	queue <- msg
}

func (uc *NotifyUsecase) Run() {
	for msg := range queue {
		// 根据任务配置发送通知
		taskType, taskTypeOk := msg["task_type"]
		_, taskReceiverIdOk := msg["task_receiver_id"]
		_, nameOk := msg["name"]
		_, outputOk := msg["output"]
		_, statusOk := msg["status"]
		if !taskTypeOk || !taskReceiverIdOk || !nameOk || !outputOk || !statusOk {
			uc.log.Errorf("#notify#参数不完整#%+v", msg)
			continue
		}
		msg["content"] = fmt.Sprintf("============\n============\n============\n任务名称: %s\n状态: %s\n输出:\n %s\n", msg["name"], msg["status"], msg["output"])
		uc.log.Debugf("%+v", msg)
		switch taskType.(int8) {
		case 1:
			// 邮件
			go uc.mailSrv.Send(msg)
		case 2:
			// Slack
			go uc.slackSrv.Send(msg)
		case 3:
			// WebHook
			go uc.webhookSrv.Send(msg)
		}
		time.Sleep(1 * time.Second)
	}
}

func parseNotifyTemplate(notifyTemplate string, msg Message) string {
	tmpl, err := template.New("notify").Parse(notifyTemplate)
	if err != nil {
		return fmt.Sprintf("解析通知模板失败: %s", err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, map[string]interface{}{
		"TaskId":   msg["task_id"],
		"TaskName": msg["name"],
		"Status":   msg["status"],
		"Result":   msg["output"],
		"Remark":   msg["remark"],
	})

	return buf.String()
}
