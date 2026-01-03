package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

type LocalStorage struct {
	dir     string
	maxSize int64
}

func NewLocalStorage(dir string, maxSize int64) *LocalStorage {
	return &LocalStorage{dir: dir, maxSize: maxSize}
}

func (s *LocalStorage) Save(ctx context.Context, upload graphql.Upload) (string, error) {
	_ = ctx

	if upload.Size > 0 && upload.Size > s.maxSize {
		return "", fmt.Errorf("file too large")
	}

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(upload.Filename))
	name := uuid.NewString() + ext
	path := filepath.Join(s.dir, name)

	if closer, ok := upload.File.(io.Closer); ok {
		defer closer.Close()
	}

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	limited := &io.LimitedReader{R: upload.File, N: s.maxSize + 1}
	written, err := io.Copy(out, limited)
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	if written > s.maxSize {
		_ = os.Remove(path)
		return "", fmt.Errorf("file too large")
	}

	return filepath.ToSlash(filepath.Join("/uploads", name)), nil
}
