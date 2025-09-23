package utils

import "os"

// Exists 判断所给路径的文件或文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// Ext 返回路径所指向文件的扩展名
// 扩展名是指路径中最后一个斜杠分隔的元素里，从最后一个点开始的后缀部分
// 如果没有点，则返回空字符串
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}