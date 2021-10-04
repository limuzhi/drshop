package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type operLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewOperLogRepo(data *Data, logger log.Logger) biz.OperLogRepo {
	return &operLogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "admin/data/oper_log")),
	}
}

func (r *operLogRepo) GetList(ctx context.Context, req *v1.OperLogListReq) (*v1.OperLogListRes, error) {
	return r.data.sc.OperLogList(ctx, req)
}
func (r *operLogRepo) GetDetail(ctx context.Context, id int64) (*v1.OperLogInfoRes, error) {
	return r.data.sc.OperLogInfo(ctx, &v1.OperLogInfoReq{OperId: id})
}
func (r *operLogRepo) BatchDelete(ctx context.Context, ids []int64) error {
	_, err := r.data.sc.OperLogDelete(ctx, &v1.OperLogDeleteReq{
		OperId: ids,
	})
	return err
}
func (r *operLogRepo) Clear(ctx context.Context) error {
	_, err := r.data.sc.OperLogClear(ctx, &common.NullReq{})
	return err
}

func (r *operLogRepo) SaveOperlogChannel(out <-chan *v1.OperLogSaveData) {
	// 只会在线程开启的时候执行一次
	Logs := make([]*v1.OperLogSaveData, 0)

	// 一直执行--收到就会执行
	for in := range out {
		Logs = append(Logs, in)
		// 每10条记录到数据库
		if len(Logs) > 5 {
			r.data.sc.OperLogSave(context.Background(), &v1.OperLogSaveReq{LogList: Logs})
			Logs = make([]*v1.OperLogSaveData, 0)
		}
	}

}
