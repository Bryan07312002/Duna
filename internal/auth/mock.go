package auth

// MockSessionStore implementation matching your style
type MockSessionStore struct {
	FuncSet    func(key string, value SessionUserUUIDAndCsrftTokenPair) error
	FuncGet    func(key string) (SessionUserUUIDAndCsrftTokenPair, error)
	FuncRemove func(key string) error
}

func (m *MockSessionStore) Set(key string, value SessionUserUUIDAndCsrftTokenPair) error {
	if m.FuncSet != nil {
		return m.FuncSet(key, value)
	}
	return nil
}

func (m *MockSessionStore) Get(key string) (SessionUserUUIDAndCsrftTokenPair, error) {
	if m.FuncGet != nil {
		return m.FuncGet(key)
	}
	return SessionUserUUIDAndCsrftTokenPair{}, nil
}

func (m *MockSessionStore) Remove(key string) error {
	if m.FuncRemove != nil {
		return m.FuncRemove(key)
	}
	return nil
}
