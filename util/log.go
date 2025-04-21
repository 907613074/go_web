package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var _dir string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_dir = strings.ReplaceAll(dir, "\\", "/")
	Log("Current directory:", _dir)
}

// Log 打印带有文件名、行号、方法名和内容的日志
func Log(v ...interface{}) {
	pc, file, line, ok := runtime.Caller(1) // 获取调用者的信息，2 表示获取调用 Log 函数的上一级
	if !ok {
		fmt.Println("无法获取调用者信息")
		return
	}

	// 获取方法名
	funcName := runtime.FuncForPC(pc).Name()
	// 提取方法名（去掉包路径）
	funcName = strings.Split(funcName, ".")[len(strings.Split(funcName, "."))-1]
	funcName = fmt.Sprintf("%-17s", funcName)

	fmt.Printf("\n%s %s:%d \n[%s] ==>  %v\n", time.Now().Format("2006-01-02 15:04:05"), file, line, funcName, v)
}

func Separator() {
	fmt.Println("--------------------------------------------------------------------------------------------")
}
