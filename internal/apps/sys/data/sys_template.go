package data

import (
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data/model"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
)

type sysTemplateRepo struct {
	data *Data
	log  *log.Helper
}

func NewSysTemplateRepo(data *Data, logger log.Logger) biz.SysTemplateRepo {
	return &sysTemplateRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "sys/data/sys_template")),
	}
}

func (r *sysTemplateRepo) InitBasicField() {
	ins := make([]*model.SysTemplate, 0)

	ins = append(ins, &model.SysTemplate{
		Code: model.SlackCode,
		Key:  model.SlackUrlKey,
	})

	ins = append(ins, &model.SysTemplate{
		Code:  model.SlackCode,
		Key:   model.SlackTemplateKey,
		Value: model.SlackTemplate,
	})

	ins = append(ins, &model.SysTemplate{
		Code:  model.MailCode,
		Key:   model.MailServerKey,
		Value: "",
	})
	ins = append(ins, &model.SysTemplate{
		Code:  model.SlackCode,
		Key:   model.SlackTemplateKey,
		Value: model.EmailTemplate,
	})

	ins = append(ins, &model.SysTemplate{
		Code:  model.WebhookCode,
		Key:   model.WebhookTemplateKey,
		Value: model.WebhookTemplate,
	})
	ins = append(ins, &model.SysTemplate{
		Code:  model.WebhookCode,
		Key:   model.WebhookUrlKey,
		Value: "",
	})
	r.data.db.Create(ins)
}

func (r *sysTemplateRepo) Slack() (*model.Slack, error) {
	list := make([]*model.SysTemplate, 0)
	err := r.data.db.Where("code = ?", model.SlackCode).Find(&list).Error
	slack := &model.Slack{}
	if err != nil {
		return nil, err
	}
	r.formatSlack(list, slack)
	return slack, err
}

func (r *sysTemplateRepo) formatSlack(list []*model.SysTemplate, slack *model.Slack) {
	for _, v := range list {
		switch v.Key {
		case model.SlackUrlKey:
			slack.Url = v.Value
		case model.SlackTemplateKey:
			slack.Template = v.Value
		default:
			slack.Channels = append(slack.Channels, model.Channel{
				Id: int(v.Id), Name: v.Value,
			})
		}
	}
}

func (r *sysTemplateRepo) UpdateSlack(url, template string) error {
	tx := r.data.db.Model(&model.SysTemplate{})
	tx = tx.Where("code = ? AND key = ?", model.SlackCode, model.SlackUrlKey)
	_ = tx.Update("value", url).Error

	tx2 := r.data.db.Model(&model.SysTemplate{})
	tx2 = tx2.Where("code = ? AND key = ?", model.SlackCode, model.SlackTemplateKey)
	_ = tx2.Update("value", template).Error
	return nil
}

func (r *sysTemplateRepo) CreateChannel(channel string) (int64, error) {
	insert := &model.SysTemplate{}
	insert.Code = model.SlackCode
	insert.Key = model.SlackChannelKey
	insert.Value = channel
	err := r.data.db.Model(&model.SysTemplate{}).Create(insert).Error
	if err != nil {
		return 0, err
	}
	return insert.Id, err
}

func (r *sysTemplateRepo) IsChannelExist(channel string) bool {
	var info *model.SysTemplate
	tx := r.data.db.Model(&model.SysTemplate{})
	tx = tx.Where("code = ? AND key = ? AND value = ?",
		model.SlackCode, model.SlackChannelKey, channel)
	if err := tx.First(&info).Error; err != nil {
		return false
	}
	return info.Id > 0
}

// 删除slack渠道
func (r *sysTemplateRepo) RemoveChannel(id int) error {
	return r.data.db.Where("code = ? AND key = ? AND id =?",
		model.SlackCode, model.SlackChannelKey, id).Delete(&model.SysTemplate{}).Error
}

func (r *sysTemplateRepo) Mail() (*model.Mail, error) {
	list := make([]*model.SysTemplate, 0)
	err := r.data.db.Where("code = ?", model.MailCode).Find(&list).Error
	mail := &model.Mail{MailUsers: make([]model.MailUser, 0)}
	if err != nil {
		return mail, err
	}
	r.formatMail(list, mail)
	return mail, err
}

func (r *sysTemplateRepo) formatMail(list []*model.SysTemplate, mail *model.Mail) {
	mailUser := model.MailUser{}
	for _, v := range list {
		switch v.Key {
		case model.MailServerKey:
			json.Unmarshal([]byte(v.Value), mail)
		case model.MailUserKey:
			json.Unmarshal([]byte(v.Value), &mailUser)
			mailUser.Id = int(v.Id)
			mail.MailUsers = append(mail.MailUsers, mailUser)
		case model.MailTemplateKey:
			mail.Template = v.Value
		}
	}
}
func (r *sysTemplateRepo) UpdateMail(config, template string) error {
	tx := r.data.db.Model(&model.SysTemplate{})
	tx = tx.Where("code = ? AND key = ?", model.MailCode, model.MailServerKey)
	_ = tx.Update("value", config).Error

	tx2 := r.data.db.Model(&model.SysTemplate{})
	tx2 = tx2.Where("code = ? AND key = ?", model.MailCode, model.MailTemplateKey)
	_ = tx2.Update("value", template).Error
	return nil
}

func (r *sysTemplateRepo) CreateMailUser(username, email string) (int64, error) {
	insert := &model.SysTemplate{}
	insert.Code = model.MailCode
	insert.Key = model.MailUserKey
	mailUser := model.MailUser{Id: 0, Username: username, Email: email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		return 0, err
	}
	insert.Value = string(jsonByte)
	err = r.data.db.Create(insert).Error
	if err != nil {
		return 0, err
	}
	return insert.Id, nil
}

func (r *sysTemplateRepo) RemoveMailUser(id int) error {
	return r.data.db.Where("code = ? AND key = ? AND id =?",
		model.MailCode, model.MailUserKey, id).Delete(&model.SysTemplate{}).Error
}

func (r *sysTemplateRepo) Webhook() (model.WebHook, error) {
	var list []*model.SysTemplate
	err := r.data.db.Where("code = ?", model.WebhookCode).Find(&list).Error
	webHook := model.WebHook{}
	if err != nil {
		return webHook, err
	}

	r.formatWebhook(list, &webHook)

	return webHook, err
}

func (r *sysTemplateRepo) formatWebhook(list []*model.SysTemplate, webHook *model.WebHook) {
	for _, v := range list {
		switch v.Key {
		case model.WebhookUrlKey:
			webHook.Url = v.Value
		case model.WebhookTemplateKey:
			webHook.Template = v.Value
		}

	}
}
func (r *sysTemplateRepo) UpdateWebHook(url, template string) error {
	tx := r.data.db.Model(&model.SysTemplate{})
	tx = tx.Where("code = ? AND key = ?", model.WebhookCode, model.WebhookUrlKey)
	_ = tx.Update("value", url).Error

	tx2 := r.data.db.Model(&model.SysTemplate{})
	tx2 = tx2.Where("code = ? AND key = ?", model.WebhookCode, model.WebhookTemplateKey)
	_ = tx2.Update("value", template).Error
	return nil
}
