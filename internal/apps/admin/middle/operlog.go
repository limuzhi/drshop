package middle

import (
	"bufio"
	"bytes"
	v1 "drpshop/api/sys/v1"
	"drpshop/pkg/token"
	"drpshop/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gconv"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 操作日志channel
var OperLogChan = make(chan *v1.OperLogSaveData, 30)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		var operParam string
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete:
			if c.Request.Method == http.MethodGet {
				query := c.Request.URL.Query()
				b, _ := json.Marshal(query)
				operParam = string(b)
			} else {
				bf := bytes.NewBuffer(nil)
				wt := bufio.NewWriter(bf)
				_, err := io.Copy(wt, c.Request.Body)
				if err != nil {
					err = nil
				}
				body, _ := ioutil.ReadAll(bf)
				contentType := c.Request.Header.Get("Content-type")
				if strings.Index(contentType, "application/json") > -1 {
					var queryMap map[string]interface{}
					_ = json.Unmarshal(body, &queryMap)
					b, _ := json.Marshal(queryMap)
					operParam = string(b)
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		if c.Request.Method == http.MethodOptions {
			return
		}

		// 执行耗时
		timeCost := endTime.Sub(startTime).Milliseconds()
		// 获取当前登录用户
		userInfo := token.FormLoginContext(c.Request.Context())
		if userInfo == nil {
			userInfo = &token.UserClaims{}
		}
		var (
			result   = ""
			status   = "1"
			errorMsg = ""
		)
		if rt, ok := c.Get("result"); ok {
			rb, _ := json.Marshal(rt)
			if rb != nil {
				result = string(rb)
			}
			var resultMap map[string]interface{}
			_ = json.Unmarshal(rb, &resultMap)

			for k, v := range resultMap {
				switch k {
				case "code":
					if gconv.Int(v) == 200 {
						status = "2"
					}
				case "err":
					errorMsg = gconv.String(v)
				}
			}
		}
		// 获取访问路径
		operLog := v1.OperLogSaveData{
			Title:         "模块标题",
			BusinessType:  "0",                //业务类型（0其它 1新增 2修改 3删除）
			Method:        c.Request.URL.Path, //方法名称
			RequestMethod: c.Request.Method,   //请求方式
			OperatorType:  "1",                //操作类别（0其它 1后台用户 2手机端用户）
			OperName:      userInfo.Username,
			OperUrl:       c.Request.RequestURI,
			OperIp:        c.ClientIP(), //主机地址
			OperLocation:  util.GetCityByIp(c.ClientIP()),
			OperParam:     operParam,
			JsonResult:    result,
			Status:        status, //操作状态（2正常 1异常）
			ErrorMsg:      errorMsg,
			OperTime:      time.Now().Unix(),
			TimeCost:      timeCost,
		}
		// 最好是将日志发送到rabbitmq或者kafka中
		// 这里是发送到channel中，开启3个goroutine处理
		OperLogChan <- &operLog
	}
}

