package model

type SysOperLog struct {
	OperId       int64  `gorm:"column:oper_id;primaryKey;autoIncrement" json:"operId"`
	Title         string `gorm:"column:title;type:varchar(50);comment:'模块标题'" json:"title"`
	BusinessType  int    `gorm:"column:business_type;type:int;default:0;comment:'业务类型（0其它 1新增 2修改 3删除）'" json:"businessType"`
	Method        string `gorm:"column:method;type:varchar(128);comment:'方法名称'" json:"method"`
	RequestMethod string `gorm:"column:request_method;type:varchar(10);comment:'请求方式'" json:"requestMethod"`
	OperatorType  int    `gorm:"column:operator_type;type:int;default:0;comment:'操作类别（0其它 1后台用户 2手机端用户）'" json:"operatorType"`
	OperName      string `gorm:"column:oper_name;type:varchar(128);comment:'操作人员'" json:"operName"`
	OperUrl       string `gorm:"column:oper_url;type:varchar(500);comment:'请求URL'" json:"operUrl"`
	OperIp        string `gorm:"column:oper_ip;type:varchar(50);comment:'主机地址'" json:"operIp"`
	OperLocation  string `gorm:"column:oper_location;type:varchar(255);comment:'操作地点'" json:"operLocation"`
	OperParam     string `gorm:"column:oper_param;type:text;comment:'请求参数'" json:"operParam"`
	JsonResult    string `gorm:"column:json_result;type:text;comment:'返回参数'" json:"jsonResult"`
	Status        string `gorm:"column:status;type:int;default:2;comment:'操作状态（2正常 1异常）'" json:"status"`
	ErrorMsg      string `gorm:"column:error_msg;type:varchar(2000);comment:'错误消息'" json:"errorMsg"`
	OperTime      int64  `gorm:"column:oper_time;type:int;comment:'操作时间'" json:"operTime"`
	TimeCost      int64  `gorm:"column:time_cost;type:int;comment:'请求耗时(ms)'" json:"timeCost"`
	UserId        int64  `gorm:"column:user_id;type:int;comment:'用户uid'" json:"userId"`
}

func (e *SysOperLog) TableName() string {
	return "sys_oper_log"
}
