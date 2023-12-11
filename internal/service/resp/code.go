package resp

const (
	/* ------------ 通用错误码 ------------- */
	Err_ParamsError         = 4000 //参数错误
	Err_Unauthorized        = 4010 //身份认证失败
	Err_Forbidden           = 4030 //无权限
	Err_NotFound            = 4040 //找不到所需资源
	Err_Conflict            = 4090 // 重复/冲突
	Err_InternalServerError = 5000 //服务器错误

	/* ------------ 自定义错误码 ------------- */
	Err_LoginFailed  = 5001 //登录失败
	Err_CreateFailed = 5002 //创建失败
	Err_DeleteFailed = 5003 //删除失败
	Err_NaviRefused  = 5004 //数字人助手繁忙，拒绝请求
)
