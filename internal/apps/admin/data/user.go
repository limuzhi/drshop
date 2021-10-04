package data

import (
	"context"
	"drpshop/api/common"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/vo"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/util/gconv"
	"github.com/mssola/user_agent"
	"strings"
	"time"
)

type userRepo struct {
	data *Data
	log  *log.Helper
	Pool *grpool.Pool
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		Pool: grpool.New(200),
		log:  log.NewHelper(log.With(logger, "module", "admin/data/user")),
	}
}

func (r *userRepo) Login(c *gin.Context, ctx context.Context, req *vo.LoginReq) (string, int64, error) {
	if req.Password == "" || req.Account == "" {
		return "", 0, errors.Unauthorized("LoginError", "账号密码不能为空")
	}
	if req.Code == "" || req.UUID == "" {
		return "", 0, errors.Unauthorized("LoginError", "验证码不能为空")
	}
	reply, err := r.data.sc.Login(ctx, &v1.LoginReq{
		Account:   strings.TrimSpace(req.Account),  //账户
		Password:  strings.TrimSpace(req.Password), // 密码
		Captcha:   strings.TrimSpace(req.Code),     // 验证码
		CaptchaId: strings.TrimSpace(req.UUID),     // 验证码id
		Code:      "",                              //短信、邮箱 验证码
	})
	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()
	ua := user_agent.New(userAgent)
	os := ua.OS()
	browser, _ := ua.Browser()
	loginlogData := &v1.LoginLogReq{
		LoginName:     req.Account,
		LoginUid:      0,
		Ipaddr:        ip,
		LoginLocation: "",
		Browser:       browser,
		Os:            ua.OS(),
		Status:        1,
		LoginTime:     time.Now().Unix(),
		Module:        "系统后台",
	}
	if err != nil {
		//保存日志（异步）
		loginlogData.Msg = err.Error()
		r.InvokeLoginLog(ctx, loginlogData)
		return "", 0, err
	}
	userInfo, err := r.data.Jk.ParseToken(reply.Token)
	uid := int64(0)
	username := ""
	if err != nil {
		loginlogData.Msg = "token 加密错误"
		loginlogData.Status = 1
	} else {
		uid = userInfo.UserID
		username = userInfo.Username
		loginlogData.Msg = "登录成功"
		loginlogData.Status = 2
	}
	loginlogData.LoginUid = uid
	loginlogData.LoginName = username
	r.InvokeLoginLog(ctx, loginlogData)
	if loginlogData.Status == 2 && uid > 0 {
		//更新登录IP
		r.InvokeLoginIpTime(ctx, uid, loginlogData.LoginTime, ip)
	}
	//保存用户在线状态(异步)
	onlineData := &v1.UserOnlineReq{
		UserId:     uid,
		Token:      reply.Token,
		CreateTime: time.Now().Unix(),
		Username:   username,
		Ip:         c.ClientIP(),
		Explorer:   browser,
		Os:         os,
	}
	r.InvokeUserOnline(ctx, onlineData)

	return reply.Token, reply.Expire, nil
}

func (r *userRepo) CaptchaImg(ctx context.Context) (id string, b64s string, err error) {
	reply, err := r.data.sc.CaptchaImg(ctx, &v1.CaptchaImgReq{
		Height: 80,
		Width:  240,
		Length: 4,
	})
	if err != nil {
		return
	}
	id = reply.CaptchaId
	b64s = reply.B64S
	return
}

func (r *userRepo) Captcha(ctx context.Context, req *vo.CaptchaReq) error {
	_, err := r.data.sc.Captcha(ctx, &v1.CaptchaReq{
		To:        strings.TrimSpace(req.To),
		Captcha:   strings.TrimSpace(req.Captcha),
		CaptchaId: strings.TrimSpace(req.CaptchaId),
	})
	return err
}

func (r *userRepo) Userinfo(ctx context.Context, uid int64) (*v1.UserInfoRes, error) {
	reply, err := r.data.sc.UserInfo(ctx, &v1.UserInfoReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *userRepo) InvokeLoginLog(ctx context.Context, data *v1.LoginLogReq) {
	r.Pool.Add(func() {
		//写入日志数据
		r.data.sc.SaveLoginlog(ctx, data)
	})
}

func (r *userRepo) InvokeLoginIpTime(ctx context.Context, userId, loginTime int64, ip string) {
	r.Pool.Add(func() {
		//写入日志数据
		r.data.sc.LoginUserUpdate(ctx, &v1.LoginUserUpdateReq{
			UserId:        userId,
			LastLoginTime: loginTime,
			LastLoginIp:   ip,
		})
	})
}

func (r *userRepo) InvokeUserOnline(ctx context.Context, data *v1.UserOnlineReq) {
	r.Pool.Add(func() {
		//写入数据
		r.data.sc.SaveUserOnline(ctx, data)
	})
}

func (r *userRepo) UserList(ctx context.Context, req *vo.UserSearchReq) ([]*v1.UserListData, int, error) {
	reply, err := r.data.sc.UserList(ctx, &v1.UserListReq{
		PageNum:   int64(req.GetPageIndex()),
		PageSize:  int64(req.GetPageSize()),
		Username:  strings.TrimSpace(req.Username),
		Mobile:    strings.TrimSpace(req.Mobile),
		Email:     "",
		Status:    req.Status,
		DeptId:    int64(req.DeptId),
		StartTime: req.GetBeginTime(),
		EndTime:   req.GetEndTime(),
	})
	if err != nil {
		return nil, 0, err
	}
	return reply.List, int(reply.Total), err
}

func (r *userRepo) ChangeStatus(ctx context.Context, uid, status int64) error {
	_, err := r.data.sc.UpdateUserStatus(ctx, &v1.UserStatusReq{
		UserId: uid,
		Status: status,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetPostAndRoleList(ctx context.Context) (*v1.UserPostRoleListRes, error) {
	return r.data.sc.GetUserPostRoleList(ctx, &common.NullReq{})
}

func (r *userRepo) BatchDeleteUser(ctx context.Context, uids []int64) error {
	_, err := r.data.sc.UserDelete(ctx, &v1.UserDeleteReq{UserIds: uids})
	return err
}

func (r *userRepo) CreateUser(ctx context.Context, req *vo.CreateUserReq) error {
	_, err := r.data.sc.UserAdd(ctx, &v1.UserAddReq{
		Username: strings.TrimSpace(req.Username),
		Nickname: strings.TrimSpace(req.Nickname),
		Mobile:   strings.TrimSpace(req.Mobile),
		Salt:     strings.TrimSpace(req.Salt),
		Email:    strings.TrimSpace(req.Email),
		DeptId:   req.DeptId,
		Sex:      gconv.Int64(req.Sex),
		Avatar:   strings.TrimSpace(req.Avatar),
		Password: strings.TrimSpace(req.Password),
		Address:  "",
		Remark:   strings.TrimSpace(req.Remark),
		IsAdmin:  req.IsAdmin,
		Status:   strings.TrimSpace(req.Status),
		CreateBy: req.CreateBy,
		RoleIds:  req.RoleIds,
		PostIds:  req.PostIds,
	})
	return err
}

func (r *userRepo) UpdateUserById(ctx context.Context, req *vo.UpdateUserReq) error {
	_, err := r.data.sc.UserUpdate(ctx, &v1.UserUpdateReq{
		UserId:   req.UserId,
		Nickname: strings.TrimSpace(req.Nickname),
		Mobile:   strings.TrimSpace(req.Mobile),
		Email:    strings.TrimSpace(req.Email),
		DeptId:   req.DeptId,
		Sex:      gconv.Int64(req.Sex),
		Avatar:   strings.TrimSpace(req.Avatar),
		Address:  "",
		Remark:   strings.TrimSpace(req.Remark),
		IsAdmin:  req.IsAdmin,
		Status:   strings.TrimSpace(req.Status),
		UpdateBy: req.UpdateBy,
		RoleIds:  req.RoleIds,
		PostIds:  req.PostIds,
	})
	return err
}

func (r *userRepo) ChangePwd(ctx context.Context, req *vo.ChangePwdReq) error {
	_, err := r.data.sc.ChangePwd(ctx, &v1.ChangePwdReq{
		UserId:      req.UserId,
		NewPassword: strings.TrimSpace(req.NewPassword),
		OldPassword: strings.TrimSpace(req.OldPassword),
	})
	return err
}

func (r *userRepo) GetUserProfile(ctx context.Context, uid int64) (*v1.UserListData, error) {
	return r.data.sc.UserProfile(ctx, &v1.UserInfoReq{UserId: uid})
}

func (r *userRepo) UserProfileSet(ctx context.Context, in *vo.UserProfileSetReq) error {
	_, err := r.data.sc.UserProfileSet(ctx, &v1.UserProfileSetReq{
		UserId:   in.UserId,
		Nickname: strings.TrimSpace(in.Nickname),
		Mobile:   strings.TrimSpace(in.Mobile),
		Email:    strings.TrimSpace(in.Email),
		Avatar:   strings.TrimSpace(in.Avatar),
		Sex:      in.Sex,
	})
	return err
}
