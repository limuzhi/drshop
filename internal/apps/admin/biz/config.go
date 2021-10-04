package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/vo"
	"drpshop/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type ConfigRepo interface {
	GetConfigByKey(ctx context.Context, key string) (*v1.ConfigInfoByKeyRes, error)
	ConfigDetail(ctx context.Context, id int64) (*v1.ConfigDetailRes, error)
	ConfigList(ctx context.Context, in *vo.ConfigListReq) (*v1.ConfigListRes, error)
	ConfigCreate(ctx context.Context, in *vo.ConfigCreateReq) error
	ConfigUpdate(ctx context.Context, in *vo.ConfigUpdateReq) error
	ConfigDelete(ctx context.Context, ids []int64) error
}

type ConfigController struct {
	configRepo ConfigRepo
	log        *log.Helper
	response.Api
}

func NewConfigController(repo ConfigRepo, logger log.Logger) *ConfigController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/config"))
	return &ConfigController{
		configRepo: repo,
		log:        log,
	}
}

func (uc *ConfigController) GetConfigByKey(c *gin.Context) {
	var req vo.ConfigKeyReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}

	info, err := uc.configRepo.GetConfigByKey(uc.NewContext(c), req.ConfigKey)
	if err != nil {
		uc.Api.Fail(c, 30002, errors.New("获取失败"))
		return
	}
	uc.Api.Success(c, info)
}

func (uc *ConfigController) ConfigList(c *gin.Context) {
	var req vo.ConfigListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}

	relpy, err := uc.configRepo.ConfigList(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 30002, errors.New("获取失败"))
		return
	}
	uc.Api.PageSuccess(c, relpy.List, int(relpy.Total), req.GetPageIndex(), req.GetPageSize())
}
func (uc *ConfigController) ConfigDetail(c *gin.Context) {
	var req vo.ConfigDetailReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}

	info, err := uc.configRepo.ConfigDetail(uc.NewContext(c), req.ConfigId)
	if err != nil {
		uc.Api.Fail(c, 30002, errors.New("获取失败"))
		return
	}
	uc.Api.Success(c, info)
}
func (uc *ConfigController) ConfigCreate(c *gin.Context) {
	var req vo.ConfigCreateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	userInfo := uc.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		uc.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	req.CreateBy = userInfo.UserID
	err := uc.configRepo.ConfigCreate(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 30002, errors.New("获取失败"))
		return
	}
	uc.Api.Success(c, nil)
}
func (uc *ConfigController) ConfigUpdate(c *gin.Context) {
	var req vo.ConfigUpdateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	userInfo := uc.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		uc.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	req.UpdateBy = userInfo.UserID
	err := uc.configRepo.ConfigUpdate(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 30002, errors.New("更新失败"))
		return
	}
	uc.Api.Success(c, nil)
}
func (uc *ConfigController) ConfigDelete(c *gin.Context) {
	var req vo.ConfigDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	if len(req.ConfigIds) == 0 {
		uc.Api.Fail(c, 20001, errors.New("参数不能为空"))
		return
	}
	err := uc.configRepo.ConfigDelete(uc.NewContext(c), req.ConfigIds)
	if err != nil {
		uc.Api.Fail(c, 30002, errors.New("删除失败"))
		return
	}
	uc.Api.Success(c, nil)
}
