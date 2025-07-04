package user

type HashStrategyMock struct {
	FuncEncode  func(str string) (string, error)
	FuncCompare func(enconded, str string) bool
}

func (h HashStrategyMock) Encode(str string) (string, error) {
	return h.FuncEncode(str)
}

func (h HashStrategyMock) Compare(enconded, str string) bool {
	return h.FuncCompare(enconded, str)
}
