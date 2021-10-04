package notify

import (
	"strconv"
	"strings"
	"time"

	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"

	"github.com/go-gomail/gomail"
	"github.com/go-kratos/kratos/v2/log"
)

type MailUsecase struct {
	log  *log.Helper
	repo biz.SysTemplateRepo
}

func NewMailUsecase(repo biz.SysTemplateRepo, logger log.Logger) *MailUsecase {
	return &MailUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/notify/mail")),
	}
}

func (uc *MailUsecase) Send(msg Message) {
	mailSetting, err := uc.repo.Mail()
	uc.log.Debugf("%+v", mailSetting)
	if err != nil {
		uc.log.Error("#mail#从数据库获取mail配置失败", err)
		return
	}
	if mailSetting.Host == "" {
		uc.log.Error("#mail#Host为空")
		return
	}
	if mailSetting.Port == 0 {
		uc.log.Error("#mail#Port为空")
		return
	}
	if mailSetting.User == "" {
		uc.log.Error("#mail#User为空")
		return
	}
	if mailSetting.Password == "" {
		uc.log.Error("#mail#Password为空")
		return
	}
	msg["content"] = parseNotifyTemplate(mailSetting.Template, msg)
	toUsers := uc.getActiveMailUsers(mailSetting, msg)
	uc.send(mailSetting, toUsers, msg)
}

func (uc *MailUsecase) send(mailSetting *model.Mail, toUsers []string, msg Message) {
	body := msg["content"].(string)
	body = strings.Replace(body, "\n", "<br>", -1)
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", mailSetting.User)
	gomailMessage.SetHeader("To", toUsers...)
	gomailMessage.SetHeader("Subject", "gocron-定时任务通知")
	gomailMessage.SetBody("text/html", body)
	mailer := gomail.NewDialer(mailSetting.Host, mailSetting.Port,
		mailSetting.User, mailSetting.Password)
	maxTimes := 3
	i := 0
	for i < maxTimes {
		err := mailer.DialAndSend(gomailMessage)
		if err == nil {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			uc.log.Errorf("mail#发送消息失败#%s#消息内容-%s", err.Error(), msg["content"])
		}
	}
}

func (uc *MailUsecase) getActiveMailUsers(mailSetting *model.Mail, msg Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	users := []string{}
	for _, v := range mailSetting.MailUsers {
		if inStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			users = append(users, v.Email)
		}
	}

	return users
}

func inStringSlice(slice []string, element string) bool {
	element = strings.TrimSpace(element)
	for _, v := range slice {
		if strings.TrimSpace(v) == element {
			return true
		}
	}

	return false
}
