package postgres

import "os"

// MockFileSystem is a mock implementation of FileSystem for testing
type MockFileSystem struct {
	ReadFileFunc func(name string) ([]byte, error)
	ReadDirFunc  func(name string) ([]os.DirEntry, error)
}

func (m MockFileSystem) ReadFile(name string) ([]byte, error) {
	return m.ReadFileFunc(name)
}

func (m MockFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return m.ReadDirFunc(name)
}

// MockDirEntry is a mock implementation of os.DirEntry for testing
type MockDirEntry struct {
	name  string
	isDir bool
}

func (m MockDirEntry) Name() string               { return m.name }
func (m MockDirEntry) IsDir() bool                { return m.isDir }
func (m MockDirEntry) Type() os.FileMode          { return 0 }
func (m MockDirEntry) Info() (os.FileInfo, error) { return nil, nil }
