package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Driver 定义文件系统驱动接口
type Driver interface {
	// 目录操作
	MakeDir(path string) error
	DeleteDir(path string) error
	DirExists(path string) (bool, error)

	// 文件操作
	Put(path string, content io.Reader) error
	Get(path string) (io.ReadCloser, error)
	Delete(path string) error
	FileExists(path string) (bool, error)

	// 通用操作
	Size(path string) (int64, error)
	Copy(src, dst string) error
	Move(src, dst string) error
	List(path string) ([]FileInfo, error)
}

// FileInfo 文件信息
type FileInfo struct {
	Name      string
	Size      int64
	IsDir     bool
	UpdatedAt int64
}

// Config 文件系统配置
type Config struct {
	Type      string // "local", "s3", "oss" 等
	RootPath  string // 本地文件系统根路径
	AccessKey string // 可用于云存储认证
	SecretKey string // 可用于云存储认证
	Bucket    string // 可用于云存储桶名
	Region    string // 可用于云存储区域
}

// NewDriver 创建新的文件系统驱动
func NewDriver(config Config) (Driver, error) {
	switch config.Type {
	case "local":
		return &LocalDriver{rootPath: config.RootPath}, nil
	case "s3":
		// 将来实现
		return nil, errors.New("S3 driver not implemented")
	case "oss":
		// 将来实现
		return nil, errors.New("OSS driver not implemented")
	default:
		return nil, errors.New("unknown filesystem driver type")
	}
}

// LocalDriver 本地文件系统驱动
type LocalDriver struct {
	rootPath string
}

// MakeDir 创建目录
func (d *LocalDriver) MakeDir(path string) error {
	fullPath := filepath.Join(d.rootPath, path)
	return os.MkdirAll(fullPath, 0755)
}

// DirExists 检查目录是否存在
func (d *LocalDriver) DirExists(path string) (bool, error) {
	fullPath := filepath.Join(d.rootPath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// DeleteDir 删除目录
func (d *LocalDriver) DeleteDir(path string) error {
	fullPath := filepath.Join(d.rootPath, path)
	return os.RemoveAll(fullPath)
}

// Put 保存文件
func (d *LocalDriver) Put(path string, content io.Reader) error {
	fullPath := filepath.Join(d.rootPath, path)

	// 确保父目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入内容
	_, err = io.Copy(file, content)
	return err
}

// Get 获取文件
func (d *LocalDriver) Get(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(d.rootPath, path)
	return os.Open(fullPath)
}

// Delete 删除文件
func (d *LocalDriver) Delete(path string) error {
	fullPath := filepath.Join(d.rootPath, path)
	return os.Remove(fullPath)
}

// FileExists 检查文件是否存在
func (d *LocalDriver) FileExists(path string) (bool, error) {
	fullPath := filepath.Join(d.rootPath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

// Size 获取文件大小
func (d *LocalDriver) Size(path string) (int64, error) {
	fullPath := filepath.Join(d.rootPath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Copy 复制文件
func (d *LocalDriver) Copy(src, dst string) error {
	srcPath := filepath.Join(d.rootPath, src)
	dstPath := filepath.Join(d.rootPath, dst)

	// 确保目标目录存在
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// 打开源文件
	sourceFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制内容
	_, err = io.Copy(destFile, sourceFile)
	return err
}

// Move 移动文件
func (d *LocalDriver) Move(src, dst string) error {
	srcPath := filepath.Join(d.rootPath, src)
	dstPath := filepath.Join(d.rootPath, dst)

	// 确保目标目录存在
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	return os.Rename(srcPath, dstPath)
}

// List 列出目录内容
func (d *LocalDriver) List(path string) ([]FileInfo, error) {
	fullPath := filepath.Join(d.rootPath, path)

	dir, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		files = append(files, FileInfo{
			Name:      entry.Name(),
			Size:      entry.Size(),
			IsDir:     entry.IsDir(),
			UpdatedAt: entry.ModTime().Unix(),
		})
	}

	return files, nil
}
