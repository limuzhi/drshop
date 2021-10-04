package code

var zhCNText = map[int]string{
	ServerError:        "内部服务器错误",
	TooManyRequests:    "请求过多",
	ParamBindError:     "参数信息错误",
	AuthorizationError: "签名信息错误",
	UrlSignError:       "参数签名错误",
	CacheSetError:      "设置缓存失败",
	CacheGetError:      "获取缓存失败",
	CacheDelError:      "删除缓存失败",
	CacheNotExist:      "缓存不存在",
	ResubmitError:      "请勿重复提交",
	HashIdsEncodeError: "HashID 加密失败",
	HashIdsDecodeError: "HashID 解密失败",
	RBACError:          "暂无访问权限",
	RedisConnectError:  "Redis 连接失败",
	MySQLConnectError:  "MySQL 连接失败",
	WriteConfigError:   "写入配置文件失败",
	SendEmailError:     "发送邮件失败",
	MySQLExecError:     "SQL 执行失败",
	GoVersionError:     "Go 版本不满足要求",
}
