package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/vo"
	"drpshop/internal/pkg/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"os"
	"path"
	"strings"
	"time"
)

type UserRepo interface {
	Login(c *gin.Context, ctx context.Context, req *vo.LoginReq) (string, int64, error)
	CaptchaImg(ctx context.Context) (id string, b64s string, err error)
	Captcha(ctx context.Context, req *vo.CaptchaReq) error
	Userinfo(ctx context.Context, uid int64) (*v1.UserInfoRes, error)
	InvokeLoginLog(ctx context.Context, data *v1.LoginLogReq)
	InvokeUserOnline(ctx context.Context, data *v1.UserOnlineReq)
	UserList(ctx context.Context, req *vo.UserSearchReq) ([]*v1.UserListData, int, error)

	ChangeStatus(ctx context.Context, uid, status int64) error
	GetPostAndRoleList(ctx context.Context) (*v1.UserPostRoleListRes, error)
	BatchDeleteUser(ctx context.Context, uids []int64) error
	CreateUser(ctx context.Context, req *vo.CreateUserReq) error
	UpdateUserById(ctx context.Context, req *vo.UpdateUserReq) error
	ChangePwd(ctx context.Context, req *vo.ChangePwdReq) error
	GetUserProfile(ctx context.Context, uid int64) (*v1.UserListData, error)
	UserProfileSet(ctx context.Context, in *vo.UserProfileSetReq) error
}

type UserController struct {
	userRepo UserRepo
	log      *log.Helper
	response.Api
}

func NewUserController(repo UserRepo, logger log.Logger) *UserController {
	log := log.NewHelper(log.With(logger, "module", "admin/biz/user"))
	return &UserController{
		userRepo: repo,
		log:      log,
	}
}

func (uc *UserController) Login(c *gin.Context) {
	var req vo.LoginReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	token, expire, err := uc.userRepo.Login(c, uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 10002, err)
		return
	}
	data := make(map[string]interface{})
	data["token"] = token
	data["expire"] = expire
	uc.Api.Success(c, data)
}

func (uc *UserController) Logout(c *gin.Context) {
	//TODO
	uc.Api.Success(c, nil)
}

func (uc *UserController) CaptchaImg(c *gin.Context) {
	id, b64s, err := uc.userRepo.CaptchaImg(uc.NewContext(c))
	if err != nil {
		uc.Api.Fail(c, 10003, err)
		return
	}
	data := make(map[string]interface{})
	data["uuid"] = id
	data["b64s"] = b64s
	uc.Api.Success(c, data)
}

func (uc *UserController) Captcha(c *gin.Context) {
	var req vo.CaptchaReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	err := uc.userRepo.Captcha(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 10003, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *UserController) GetUserDetail(c *gin.Context) {
	userId := uc.GetUserId(c)
	if userId <= 0 {
		uc.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	info, err := uc.userRepo.Userinfo(uc.NewContext(c), userId)
	if err != nil {
		uc.Api.Fail(c, 20001, errors.New("获取用户信息失败"))
		return
	}
	data := make(map[string]interface{})
	if info.Avatar == "" {
		info.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	}
	data["uid"] = userId
	data["avatar"] = info.Avatar
	data["username"] = info.UserName
	data["introduction"] = info.Introduction
	data["roles"] = info.RoleKeys
	data["permissions"] = info.Permissions
	uc.Api.Success(c, data)

}

func (uc *UserController) UserList(c *gin.Context) {
	var req vo.UserSearchReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	list, total, err := uc.userRepo.UserList(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 10101, err)
		return
	}
	uc.Api.PageSuccess(c, list, total, req.GetPageIndex(), req.GetPageSize())
}

func (uc *UserController) ChangeStatus(c *gin.Context) {
	var req vo.UserChangeStatusReq
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
	if userInfo.UserID == req.UserId {
		uc.Api.Fail(c, 10002, errors.New("不能更改自己的状态"))
		return
	}
	//TODO  用户角色以下用户
	err := uc.userRepo.ChangeStatus(uc.NewContext(c), req.UserId, req.Status)
	if err != nil {
		uc.Api.Fail(c, 20003, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *UserController) GetPostAndRoleList(c *gin.Context) {
	res, err := uc.userRepo.GetPostAndRoleList(uc.NewContext(c))
	if err != nil {
		uc.Api.Fail(c, 20001, err)
		return
	}
	uc.Api.Success(c, res)
}

func (uc *UserController) BatchDeleteUser(c *gin.Context) {
	var req vo.BatchDeleteUserReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	if len(req.UserIds) == 0 {
		uc.Api.Fail(c, 10001, errors.New("参数不能为空"))
		return
	}
	err := uc.userRepo.BatchDeleteUser(uc.NewContext(c), req.UserIds)
	if err != nil {
		uc.Api.Fail(c, 20003, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *UserController) UserCreate(c *gin.Context) {
	var req vo.CreateUserReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	req.CreateBy = uc.GetUserId(c)
	err := uc.userRepo.CreateUser(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 20005, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var req vo.UpdateUserReq
	// 参数绑定
	if err := vo.DefalutGetValidParams(c, &req); err != nil {
		uc.Api.Fail(c, 10001, err)
		return
	}
	req.UpdateBy = uc.GetUserId(c)
	err := uc.userRepo.UpdateUserById(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 20005, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *UserController) ChangePwd(c *gin.Context) {
	var req vo.ChangePwdReq
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
	req.UserId = userInfo.UserID
	err := uc.userRepo.ChangePwd(uc.NewContext(c), &req)
	if err != nil {
		uc.Api.Fail(c, 20005, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (uc *UserController) UserPofile(c *gin.Context) {
	userInfo := uc.GetUserInfo(c)
	if userInfo == nil || userInfo.UserID <= 0 {
		uc.Api.Fail(c, 20001, errors.New("未登录"))
		return
	}
	info, err := uc.userRepo.GetUserProfile(uc.NewContext(c), userInfo.UserID)
	if err != nil {
		uc.Api.Fail(c, 20005, err)
		return
	}
	uc.Api.Success(c, info)
}

func (uc *UserController) UserProfileSet(c *gin.Context) {
	var req vo.UserProfileSetReq
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
	req.UserId = userInfo.UserID
	if err := uc.userRepo.UserProfileSet(uc.NewContext(c), &req); err != nil {
		uc.Api.Fail(c, 20005, err)
		return
	}
	uc.Api.Success(c, nil)
}

func (e *UserController) UploadAvatar(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		e.Api.Fail(c, 10001, errors.New("上传文件参数未找到"))
		return
	}

	//builder只可以追加，不可以替换
	var filepath strings.Builder
	fileExt := strings.ToLower(path.Ext(f.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		e.Api.Fail(c, 10002, errors.New("上传失败!只允许png,jpg,gif,jpeg文件"))
		return
	}

	fileDir := fmt.Sprintf("%s/%d-%d-%d/",
		"static/uploadfile/", time.Now().Year(), time.Now().Month(), time.Now().Day())
	filepath.WriteString(fileDir)
	fileName := uuid.New().String()
	filepath.WriteString(fileName)
	if _, err := os.Stat(fileDir); err != nil {
		if err := os.Mkdir(fileDir, os.ModePerm); err != nil {
			e.Api.Fail(c, 10002, errors.New("文件夹创建失败"))
			return
		}
	}
	filepath.WriteString(fileExt)
	if err := c.SaveUploadedFile(f, filepath.String()); err != nil {
		e.Api.Fail(c, 10002, errors.New("上传失败"))
		return
	}

	e.Api.Success(c, map[string]string{"path": filepath.String()})
}
