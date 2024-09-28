package exception

import "fmt"

func ErrServerInternal(format string, a ...any) *ApiException {
	return &ApiException{
		Code:    50000,
		Message: fmt.Sprintf(format, a...),
	}
}

func ErrNotFound(format string, a ...any) *ApiException {
	return &ApiException{
		Code:    404,
		Message: fmt.Sprintf(format, a...),
	}
}

func ErrValidateFailed(format string, a ...any) *ApiException {
	return &ApiException{
		Code:    400,
		Message: fmt.Sprintf(format, a...),
	}
}

//没能成功读取文件
func ErrOpenFileFailed(format string, a ...any) *ApiException {
	return &ApiException{
		Code:    50005,
		Message: fmt.Sprintf(format, a...),
	}
}

//没能成功解析设备登录的用户信息的yaml文件
func ErrParseFileFailed(format string, a ...any) *ApiException {
	return &ApiException{
		Code:    50008,
		Message: fmt.Sprintf(format, a...),
	}
}

//正则匹配失败
func ErrRegularMatchFailed(format string, a ...any) *ApiException {
	return &ApiException{
		Code:    50444,
		Message: fmt.Sprintf(format, a...),
	}
}
