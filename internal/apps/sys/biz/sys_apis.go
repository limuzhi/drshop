package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/internal/apps/sys/global"
	"drpshop/pkg/errors/normal"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/thoas/go-funk"
)

type SysApisRepo interface {
	//列表
	ListApis(ctx context.Context, in *v1.ApisListReq) ([]*model.SysApis, int64, error)
	//创建
	CreateApis(ctx context.Context, in *model.SysApis) error
	//更新
	UpdateApis(ctx context.Context, apiId int64, in *model.SysApis) error
	//根据接口ID获取接口列表
	GetApisById(ctx context.Context, ids []int64) ([]*model.SysApis, error)
	//批量删除接口
	BatchDeleteByIds(ctx context.Context, ids []int64) error
	//获取接口树(按接口Category字段分类)
	GetListOrderCategory(ctx context.Context) ([]*model.SysApis, error)
	//根据接口路径和请求方式获取接口描述
	GetApiDescByPath(ctx context.Context, path string, method string) (string, error)

	GetListByPath(ctx context.Context, paths []string) ([]*model.SysApis, error)
}

type SysApisUsecase struct {
	repo SysApisRepo
	log  *log.Helper
}

func NewSysApisUsecase(repo SysApisRepo, logger log.Logger) *SysApisUsecase {
	return &SysApisUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "sys/biz/sys_apis"))}
}

//列表
func (uc *SysApisUsecase) ListApis(ctx context.Context, in *v1.ApisListReq) ([]*model.SysApis, int64, error) {
	if in.PageInfo.PageNum <= 0 {
		in.PageInfo.PageNum = 1
	}
	if in.PageInfo.PageSize <= 0 {
		in.PageInfo.PageSize = 10
	}
	return uc.repo.ListApis(ctx, in)
}

//创建
func (uc *SysApisUsecase) CreateApis(ctx context.Context, in *model.SysApis) error {
	return uc.repo.CreateApis(ctx, in)
}

//更新
func (uc *SysApisUsecase) UpdateApis(ctx context.Context, apiId int64, in *model.SysApis) error {
	return uc.repo.UpdateApis(ctx, apiId, in)
}

//根据接口ID获取接口列表
func (uc *SysApisUsecase) GetApisById(ctx context.Context, ids []int64) ([]*model.SysApis, error) {
	return uc.repo.GetApisById(ctx, ids)
}

//批量删除接口
func (uc *SysApisUsecase) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteByIds(ctx, ids)
}

//获取接口树(按接口Category字段分类)
func (uc *SysApisUsecase) GetListOrderCategory(ctx context.Context) ([]*model.SysApis, error) {
	return uc.repo.GetListOrderCategory(ctx)
}

//根据接口路径和请求方式获取接口描述
func (uc *SysApisUsecase) GetApiDescByPath(ctx context.Context, path string, method string) (string, error) {
	return uc.repo.GetApiDescByPath(ctx, path, method)
}

func (uc *SysApisUsecase) GetListByPath(ctx context.Context, paths []string) ([]*model.SysApis, error) {
	return uc.repo.GetListByPath(ctx, paths)
}

//TODO---
func (uc *SysApisUsecase) GetApisInfo(ctx context.Context, id int64) (*model.SysApis, error) {
	list, err := uc.repo.GetApisById(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, normal.RecordNotFound
	}
	return list[0], nil
}

//获取接口树(按接口Category字段分类)
func (uc *SysApisUsecase) GetApisListTree(ctx context.Context) ([]*v1.ApisTreeData, error) {
	list, err := uc.repo.GetListOrderCategory(ctx)
	if err != nil {
		return nil, err
	}
	// 获取所有的分类
	var categoryList []string
	for _, api := range list {
		categoryList = append(categoryList, api.Category)
	}
	// 获取去重后的分类
	categoryUniq := funk.UniqString(categoryList)

	apiTree := make([]*v1.ApisTreeData, len(categoryUniq))

	for i, category := range categoryUniq {
		apiTree[i] = &v1.ApisTreeData{
			ApiId:    int64(-i),
			Title:    category,
			Category: category,
			Children: nil,
		}
		for _, api := range list {
			if category == api.Category {
				info := &v1.ApisTreeData{
					ApiId:    api.ApiId,
					Title:    api.Title + "-(" + api.Path + "|" + api.Method + ")",
					Category: api.Category,
				}
				apiTree[i].Children = append(apiTree[i].Children, info)
			}
		}
	}
	return apiTree, err
}

func (uc *SysApisUsecase) DtoOut(data *model.SysApis) *v1.ApisInfoRes {
	info := &v1.ApisInfoRes{
		ApiId:      data.ApiId,
		Title:      data.Title,
		Path:       data.Path,
		Method:     data.Method,
		Category:   data.Category,
		Permission: data.Permission,
		CreatedAt:  global.GetDateByUnix(data.CreatedAt),
		UpdatedAt:  global.GetDateByUnix(data.UpdatedAt),
		CreateBy:   data.CreateBy,
		UpdateBy:   data.UpdateBy,
	}
	return info
}
