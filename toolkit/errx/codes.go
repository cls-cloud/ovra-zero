package errx

const (
	CodeUnknown  = 1000 // 未知错误
	CodeInvalid  = 1001 // 参数错误
	CodeInternal = 1002 // 系统错误

	CodeLoginErr = 1100 // 登录失败
	CodeNoUser   = 1101 // 用户不存在

	CodeBizErr = 1200 // 业务错误
	CodeNoPerm = 1201 // 无权限

	CodeNoData     = 1300 // 数据未找到
	CodeOrmInvalid = 1301 // ORM错误
)
