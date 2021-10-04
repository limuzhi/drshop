package data

import (
	"context"
	"time"

	taskv1 "drpshop/api/tasknode/v1"
	"drpshop/internal/conf"
	"drpshop/pkg/token"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	grpct "google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDbClient, NewSysUserRepo,
	NewSysRoleRepo, NewSysLoginLogRepo, NewSysApisRepo,
	NewSysDictDataRepo, NewSysDictTypeRepo, NewSysJobRepo,
	NewSysMenuRepo, NewSysPostRepo, NewSysHostRepo,
	NewSysTaskRepo)

// Data .
type Data struct {
	db             *gorm.DB
	rdb            *redis.Client
	tk             *token.Token
	casbinEnforcer *casbin.Enforcer
}

func NewDbClient(conf *conf.Data, logger log.Logger) *gorm.DB {
	log := log.NewHelper(log.With(logger, "module", "sys-service/data/db"))
	opts := gormConfig()
	db, err := gorm.Open(mysql.Open(conf.Database.Source), opts)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	db.Debug()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed opening connection to sqlDB: %v", err)
	}
     //最大链接数
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(int(conf.Database.MaxOpenConnections))

	//最大可复用时间
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(conf.Database.MaxConnectionLifeTime.AsDuration())

	//空闲最大数量
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(int(conf.Database.MaxIdleConnections))

	return db
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, uc *conf.UserConfig) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "sys-service/data"))
	tk, err := token.New(uc.Cert.Key, uc.Cert.Cert)
	if err != nil {
		log.Errorf("failed create jwt ", err)
		return nil, func() {}, err
	}
	db := NewDbClient(c, logger)
	casbinEnforcer, err := NewCasbinEnforcer(c.CasbinModelPath, db)
	if err != nil {
		log.Errorf("failed init casbinEnforcer ", err)
		return nil, func() {}, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})

	d := &Data{
		db:             db,
		rdb:            rdb,
		tk:             tk,
		casbinEnforcer: casbinEnforcer,
	}
	return d, func() {
		if sqlDB, err := db.DB(); err != nil {
			log.Error(err)
		} else {
			if err = sqlDB.Close(); err != nil {
				log.Error(err)
			}
		}
	}, nil
}

func NewDiscovery(addr, scheme string) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = addr
	c.Scheme = scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func NewTaskServiceClient(r registry.Discovery) (*grpct.ClientConn, taskv1.TaskClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint("discovery:///task.node.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	c := taskv1.NewTaskClient(conn)
	return conn, c, nil
}

func NewCasbinEnforcer(rbacModelPath string, db *gorm.DB) (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(rbacModelPath, a)
	if err != nil {
		return nil, err
	}
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
