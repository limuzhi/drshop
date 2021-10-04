package vo

type ConfigKeyReq struct {
	ConfigKey string `form:"configKey" json:"configKey" validate:"required"`
}

type ConfigListReq struct {
	ConfigName string `form:"configName" json:"configName" validate:""`
	ConfigKey  string `form:"configKey" json:"configKey" validate:""`
	ConfigType string `form:"configType" json:"configType" validate:""`
	Pagination
}

type ConfigCreateReq struct {
	ConfigName  string `form:"configName" json:"configName" validate:"required,min=1,max=100"`
	ConfigKey   string `form:"configKey" json:"configKey" validate:"required,min=1,max=100"`
	ConfigValue string `form:"configValue" json:"configValue" validate:"min=1,max=500"`
	ConfigType  string `form:"configType" json:"configType" validate:"oneof=1 2"`
	IsFrontend  string `form:"isFrontend" json:"isFrontend" validate:"oneof=1 2"`
	Remark      string `form:"remark" json:"remark" validate:""`
	CreateBy    int64  `form:"-" json:"-"`
}

type ConfigUpdateReq struct {
	ConfigCreateReq
	ConfigId int64 `form:"configId" json:"configId" validate:"required,gte=1"`
	UpdateBy int64 `form:"-" json:"-"`
}

type ConfigDetailReq struct {
	ConfigId int64 `form:"configId" json:"configId" validate:"required,gte=1"`
}

type ConfigDeleteReq struct {
	ConfigIds []int64 `form:"configIds" json:"configIds" validate:"required,gte=1"`
}
