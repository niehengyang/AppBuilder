package resp

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, msg string) error {
	return &Error{Code: code, Message: msg}
}

func NewError409() error {
	return NewError(Err_Conflict, "请勿重复操作")
}

func NewError404() error {
	return NewError(Err_NotFound, "找不到所需资源")
}

func NewError403() error {
	return NewError(Err_Forbidden, "无权限操作资源")
}

func NewError400(msg string) error {
	return NewError(Err_ParamsError, "请求参数错误:"+msg)
}

func NewError401(msg string) error {
	if len(msg) == 0 {
		msg = "认证失败"
	}
	return NewError(Err_Unauthorized, msg)
}

func NewError500() error {
	return NewError(Err_InternalServerError, "未知错误，请联系管理员")
}

func NewCreatedError() error {
	return NewError(Err_CreateFailed, "创建失败,请联系管理员")
}

func NewUpdateError() error {
	return NewError(Err_CreateFailed, "更新失败,请联系管理员")
}

func NewDeleteError() error {
	return NewError(Err_DeleteFailed, "删除失败,请联系管理员")
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Data() *ErrorResponse {
	return &ErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}

/*
全局统一的错误响应处理，
使用httpx.Error()时触发
*/
func SetErrorHandler() {
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *Error:
			errData := e.Data()
			switch errData.Code {
			case Err_ParamsError:
				return http.StatusBadRequest, e.Data()
			case Err_Unauthorized:
				return http.StatusUnauthorized, e.Data()
			case Err_Forbidden:
				return http.StatusForbidden, e.Data()
			case Err_NotFound:
				return http.StatusNotFound, e.Data()
			case Err_Conflict:
				return http.StatusConflict, e.Data()
			default:
				return http.StatusInternalServerError, e.Data()
			}

		default:
			return http.StatusInternalServerError, NewError500()
		}
	})
}
