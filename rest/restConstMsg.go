package rest

// 消息常量
const (
	msgLoginFailed       = "登陆失败，请输入正确的用户code和用户密码。"
	msgLoginSuccess      = "登陆成功。"
	msgLoginExpired      = "登陆超时，请重新登陆。"
	msgLoginTokenInvalid = "令牌无效，请登陆后再操作。"
)

// 返回码
const (
	retCodeSuccess           = 200
	retCodeFailed            = 500
	retCodeLoginExpired      = 501
	retCodeLoginTokenInvalid = 502
)
