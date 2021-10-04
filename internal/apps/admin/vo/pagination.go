/**
* @Author: lh
* @Description:
* @File: pagination
* @Version: 1.0.0
* @Date: 2021/4/10 15:44
 */
package vo

import (
	"strings"
	"time"
)

type Pagination struct {
	PageNum   int    `json:"pageNum" form:"pageNum"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	BeginTime string `json:"beginTime" form:"beginTime"` //开始日期
	EndTime   string `json:"endTime" form:"endTime"`     //结束日期
}

func (m *Pagination) GetPageIndex() int {
	if m.PageNum <= 0 {
		m.PageNum = 1
	}
	return m.PageNum
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

func (m *Pagination) GetBeginTime() int64 {
	str := strings.TrimSpace(m.BeginTime)
	if str == "" {
		return 0
	}
	beginTime, _ := time.ParseInLocation("2006-01-02", str, time.Local)
	return beginTime.Unix()
}

func (m *Pagination) GetEndTime() int64 {
	str := strings.TrimSpace(m.EndTime)
	if str == "" {
		return 0
	}
	endTime, _ := time.ParseInLocation("2006-01-02", str, time.Local)
	return endTime.Unix() + 86400
}
