package uflog

import (
	"fmt"
	"os"
)

const (
	defaultLogDir         = "log"
	defaultLogPrefix      = "ufile_"
	defaultLogSuffix      = ".log"
	defaultLogSize        = 50 // MB
	defaultLogLevelString = "DEBUG"
)

var Glogger *Logger

func InitLogger(dir, prefix, suffix string, size int64, level string) {
	if Glogger != nil {
		Glogger.Close()
	}
	if dir == "" {
		dir = defaultLogDir
	}
	if prefix == "" {
		prefix = defaultLogPrefix
	}
	if suffix == "" {
		suffix = defaultLogSuffix
	}
	if size <= 0 {
		size = defaultLogSize
	}
	if level == "" {
		level = defaultLogLevelString
	}
	logger, err := NewRotate(dir, prefix, suffix, size)
	if err != nil {
		fmt.Println("Init Logger fail:", err)
		os.Exit(-1)
	}
	Glogger = logger
	SetLogLevel(level)
}

func InitDefaultLogger() {
	InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString)
}

//设置日志中输出代码文件名的层次
func SetCodeCallLevel(call_level int) {
	Glogger.CallLevel = call_level
}

func SetLogLevel(level string) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, level)
	}
	Glogger.SetOutputLevelString(level)
}

func SetDailyRotate(daily bool) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString)
	}
	Glogger.SetDailyRotate(daily)
}

func INFOF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Infof(format, v...)
}

func INFO(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}

	Glogger.Info(v...)
}

func ERRORF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Errorf(format, v...)
}

func ERROR(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Error(v...)
}

func WARN(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warn(v...)
}

func WARNF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warnf(format, v...)
}

func DEBUG(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debug(v...)
}

func DEBUGF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debugf(format, v...)
}

//下面的日志接口，可以把session seq 输出，方便跟踪日志
//add by delex 20180321
//===========================================================
func INFOF2(seq_id string, format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Infof2(seq_id, format, v...)
}

func INFO2(seq_id string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Info2(seq_id, v...)
}

func ERRORF2(seq_id string, format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Errorf2(seq_id, format, v...)
}

func ERROR2(seq_id string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Error2(seq_id, v...)
}

func WARN2(seq_id string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warn2(seq_id, v...)
}

func WARNF2(seq_id string, format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warnf2(seq_id, format, v...)
}

func DEBUG2(seq_id string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debug2(seq_id, v...)
}

func DEBUGF2(seq_id string, format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debugf2(seq_id, format, v...)
}
