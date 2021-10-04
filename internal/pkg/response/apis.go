package response

import (
	"context"
	"drpshop/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Api struct {
}

//通常错误数据处理
func (e *Api) Fail(c *gin.Context, code int, err error) {
	Fail(c, code, err)
}

//通常成功数据处理
func (e *Api) Success(c *gin.Context, data interface{}) {
	Success(c, data)
}

//分页数据处理
func (e *Api) PageSuccess(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int) {
	PageSuccess(c, result, count, pageIndex, pageSize)
}

//兼容函数
func (e *Api) Custum(c *gin.Context, data gin.H) {
	RJosn(c, http.StatusOK, data)
}

func (e *Api) NewContext(c *gin.Context) context.Context {
	return token.NewContextMetadataClientToken(c.Request.Context())
}

func (e *Api) GetUserId(c *gin.Context) int64 {
	return token.FormGlobalUidContext(c.Request.Context())
}

func (e *Api) GetUserInfo(c *gin.Context) *token.UserClaims {
	return token.FormLoginContext(c.Request.Context())
}
