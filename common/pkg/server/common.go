package server

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"
)

// 定义日志/堆栈打印中用到的占位符字节切片
var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

const (
	// Version ApiServer框架版本号
	Version = "v1.0.0"
	Module  = "discron/api-server"
)

func formatTime(t time.Time) string {
	var timeString = t.Format("2006/01/01 - 15:04:05")
	return timeString
}


// stack 生成当前程序调用栈的详细信息（包含文件路径、行号、函数名、源码内容）
// 作用：当程序发生panic时，用于打印详细的调用栈，帮助开发者定位问题
// 参数skip：需要跳过的栈帧数量（比如跳过当前stack函数本身，避免栈信息冗余）
// 返回值：包含完整调用栈信息的字节切片（可直接转为字符串打印）
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source 从文件行切片中获取指定行号的源码，并去除首尾空白
// 作用：在调用栈中展示具体出错行的代码，帮助快速定位问题代码
// 参数lines：文件内容按行分割后的字节切片（0索引）
// 参数n：调用栈中获取的行号（1索引，需转换为0索引）
// 返回值：指定行的源码字节切片（若行号无效，返回dunno占位符）
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function 通过程序计数器（PC）获取对应的函数名，并简化格式（去除路径和包名，保留核心函数名）
// 作用：在调用栈中展示简洁的函数名，避免冗长的全路径函数名（如将github.com/xxx/pkg·func简化为func）
// 参数pc：程序计数器（从runtime.Caller获取）
// 返回值：简化后的函数名字节切片（若无法获取函数名，返回dunno占位符）
// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}