package biz

import (
	"drpshop/internal/apps/sys/model"
	"drpshop/pkg/httpclient"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"html"
	"strconv"
	"strings"
	"time"
)

type SlackUsecase struct {
	log  *log.Helper
	repo SysTemplateRepo
}

func NewSlackUsecase(repo SysTemplateRepo, logger log.Logger) *SlackUsecase {
	return &SlackUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/notify/slack")),
	}
}

func (uc *SlackUsecase) Send(msg Message) {
	slackSetting, err := uc.repo.Slack()
	if err != nil {
		uc.log.Error("#slack#从数据库获取slack配置失败", err)
		return
	}
	if slackSetting.Url == "" {
		uc.log.Error("#slack#webhook-url为空")
		return
	}
	if len(slackSetting.Channels) == 0 {
		uc.log.Error("#slack#channels配置为空")
		return
	}
	uc.log.Debugf("%+v", slackSetting)
	channels := uc.getActiveSlackChannels(slackSetting, msg)
	uc.log.Debugf("%+v", channels)
	msg["content"] = parseNotifyTemplate(slackSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	for _, channel := range channels {
		uc.send(msg, slackSetting.Url, channel)
	}
}

func (uc *SlackUsecase) send(msg Message, slackUrl string, channel string) {
	formatBody := uc.format(msg["content"].(string), channel)
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(slackUrl, formatBody, timeout)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			uc.log.Errorf("slack#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}

func (uc *SlackUsecase) getActiveSlackChannels(slackSetting *model.Slack, msg Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	channels := []string{}
	for _, v := range slackSetting.Channels {
		if inStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			channels = append(channels, v.Name)
		}
	}

	return channels
}

// 格式化消息内容
func (uc *SlackUsecase) format(content string, channel string) string {
	content = escapeJson(content)
	specialChars := []string{"&", "<", ">"}
	replaceChars := []string{"&amp;", "&lt;", "&gt;"}
	content = replaceStrings(content, specialChars, replaceChars)

	return fmt.Sprintf(`{"text":"%s","username":"gocron", "channel":"%s"}`, content, channel)
}

// 转义json特殊字符
func escapeJson(s string) string {
	specialChars := []string{"\\", "\b", "\f", "\n", "\r", "\t", "\""}
	replaceChars := []string{"\\\\", "\\b", "\\f", "\\n", "\\r", "\\t", "\\\""}

	return replaceStrings(s, specialChars, replaceChars)
}

// 批量替换字符串
func replaceStrings(s string, old []string, replace []string) string {
	if s == "" {
		return s
	}
	if len(old) != len(replace) {
		return s
	}

	for i, v := range old {
		s = strings.Replace(s, v, replace[i], 1000)
	}

	return s
}
