package diskfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/disk"
	"github.com/goatcms/goatcore/varutil"
)

// Filespace is a files set on local filesystem
type Filespace struct {
	path string
}

// NewFilespace create new Filespace instance
func NewFilespace(path string) (filesystem.Filespace, error) {
	varutil.FixDirPath(&path)
	return filesystem.Filespace(&Filespace{
		path: path,
	}), nil
}

// Copy a file or a directory inside the filespace
func (fs *Filespace) Copy(src, dest string) error {
	return disk.Copy(fs.path+src, fs.path+dest)
}

func (fs *Filespace) CopyDirectory(src, dest string) error {
	return disk.CopyDirectory(fs.path+src, fs.path+dest)
}

func (fs *Filespace) CopyFile(src, dest string) error {
	return disk.CopyFile(fs.path+src, fs.path+dest)
}

func (fs *Filespace) ReadDir(subPath string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(fs.path + subPath)
}

func (fs *Filespace) IsExist(subPath string) bool {
	return disk.IsExist(fs.path + subPath)
}

func (fs *Filespace) IsFile(subPath string) bool {
	return disk.IsFile(fs.path + subPath)
}

func (fs *Filespace) IsDir(subPath string) bool {
	return disk.IsDir(fs.path + subPath)
}

func (fs *Filespace) MkdirAll(subPath string, filemode os.FileMode) error {
	return disk.MkdirAll(fs.path+subPath, filemode)
}

func (fs *Filespace) Writer(subPath string) (filesystem.Writer, error) {
	return os.OpenFile(fs.path+subPath, os.O_WRONLY|os.O_CREATE, filesystem.DefaultUnixFileMode)
}

func (fs *Filespace) Reader(subPath string) (filesystem.Reader, error) {
	return os.OpenFile(fs.path+subPath, os.O_RDONLY, filesystem.DefaultUnixFileMode)
}

func (fs *Filespace) ReadFile(subPath string) ([]byte, error) {
	return ioutil.ReadFile(fs.path + subPath)
}

func (fs *Filespace) WriteFile(subPath string, data []byte, perm os.FileMode) (err error) {
	path := fs.path + subPath
	if err = fs.MkdirAll(filepath.Dir(path), perm); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}

func (fs *Filespace) Remove(subPath string) error {
	return os.Remove(fs.path + subPath)
}

func (fs *Filespace) RemoveAll(subPath string) error {
	return os.RemoveAll(fs.path + subPath)
}

func (fs *Filespace) Filespace(subPath string) (filesystem.Filespace, error) {
	if !fs.IsDir(subPath) {
		return nil, fmt.Errorf("Path is not a directory " + fs.path + subPath)
	}
	return NewFilespace(fs.path + subPath)
}
