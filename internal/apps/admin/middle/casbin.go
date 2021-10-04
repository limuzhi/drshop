package middle

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

func CasbinMiddleware(casbinc *biz.CasbinController) gin.HandlerFunc {
	return func(c *gin.Context) {
		e := &response.Api{}
		userInfo := e.GetUserInfo(c)
		if userInfo == nil || userInfo.UserID <= 0 {
			e.Fail(c, 20001, errors.New("用户未登录"))
			c.Abort()
			return
		}
		if userInfo.UserID != 1 {
			// 获得请求路径URL
			obj := c.FullPath()
			// 获取请求方式
			act := c.Request.Method //c.Request.URL.Path
			log.Println("obj:", obj, "====act:", act, "====")
			isPass := casbinc.CheckAuth(e.NewContext(c), userInfo.UserID, obj, act)
			if !isPass {
				e.Fail(c, 20002, errors.New("没有权限"))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
