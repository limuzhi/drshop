package biz

import (
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysTemplateRepo interface {
	InitBasicField()
	Slack() (*model.Slack, error)
	UpdateSlack(url, template string) error
	//创建slack渠道
	CreateChannel(channel string) (int64, error)
	IsChannelExist(channel string) bool
	RemoveChannel(id int) error
	Mail() (*model.Mail, error)
	UpdateMail(config, template string) error
	CreateMailUser(username, email string) (int64, error)
	RemoveMailUser(id int) error
	Webhook() (model.WebHook, error)
	UpdateWebHook(url, template string) error
}

type SysTemplateUsecase struct {
	repo SysTemplateRepo
	log  *log.Helper
}

func NewSysTemplateUsecase(repo SysTemplateRepo, logger log.Logger) *SysTemplateUsecase {
	return &SysTemplateUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/notify/sys_template")),
	}
}

func (uc *SysTemplateUsecase) InitBasicField() {
	uc.repo.InitBasicField()
}
func (uc *SysTemplateUsecase) Slack() (*model.Slack, error) {
	return uc.repo.Slack()
}
func (uc *SysTemplateUsecase) UpdateSlack(url, template string) error {
	return uc.repo.UpdateSlack(url, template)
}

//创建slack渠道
func (uc *SysTemplateUsecase) CreateChannel(channel string) (int64, error) {
	return uc.repo.CreateChannel(channel)
}

func (uc *SysTemplateUsecase) IsChannelExist(channel string) bool {
	return uc.repo.IsChannelExist(channel)
}

func (uc *SysTemplateUsecase) RemoveChannel(id int) error {
	return uc.repo.RemoveChannel(id)
}

func (uc *SysTemplateUsecase) Mail() (*model.Mail, error) {
	return uc.repo.Mail()
}

func (uc *SysTemplateUsecase) UpdateMail(config, template string) error {
	return uc.repo.UpdateMail(config, template)
}

func (uc *SysTemplateUsecase) CreateMailUser(username, email string) (int64, error) {
	return uc.repo.CreateMailUser(username, email)
}

func (uc *SysTemplateUsecase) RemoveMailUser(id int) error {
	return uc.repo.RemoveMailUser(id)
}

func (uc *SysTemplateUsecase) Webhook() (model.WebHook, error) {
	return uc.repo.Webhook()
}

func (uc *SysTemplateUsecase) UpdateWebHook(url, template string) error {
	return uc.repo.UpdateWebHook(url, template)
}
