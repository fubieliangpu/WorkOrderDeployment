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
