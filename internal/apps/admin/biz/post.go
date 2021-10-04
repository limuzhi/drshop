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

type PostRepo interface {
	PostList(ctx context.Context, in *vo.PostListReq) (*v1.PostListRes, error)
	PostDetail(ctx context.Context, id int64) (*v1.PostDetailRes, error)
	PostCreate(ctx context.Context, in *vo.PostCreateReq) error
	PostUpdate(ctx context.Context, in *vo.PostUpdateReq) error
	PostDelete(ctx context.Context, ids []int64) error
}

type PostController struct {
	postRepo PostRepo
	log      *log.Helper
	response.Api
}

func NewPostController(repo PostRepo, logger log.Logger) *PostController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/post"))
	return &PostController{
		postRepo: repo,
		log:      log,
	}
}

func (e *PostController) PostList(c *gin.Context) {
	var req vo.PostListReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.postRepo.PostList(e.NewContext(c), &req)
	if err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.PageSuccess(c, reply.List, int(reply.Total), req.GetPageIndex(), req.GetPageSize())
}
func (e *PostController) PostDetail(c *gin.Context) {
	var req vo.PostDetail
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	reply, err := e.postRepo.PostDetail(e.NewContext(c), req.PostId)
	if err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.Success(c, reply)
}
func (e *PostController) PostCreate(c *gin.Context) {
	var req vo.PostCreateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	userInfo := e.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	req.CreateBy = userInfo.UserID
	if err := e.postRepo.PostCreate(e.NewContext(c), &req); err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *PostController) PostUpdate(c *gin.Context) {
	var req vo.PostUpdateReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	userInfo := e.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		e.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	req.UpdateBy = userInfo.UserID
	if err := e.postRepo.PostUpdate(e.NewContext(c), &req); err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.Success(c, nil)
}
func (e *PostController) PostDelete(c *gin.Context) {
	var req vo.PostDeleteReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		e.Api.Fail(c, 10001, err)
		return
	}
	if len(req.PostIds) == 0 {
		e.Api.Fail(c, 20001, errors.New("参数不能为空"))
		return
	}
	if err := e.postRepo.PostDelete(e.NewContext(c), req.PostIds); err != nil {
		e.Api.Fail(c, 50001, err)
		return
	}
	e.Api.Success(c, nil)
}
