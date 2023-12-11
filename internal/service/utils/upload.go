package utils

import (
	"appBuilder/ebyte/logger"
	"appBuilder/internal/svc"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadFileInfo struct {
	Path       string
	Name       string
	OriginName string
	Size       int64
}

// UploadImage 图片文件上传
func UploadImage(svcCtx *svc.ServiceContext, req *http.Request) (fileInfo UploadFileInfo, err error) {
	var info UploadFileInfo

	// 设置表单最大10MB
	errS := req.ParseMultipartForm(10 << 20)
	if errS != nil {
		logger.Error("文件大小超出范围:", zap.Error(errS))
		return fileInfo, errS
	}

	// 读取文件
	file, header, errF := req.FormFile("file")
	if errF != nil {
		logger.Error("文件读取失败:", zap.Error(errS))
		return fileInfo, errF
	}
	defer func(file multipart.File) {
		if errC := file.Close(); errC != nil {
			logger.Error("文件关闭失败:", zap.Error(errC))
		}
	}(file)

	// 获取上传的文件类型
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/png"} // 允许的图片类型
	// 判断文件类型是否允许
	var isValidType bool
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		logger.Error(fmt.Sprintf("不支持的图片类型"))
		return fileInfo, errors.New("不支持的图片类型")
	}

	// 为上传的文件生成一个唯一的文件名
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	uploadPath := svcCtx.Config.FileServer.ImageFileSavePath + "Upload/"

	// 如果没有则创建目录
	exists, _ := PathExists(uploadPath)
	if !exists {
		if errM := os.MkdirAll(uploadPath, 0777); errM != nil {
			return fileInfo, errM
		}
	}

	// 保存文件

	if _, errI := os.Stat(uploadPath); os.IsNotExist(errI) {
		if errMk := os.Mkdir(uploadPath, os.ModePerm); errMk != nil {
			return fileInfo, errMk
		}
	}

	out, errC := os.Create(filepath.Join(uploadPath, fileName))
	if errC != nil {
		logger.Error("文件创建失败:", zap.Error(errC))
		return fileInfo, errC
	}
	defer func(out *os.File) {
		if errOut := out.Close(); errOut != nil {
			logger.Error("文件关闭失败:", zap.Error(errOut))
		}
	}(out)

	// 将文件内容复制到输出文件
	_, errCopy := io.Copy(out, file)
	if errCopy != nil {
		logger.Error("文件复制失败:", zap.Error(errCopy))
		return fileInfo, errCopy
	}

	// 获取文件信息
	info.Path = filepath.Join(uploadPath, fileName)
	info.Name = fileName
	info.OriginName = header.Filename
	info.Size = header.Size

	return info, nil
}

// UploadPpt PPT文件上传
func UploadPpt(svcCtx *svc.ServiceContext, req *http.Request) (fileInfo UploadFileInfo, err error) {
	var info UploadFileInfo

	// 设置表单最大100MB
	errS := req.ParseMultipartForm(100 * 1024 * 1024)
	if errS != nil {
		logger.Error("文件大小超出范围:", zap.Error(errS))
		return fileInfo, errS
	}

	// 读取文件
	file, header, errF := req.FormFile("file")
	if errF != nil {
		logger.Error("文件读取失败:", zap.Error(errS))
		return fileInfo, errF
	}
	defer func(file multipart.File) {
		if errC := file.Close(); errC != nil {
			logger.Error("文件关闭失败:", zap.Error(errC))
		}
	}(file)

	// 检查文件是否具有有效的扩展名
	ext := filepath.Ext(header.Filename)
	if ext != ".ppt" && ext != ".pptx" {
		logger.Error(fmt.Sprintf("不支持的文件类型"))
		return fileInfo, errors.New("不支持的文件类型")
	}

	// 为上传的文件生成一个唯一的文件名
	fileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	uploadPath := svcCtx.Config.FileServer.PptFileSavePath + "Upload/"

	// 如果没有则创建目录
	exists, _ := PathExists(uploadPath)
	if !exists {
		if errM := os.MkdirAll(uploadPath, 0777); errM != nil {
			return fileInfo, errM
		}
	}

	// 保存文件
	if _, errI := os.Stat(uploadPath); os.IsNotExist(errI) {
		if errMk := os.Mkdir(uploadPath, os.ModePerm); errMk != nil {
			return fileInfo, errMk
		}
	}

	out, errC := os.Create(filepath.Join(uploadPath, fileName))
	if errC != nil {
		logger.Error("文件创建失败:", zap.Error(errC))
		return fileInfo, errC
	}
	defer func(out *os.File) {
		if errOut := out.Close(); errOut != nil {
			logger.Error("文件关闭失败:", zap.Error(errOut))
		}
	}(out)

	// 将文件内容复制到输出文件
	_, errCopy := io.Copy(out, file)
	if errCopy != nil {
		logger.Error("文件复制失败:", zap.Error(errCopy))
		return fileInfo, errCopy
	}

	// 获取文件信息
	info.Path = filepath.Join(uploadPath, fileName)
	info.Name = fileName
	info.OriginName = header.Filename
	info.Size = header.Size

	return info, nil
}

func DelFile(svcCtx *svc.ServiceContext, filePath string) error {

	// 检查文件路径是否在目标目录下
	uploadDir := svcCtx.Config.FileServer.RelativePath
	if !IsPathInDirectory(filePath, uploadDir) {
		logger.Error("无效的文件路径:")
		return errors.New("无效的文件路径")
	}

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err == nil {
		// 文件存在，删除文件
		err := os.Remove(filePath)
		if err != nil {
			logger.Error("无法删除文件:", zap.Error(err))
			return err
		}
	} else if os.IsNotExist(err) {
		// 文件不存在，不进行处理
		logger.Error("文件不存在，无需处理:", zap.Error(err))
		return nil
	} else {
		// 发生其他错误
		logger.Error("发生错误:", zap.Error(err))
		return err
	}

	return nil
}
