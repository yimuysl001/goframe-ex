/**
兼容 fs 风格前端数据，todo 目前未测试过待完善
*/

package eswagger

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

func TarGzipEmbedFS(fsys fs.FS, root string) ([]byte, error) {
	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)
	defer g.Try(context.TODO(), func(ctx context.Context) {
		tw.Close()
		gw.Close()
	})

	// 遍历文件系统
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath := strings.TrimPrefix(path, root+"/")
		if relPath == path { // 处理根目录情况
			relPath = filepath.Base(path)
		}

		// 获取文件信息
		info, err := d.Info()
		if err != nil {
			return err
		}

		// 创建tar头
		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 写入tar头
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// 如果是普通文件，写入内容
		if !d.IsDir() {
			file, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk embedded FS: %w", err)
	}

	return buf.Bytes(), nil
}

func InitPublic(ctx context.Context, staticFiles embed.FS, rootpaths ...string) bool {
	serverRoot := g.Cfg().MustGet(ctx, "server.serverRoot").String()
	if serverRoot == "" { // 未配置路径，不需要处理swagger
		g.Log().Error(ctx, "未开启serverRoot")
		return false
	}

	fspath := "static"

	if len(rootpaths) > 0 && rootpaths[0] != "" {
		fspath = rootpaths[0]
	}

	embedFS, err := TarGzipEmbedFS(staticFiles, fspath)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}

	if !strings.Contains(serverRoot, ":") { // 全路径不处理
		var pwd = gfile.Pwd()
		pwd = strings.ReplaceAll(pwd, "\\", "/")
		serverRoot = path.Join(pwd, serverRoot)
	}
	if !strings.HasSuffix(serverRoot, "/") {
		serverRoot = serverRoot + "/"
	}

	err = gres.Add(string(embedFS), serverRoot)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	g.Log().Debug(ctx, "静态资源添加完成 添加完成")
	return true

}
