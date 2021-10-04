package notify

import (
	"drpshop/internal/apps/sys/biz"
	"drpshop/pkg/httpclient"
	"github.com/go-kratos/kratos/v2/log"
	"html"
	"time"
)

type WebhookUsecase struct {
	log  *log.Helper
	repo biz.SysTemplateRepo
}

func NewWebhookUsecase(repo biz.SysTemplateRepo, logger log.Logger) *WebhookUsecase {
	return &WebhookUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/notify/webhook")),
	}
}


func (uc *WebhookUsecase) Send(msg Message) {
	webHookSetting, err := uc.repo.Webhook()
	if err != nil {
		uc.log.Error("#webHook#从数据库获取webHook配置失败", err)
		return
	}
	if webHookSetting.Url == "" {
		uc.log.Error("#webHook#webhook-url为空")
		return
	}
	uc.log.Debugf("%+v", webHookSetting)
	msg["name"] = escapeJson(msg["name"].(string))
	msg["output"] = escapeJson(msg["output"].(string))
	msg["content"] = parseNotifyTemplate(webHookSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	uc.send(msg, webHookSetting.Url)
}

func (uc *WebhookUsecase) send(msg Message, url string) {
	content := msg["content"].(string)
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(url, content, timeout)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			uc.log.Errorf("webHook#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}
