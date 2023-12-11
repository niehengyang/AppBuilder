package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// 判断路径是否在指定目录下
func IsPathInDirectory(path string, directory string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		return false
	}

	return filepath.HasPrefix(absPath, absDirectory)
}

// 清空目录
func ClearDirectory(directoryPath string, delRoot bool) error {
	// Get a list of all entries (files and directories) in the directory
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	// Iterate through the entries and remove them
	for _, entry := range entries {
		entryPath := filepath.Join(directoryPath, entry.Name())

		// Remove files
		if entry.IsDir() {
			// Recursively clear subdirectories
			if err := ClearDirectory(entryPath, false); err != nil {
				return err
			}
		} else {
			// Remove the file
			if err := os.Remove(entryPath); err != nil {
				return err
			}
		}
	}

	// Remove the directory itself
	if delRoot {
		if err := os.Remove(directoryPath); err != nil {
			return err
		}
	}

	return nil
}

// 删除文件
func DeleteFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，尝试删除
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("无法删除文件: %v\n", err)
			return err
		}

	} else if os.IsNotExist(err) {
		fmt.Printf("文件不存在: %s\n", filePath)
		return nil
	} else {
		fmt.Printf("无法确定文件是否存在: %v\n", err)
		return err
	}
	return nil
}
