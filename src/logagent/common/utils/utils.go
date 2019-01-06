package utils

import "github.com/astaxie/beego/logs"


// 在 app.conf 的Beego日志配置level是字符串
// Beego框架的日志级别是 int
// 在配置文件写int 可读性不强
func BeegoConvertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}
