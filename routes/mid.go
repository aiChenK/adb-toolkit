package routes

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServeFileSystem 定义提供文件服务的抽象接口
type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

// EmbedFolder 将内嵌静态文件系统包装为 ServeFileSystem
func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

// Serve 返回一个中间件处理函数，用于提供静态文件服务
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
