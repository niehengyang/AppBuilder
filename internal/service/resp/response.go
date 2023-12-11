package resp

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Message: "success",
		Data:    data,
	}
}

//func NewError400() *Response{
//	return &Response{
//		Code: 40000,
//		Message: "请求参数错误",
//	}
//}
//
//func NewError401() *Response{
//	return &Response{
//		Code: 40101,
//		Message: "认证失败",
//	}
//}
//
//func NewError403() *Response{
//	return &Response{
//		Code: 40404,
//		Message: "无权限操作资源",
//	}
//}
//
//func NewError404() *Response{
//	return &Response{
//		Code: 40404,
//		Message: "找不到所需的资源",
//	}
//}
//
//func NewError409() *Response{
//	return &Response{
//		Code: 40909,
//		Message: "请勿重复创建",
//	}
//}
//
//func NewError500() *Response{
//	return &Response{
//		Code: 50000,
//		Message: "请勿重复创建",
//	}
//}
