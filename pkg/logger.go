package logger

import (
	"time"

	"github.com/nicexiaonie/glog"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func LoggerInit(path string) {

	Logger, _ = glog.New(&glog.Config{
		// 日志存放目录
		Path: path,
		// 日志文件名
		Filename: "app.log",
		// 日志级别 超过并包含设置级别以上的日志会处理保存
		Level: "info",
		// json text formatter
		Format: "json",
		// 自定义格式化接口
		//Formatter:    logrus.Formatter,

		// option: file
		Output:       "file",
		ReportCaller: true,

		// 可以为空，非空按照日期格式生成日志文件
		Split: ".2006010215",
		// 日志最后修改时间超过多少时间进行清理
		Lifetime: 5 * time.Second,
	})
}
