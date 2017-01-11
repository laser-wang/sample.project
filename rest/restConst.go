// Package rest 包含本http服务器例子所需要用到的函数和常量.
package rest

// 普通常量
const (
	// checkErr 函数用到的常量 1表示记录日志并退出 2表示仅记录日志，程序不退出
	checkFlagExit = 1
	checkFlagLog  = 2

	// 分隔符
	splitChar = "_"

	// token有效时长
	tokenValidTime = 30
)
