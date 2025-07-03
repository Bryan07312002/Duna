package postgres

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestNewMigration(t *testing.T) {
	t.Run("with nil filesystem", func(t *testing.T) {
		m := NewMigration("test", 123, "/path", nil)
		if m.name != "test" {
			t.Errorf("expected name 'test', got '%s'", m.name)
		}
		if m.timestamp != 123 {
			t.Errorf("expected timestamp 123, got %d", m.timestamp)
		}
		if m.path != "/path" {
			t.Errorf("expected path '/path', got '%s'", m.path)
		}
		if _, ok := m.fs.(DefaultFileSystem); !ok {
			t.Error("expected DefaultFileSystem when nil is passed")
		}
	})

	t.Run("with custom filesystem", func(t *testing.T) {
		// Create a custom filesystem with a unique behavior we can test
		customCalled := false
		fs := MockFileSystem{
			ReadFileFunc: func(name string) ([]byte, error) {
				customCalled = true
				return nil, nil
			},
		}

		m := NewMigration("test", 123, "/path", fs)

		// Instead of comparing the filesystem directly, we'll verify behavior
		_, _ = m.fs.ReadFile("test") // Should call our custom implementation
		if !customCalled {
			t.Error("custom filesystem implementation was not used")
		}
	})
}

func TestFullName(t *testing.T) {
	m := &Migration{
		name:      "create-users-table",
		timestamp: 1234567890,
	}
	expected := "1234567890-create-users-table"
	if actual := m.FullName(); actual != expected {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}

func TestGetUpQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedQuery := "CREATE TABLE users (id SERIAL PRIMARY KEY);"
		fs := MockFileSystem{
			ReadFileFunc: func(name string) ([]byte, error) {
				if filepath.Base(name) != "up.sql" {
					t.Errorf("expected to read up.sql, got %s", name)
				}
				return []byte(expectedQuery), nil
			},
		}

		m := &Migration{
			name:      "create-users-table",
			timestamp: 1234567890,
			path:      "/migrations",
			fs:        fs,
		}

		query, err := m.GetUpQuery()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if query != expectedQuery {
			t.Errorf("expected query '%s', got '%s'", expectedQuery, query)
		}
	})

	t.Run("file read error", func(t *testing.T) {
		expectedErr := errors.New("file not found")
		fs := MockFileSystem{
			ReadFileFunc: func(name string) ([]byte, error) {
				return nil, expectedErr
			},
		}

		m := &Migration{
			name:      "create-users-table",
			timestamp: 1234567890,
			path:      "/migrations",
			fs:        fs,
		}

		_, err := m.GetUpQuery()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error '%v', got '%v'", expectedErr, err)
		}
	})
}

func TestGetDownQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedQuery := "DROP TABLE users;"
		fs := MockFileSystem{
			ReadFileFunc: func(name string) ([]byte, error) {
				if filepath.Base(name) != "down.sql" {
					t.Errorf("expected to read down.sql, got %s", name)
				}
				return []byte(expectedQuery), nil
			},
		}

		m := &Migration{
			name:      "create-users-table",
			timestamp: 1234567890,
			path:      "/migrations",
			fs:        fs,
		}

		query, err := m.GetDownQuery()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if query != expectedQuery {
			t.Errorf("expected query '%s', got '%s'", expectedQuery, query)
		}
	})

	t.Run("file read error", func(t *testing.T) {
		expectedErr := errors.New("file not found")
		fs := MockFileSystem{
			ReadFileFunc: func(name string) ([]byte, error) {
				return nil, expectedErr
			},
		}

		m := &Migration{
			name:      "create-users-table",
			timestamp: 1234567890,
			path:      "/migrations",
			fs:        fs,
		}

		_, err := m.GetDownQuery()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error '%v', got '%v'", expectedErr, err)
		}
	})
}

func TestReadMigrationDir(t *testing.T) {
	tests := []struct {
		name        string
		dirEntries  []os.DirEntry
		want        []*Migration
		wantErr     bool
		expectError error
	}{
		{
			name: "valid migrations",
			dirEntries: []os.DirEntry{
				MockDirEntry{name: "123-create-table", isDir: true},
				MockDirEntry{name: "456-drop-table", isDir: true},
				MockDirEntry{name: "789-alter-table", isDir: true},
			},
			want: []*Migration{
				{name: "create-table", timestamp: 123, path: "/migrations"},
				{name: "drop-table", timestamp: 456, path: "/migrations"},
				{name: "alter-table", timestamp: 789, path: "/migrations"},
			},
			wantErr: false,
		},
		{
			name: "ignore files",
			dirEntries: []os.DirEntry{
				MockDirEntry{name: "123-create-table", isDir: true},
				MockDirEntry{name: "README.md", isDir: false},
				MockDirEntry{name: "456-drop-table", isDir: true},
			},
			want: []*Migration{
				{name: "create-table", timestamp: 123, path: "/migrations"},
				{name: "drop-table", timestamp: 456, path: "/migrations"},
			},
			wantErr: false,
		},
		{
			name: "invalid timestamp",
			dirEntries: []os.DirEntry{
				MockDirEntry{name: "abc-create-table", isDir: true},
			},
			wantErr:     true,
			expectError: strconv.ErrSyntax,
		},
		{
			name: "sorting order",
			dirEntries: []os.DirEntry{
				MockDirEntry{name: "123-old-migration", isDir: true},
				MockDirEntry{name: "456-middle-migration", isDir: true},
				MockDirEntry{name: "789-new-migration", isDir: true},
			},
			want: []*Migration{
				{name: "old-migration", timestamp: 123, path: "/migrations"},
				{name: "middle-migration", timestamp: 456, path: "/migrations"},
				{name: "new-migration", timestamp: 789, path: "/migrations"},
			},
			wantErr: false,
		},
		{
			name: "realistic timestamps",
			dirEntries: []os.DirEntry{
				MockDirEntry{name: "1000-very-old", isDir: true},
				MockDirEntry{name: "2000-old", isDir: true},
				MockDirEntry{name: "3000-recent", isDir: true},
			},
			want: []*Migration{
				{name: "very-old", timestamp: 1000, path: "/migrations"},
				{name: "old", timestamp: 2000, path: "/migrations"},
				{name: "recent", timestamp: 3000, path: "/migrations"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := MockFileSystem{
				ReadDirFunc: func(name string) ([]os.DirEntry, error) {
					return tt.dirEntries, nil
				},
			}

			got, err := ReadMigrationDir("/migrations", fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMigrationDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectError != nil && !errors.Is(err, tt.expectError) {
					t.Errorf("ReadMigrationDir() error = %v, expectError %v", err, tt.expectError)
			}

			if len(got) != len(tt.want) {
				t.Errorf("ReadMigrationDir() returned %d migrations, want %d", len(got), len(tt.want))
				return
			}

			for i, gotMigration := range got {
				if gotMigration.name != tt.want[i].name {
					t.Errorf("Migration %d name = %v, want %v", i, gotMigration.name, tt.want[i].name)
				}
				if gotMigration.timestamp != tt.want[i].timestamp {
					t.Errorf("Migration %d timestamp = %v, want %v", i, gotMigration.timestamp, tt.want[i].timestamp)
				}
			}

			// Verify sorting (newest first)
			for i := 0; i < len(got)-1; i++ {
				if got[i].timestamp > got[i+1].timestamp {
					t.Errorf("Migrations not sorted correctly: %d comes before %d", got[i].timestamp, got[i+1].timestamp)
				}
			}
		})
	}

	t.Run("read dir error", func(t *testing.T) {
		expectedErr := errors.New("directory not found")
		fs := MockFileSystem{
			ReadDirFunc: func(name string) ([]os.DirEntry, error) {
				return nil, expectedErr
			},
		}

		_, err := ReadMigrationDir("/nonexistent", fs)
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error '%v', got '%v'", expectedErr, err)
		}
	})
}
