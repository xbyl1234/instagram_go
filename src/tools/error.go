package tools

type MakeMoneyError struct {
	ErrStr  string
	ErrType int
}

// 实现 `error` 接口
func (this *MakeMoneyError) Error() string {
	return this.ErrStr
}
