package error_codes

var codeTextDict = map[interface{}]string{
	SUCCESS:            "成功",
	FAIL:               "服务内部错误",
	InvalidParam:       "非法请求参数",
	UnAuth:             "无访问权限",
	NotFound:           "找不到资源",
	DbErr:              "数据库出错",
	CacheErr:           "缓存出错",
	CreateFileFail:     "创建文件失败",
	SignError:          "签名验证失败",
	GrpcSysErr:         "系统错误",
	ConfigErr:          "配置错误",
	Unknown:            "未知错误",
	DeadlineExceeded:   "操作超时",
	AccessDenied:       "拒绝访问",
	LimitExceed:        "请求过多，请稍后重试",
	MethodNotAllowed:   "方法不被允许",
	ServiceUnavailable: "服务暂不可用，请稍后重试",
	TokenExpired:       "TOKEN过期",
	TokenInvalid:       "非法TOKEN",
	TicketInvalid:      "非法Ticket",
	PhoneEmpty:         "手机号为空",
	LicenseExpired:     "License非法或者过期",
}

var (
	Success            = NewWithCode(SUCCESS)
	SystemError        = NewWithCode(FAIL)
	InvalidParamError  = NewWithCode(InvalidParam)
	UnAuthError        = NewWithCode(UnAuth)
	NotFoundError      = NewWithCode(NotFound)
	DatabaseError      = NewWithCode(DbErr)
	InvalidTicketError = NewWithCode(TicketInvalid)
)
