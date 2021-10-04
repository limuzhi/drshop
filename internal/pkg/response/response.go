package response

import (
	"drpshop/internal/pkg/code"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code      int         `json:"code"`      // 错误码((200:成功, 400:失败, 非200:错误码))
	Data      interface{} `json:"data"`      // 返回数据(业务接口定义具体数据结构)
	Msg       string      `json:"msg"`       //提示信息
	RequestId string      `json:"requestId"` //请求requestId
	Err       string      `json:"-"`         //err错误信息
}

type PageResponse struct {
	List      interface{} `json:"list"`      // 返回数据(业务接口定义具体数据结构)
	Count     int         `json:"count"`     // 总数
	PageIndex int         `json:"pageIndex"` //当前页码
	PageSize  int         `json:"pageSize"`  //当前显示页数
}

// 返回前端-成功
func Success(c *gin.Context, data interface{}) {
	var res Response
	res.Data = data
	res.Msg = "成功"
	res.RequestId = ""
	res.Code = SUCCESS
	RJosn(c, http.StatusOK, res)
}

func RJosn(c *gin.Context, code int, data interface{}) {
	c.Set("result", data)
	c.AbortWithStatusJSON(code, data)
}

// 返回前端-失败
func Fail(c *gin.Context, BusinessCode int, err error) {
	var res Response
	if err != nil {
		res.Err = err.Error()
		res.Msg = res.Err
	}
	msg := code.Text(BusinessCode)
	if msg != "" {
		res.Msg = msg
	}
	res.Code = BusinessCode
	RJosn(c, http.StatusOK, res)
}

//分页数据处理
func PageSuccess(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int) {
	var res PageResponse
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	Success(c, res)
}
