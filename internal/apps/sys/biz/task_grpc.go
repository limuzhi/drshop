package biz

import (
	v1 "drpshop/api/task/v1"
	"github.com/go-kratos/kratos/v2/log"
)

type TaskGrpcPoolRepo interface {
	Get(addr string) (v1.TaskServiceClient, error)
	Release(addr string) //释放连接
	Stop(ip string, port int, id int64)
	Exec(ip string, port int, taskReq *v1.TaskReq) (string, error)
}

type TaskGrpcUsecase struct {
	repo TaskGrpcPoolRepo
	log  *log.Helper
}

func NewTaskGrpcUsecase(repo TaskGrpcPoolRepo, logger log.Logger) *TaskGrpcUsecase {
	return &TaskGrpcUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "sys/biz/task_grpc")),
	}
}
