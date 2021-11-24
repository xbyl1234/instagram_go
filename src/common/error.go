package common

type MakeMoneyError struct {
	ErrStr  string
	ErrType ErrType
}

// 实现 `error` 接口
func (this *MakeMoneyError) Error() string {
	return this.ErrStr
}

type ErrType int

var (
	ApiError   ErrType = 0
	PorxyError ErrType = 1
	OtherError ErrType = 2
)
