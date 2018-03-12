package filestore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const metaFileSuffix = ".contenttype"

func NewLocal(rootPath string) Interface {
	return &localFS{rootPath: rootPath}
}

type localFS struct {
	rootPath string
}

func (fs *localFS) PutUserFile(ctx context.Context, user, path string, src io.Reader, contentType string) error {
	path, err := fs.resolvePath(user, path)
	if err != nil {
		return err
	}
	dest, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path+".contenttype", []byte(contentType), 0600)
	if err != nil {
		return err
	}
	return nil
}

func (fs *localFS) GetUserFile(ctx context.Context, user, path string) (io.ReadCloser, string, error) {
	path, err := fs.resolvePath(user, path)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(path) // read-only
	if err != nil {
		return nil, "", err
	}
	contentType, err := ioutil.ReadFile(path + metaFileSuffix)
	return file, string(contentType), err
}

func (fs *localFS) DeleteUserFile(ctx context.Context, user, path string) error {
	path, err := fs.resolvePath(user, path)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return os.Remove(path + metaFileSuffix)
}

func (fs *localFS) ListFiles(ctx context.Context, user string) ([]string, error) {
	path, err := fs.resolvePath(user, "")
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(path)
	var fnames []string
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), metaFileSuffix) {
			fnames = append(fnames, f.Name())
		}
	}
	return fnames, err
}

func (fs *localFS) resolvePath(user, path string) (string, error) {
	userRoot := filepath.Join(fs.rootPath, user)
	os.MkdirAll(userRoot, 0700)
	path = filepath.Join(userRoot, path)
	rel, err := filepath.Rel(userRoot, path)
	if err != nil {
		return "", err
	}
	if strings.Contains(rel, "..") {
		return "", fmt.Errorf("path may not contain ..") // path escapes user's root
	}
	if strings.Contains(rel, "/") {
		return "", fmt.Errorf("path may not have subdirectories") // don't allow subfolders
	}
	if strings.HasSuffix(rel, metaFileSuffix) {
		return "", fmt.Errorf("path may not end with " + metaFileSuffix) // don't allow intersection with our meta files
	}
	return path, nil
}
