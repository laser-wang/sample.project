//一个可用于生产简单的http服务器
//1.带文件日志
//2.支持平滑重启和安全退出
//
//后续待追加
//1.鉴权处理
//2.https支持

package main

import (
	"sample.project/rest"
	_ "sample.project/rest"
)

//主函数
func main() {
	rest.StartServer()
}
