package error_codes

// 公共错误代码
const (
	SUCCESS            = 0
	FAIL               = 1
	InvalidParam       = 2
	UnAuth             = 3
	NotFound           = 4
	DbErr              = 5
	CacheErr           = 6
	CreateFileFail     = 7
	SignError          = 8
	GrpcSysErr         = 9
	ConfigErr          = 10
	Unknown            = 11
	DeadlineExceeded   = 12
	AccessDenied       = 13
	LimitExceed        = 14
	MethodNotAllowed   = 15
	ServiceUnavailable = 16
	TokenExpired       = 17
	TokenInvalid       = 18
	TicketInvalid      = 19
	PhoneEmpty         = 20
	LicenseExpired     = 21
)
