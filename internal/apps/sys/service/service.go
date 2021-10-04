package service

import (
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/data"
	"drpshop/internal/apps/sys/service/notify"
	"drpshop/internal/conf"
	"drpshop/pkg/captcha"
	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewSysService, NewTaskService, NewNotifyService)

type SysService struct {
	v1.UnimplementedSysServiceServer
	apic        *biz.SysApisUsecase
	configc     *biz.SysConfigUsecase
	deptc       *biz.SysDeptUsecase
	dictDatac   *biz.SysDictDataUsecase
	dictTypec   *biz.SysDictTypeUsecase
	hostc       *biz.SysHostUsecase
	jobc        *biz.SysJobUsecase
	loginLogc   *biz.SysLoginLogUsecase
	menuc       *biz.SysMenuUsecase
	operlogc    *biz.SysOperLogUsecase
	postc       *biz.SysPostUsecase
	rolec       *biz.SysRoleUsecase
	taskc       *biz.SysTaskUsecase
	templcatec  *biz.SysTemplateUsecase
	userc       *biz.SysUserUsecase
	userOnlinec *biz.SysUserOnlineUsecase

	taskSrc *TaskService
	log     *log.Helper
	store   *captcha.RedisStore
	domain  string
}

func NewSysService(c *conf.Data, u *conf.UserConfig, logger log.Logger) (*SysService, func(), error) {
	dataData, cleanup, err := data.NewData(c, logger, u)
	if err != nil {
		return nil, nil, err
	}
	apisRepo := data.NewSysApisRepo(dataData, logger)
	configRepo := data.NewSysConfigRepo(dataData, logger)
	deptRepo := data.NewSysDeptRepo(dataData, logger)
	dictDataRepo := data.NewSysDictDataRepo(dataData, logger)
	dictTypeRepo := data.NewSysDictTypeRepo(dataData, logger)
	hostRepo := data.NewSysHostRepo(dataData, logger)
	jobRepo := data.NewSysJobRepo(dataData, logger)
	loginLogRepo := data.NewSysLoginLogRepo(dataData, logger)
	menuRepo := data.NewSysMenuRepo(dataData, logger)
	operlogRepo := data.NewSysOperLogRepo(dataData, logger)
	postRepo := data.NewSysPostRepo(dataData, logger)
	roleRepo := data.NewSysRoleRepo(dataData, logger)
	taskRepo := data.NewSysTaskRepo(dataData, logger)
	templateRepo := data.NewSysTemplateRepo(dataData, logger)
	userRepo := data.NewSysUserRepo(dataData, logger)
	userOnlineRepo := data.NewSysUserOnlineRepo(dataData, logger)

	taskSrc := NewTaskService(taskRepo, logger)
	taskSrc.Initialize()

	go func() {
		notifySrcc := notify.NewNotifyUsecase(templateRepo, logger)
		notifySrcc.Run()
	}()

	return &SysService{
		apic:        biz.NewSysApisUsecase(apisRepo, logger),
		configc:     biz.NewSysConfigUsecase(configRepo, logger),
		deptc:       biz.NewSysDeptUsecase(deptRepo, logger),
		dictDatac:   biz.NewSysDictDataUsecase(dictDataRepo, logger),
		dictTypec:   biz.NewSysDictTypeUsecase(dictTypeRepo, logger),
		hostc:       biz.NewSysHostUsecase(hostRepo, logger),
		jobc:        biz.NewSysJobUsecase(jobRepo, logger),
		loginLogc:   biz.NewSysLoginLogUsecase(loginLogRepo, logger),
		menuc:       biz.NewSysMenuUsecase(menuRepo, userRepo, roleRepo, logger),
		operlogc:    biz.NewSysOperLogUsecase(operlogRepo, logger),
		postc:       biz.NewSysPostUsecase(postRepo, logger),
		rolec:       biz.NewSysRoleUsecase(roleRepo, logger),
		taskc:       biz.NewSysTaskUsecase(taskRepo, logger),
		templcatec:  biz.NewSysTemplateUsecase(templateRepo, logger),
		userc:       biz.NewSysUserUsecase(userRepo, logger),
		userOnlinec: biz.NewSysUserOnlineUsecase(userOnlineRepo, logger),
		taskSrc:     taskSrc,
		log:         log.NewHelper(log.With(logger, "module", "service/sys-service")),
		store: captcha.NewRedisStore(redis.NewClient(&redis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DB:           int(c.Redis.Db),
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})),
		domain: "drpshop-sys",
	}, cleanup, err
}

type NotifyService struct {
	log *log.Helper
}

func NewNotifyService(logger log.Logger) *NotifyService {
	return &NotifyService{
		log: log.NewHelper(log.With(logger, "module", "service/notify-service")),
	}
}
