package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/global"
	"drpshop/internal/apps/admin/middle"
	"drpshop/internal/conf"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/sdk/trace"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdminService)

type AdminService struct {
	Logger   log.Logger
	DataData *data.Data
	Middle   middleware.Middleware
	Router   *gin.Engine
	Casbinc  *biz.CasbinController
}

func NewAdminService(c *conf.Data, registry *conf.Registry, logger log.Logger,
	tp *trace.TracerProvider, m middleware.Middleware) (*AdminService, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	// 初始化Validator数据校验
	global.InitValidate()

	discovery := data.NewDiscovery(registry)
	sysClient := data.NewSysServiceClient(discovery, tp)
	dataData, err := data.NewData(c, logger, rdb, sysClient)
	if err != nil {
		return nil, nil, err
	}
	casbinRepo := data.NewCasbinRepo(dataData, logger)
	casbinc := biz.NewCasbinController(casbinRepo, logger)
	return &AdminService{
		Logger:   logger,
		DataData: dataData,
		Router:   gin.Default(),
		Middle:   m,
		Casbinc:  casbinc,
	}, func() {}, nil
}

func (s *AdminService) RegisterSysServiceServer() *gin.Engine {
	s.Router.Use(kgin.Middlewares(s.Middle), middle.OperationLogMiddleware())
	//TODO 检测TOKEN 中间件还需一个

	// 无需认证的路由
	noCheckRoleRouter(s)
	//需要认证的路由
	checkRoleRouter(s)
	return s.Router
}
