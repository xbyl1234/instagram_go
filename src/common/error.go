package common

import "fmt"

type MakeMoneyError struct {
	ErrStr    string
	ErrType   ErrType
	ExternErr error
}

// 实现 `error` 接口
func (this *MakeMoneyError) Error() string {
	errHead := ""
	switch this.ErrType {
	case ApiError:
		errHead = "api error"
		break
	case PorxyError:
		errHead = "proxy error"
		break
	case NoMoreError:
		errHead = "no more error"
		break
	case OtherError:
		errHead = "other error"
		break
	case RequestError:
		errHead = "request error"
		break
	}
	if this.ErrStr != "" {
		errHead += "," + this.ErrStr
	}
	if this.ExternErr != nil {
		errHead += "," + fmt.Sprintf("%v", this.ExternErr)
	}
	return errHead
}

type ErrType int

var (
	ApiError     ErrType = 0
	PorxyError   ErrType = 1
	NoMoreError  ErrType = 2
	OtherError   ErrType = 3
	RequestError ErrType = 4
)

func IsNoMoreError(err error) bool {
	e, ok := err.(*MakeMoneyError)
	if ok {
		return e.ErrType == NoMoreError
	}
	return false
}

func IsError(err error, errType ErrType) bool {
	e, ok := err.(*MakeMoneyError)
	if ok {
		return e.ErrType == errType
	}
	return false
}
