/* Package xylog实现了业务debug日志接口
   debug 主要用于在开发过程中调试业务功能，提供各种级别的日志输出。
   authored by 田贵茂 2015.08.13
*/
package xylog

import (
	"fmt"
	"os"

	logs "github.com/astaxie/beego/logs"
)

//日志级别枚举值
type LogLevel int

const (
	FatalLevel   LogLevel = logs.LevelCritical
	ErrorLevel   LogLevel = logs.LevelError
	WarningLevel LogLevel = logs.LevelWarn
	InfoLevel    LogLevel = logs.LevelInfo
	DebugLevel   LogLevel = logs.LevelDebug
	TraceLevel   LogLevel = logs.LevelTrace
)

func (l LogLevel) String() (str string) {
	switch l {
	case logs.LevelTrace:
		str = fmt.Sprintf("Trace (%d)", int(l))
	case logs.LevelInfo:
		str = fmt.Sprintf("Info (%d)", int(l))
	case logs.LevelDebug:
		str = fmt.Sprintf("Debug (%d)", int(l))
	case logs.LevelWarn:
		str = fmt.Sprintf("Warning (%d)", int(l))
	case logs.LevelError:
		str = fmt.Sprintf("Error (%d)", int(l))
	case logs.LevelCritical:
		str = fmt.Sprintf("Critical (%d)", int(l))
	default:
		str = fmt.Sprintf("Undefined (%d)", int(l))
	}
	return
}

//日志类型
type XYLogger struct {
	config *LoggerConfig
	beelog *logs.BeeLogger
}

var (
	def_log_level LogLevel = DebugLevel

	def *XYLogger = NewLogger(DefConfig, 1000)
)

func NewLogger(lc *LoggerConfig, chanlen int64) (l *XYLogger) {
	l = &XYLogger{
		config: lc,
		beelog: logs.NewLogger(chanlen),
	}

	//l.ApplyConfig(l.config)
	return
}

func (l *XYLogger) Logger() *logs.BeeLogger {
	return l.beelog
}

func (l *XYLogger) Config() *LoggerConfig {
	return l.config
}

func (l *XYLogger) ApplyConfig(lc *LoggerConfig) {
	if lc == nil {
		lc = DefConfig
	}
	l.beelog.SetLevel(int(l.config.Level))

	if l.config.Verbose {
		l.beelog.EnableFuncCallDepth(true)
		l.beelog.SetLogFuncCallDepth(4)
	} else {
		l.beelog.EnableFuncCallDepth(false)
		l.beelog.SetLogFuncCallDepth(0)
	}

	var strconfig string
	if l.config.Stdout {
		//strconfig = fmt.Sprintf(`{"level":%v}`, l.config.Level)
		l.beelog.SetLogger("console", strconfig)
	} else {
		if l.config.Filename == "" {
			if l.config.NodeId >= 0 {
				l.config.LogId = l.config.NodeId
			} else {
				//l.config.Filename = fmt.Sprintf("%s.%d.log", l.config.AppName, os.Getpid())
				l.config.LogId = os.Getpid()
			}
			l.config.Filename = fmt.Sprintf("%s.%d.log", l.config.AppName, l.config.LogId)
		}

		strconfig = fmt.Sprintf(`{"filename":"%s/%s","maxlines":%v,"maxsize":%v,"daily":%v,"maxdays":%v,"rotate":%v}`,
			l.config.Path,
			l.config.Filename,
			l.config.Maxlines,
			l.config.Maxsize,
			l.config.Daily,
			l.config.Maxdays,
			l.config.Rotate) //,
		//logconfig.Level)
		l.beelog.SetLogger("file", strconfig)
	}
}

func (l *XYLogger) SetLogLevel(level LogLevel) {
	l.beelog.SetLevel((int)(level))
}

func (l *XYLogger) Log(log_level LogLevel, format string, v ...interface{}) {
	switch log_level {
	case logs.LevelCritical:
		l.beelog.Critical(format, v...)
	case logs.LevelDebug:
		l.beelog.Debug(format, v...)
	case logs.LevelError:
		l.beelog.Error(format, v...)
	case logs.LevelInfo:
		l.beelog.Info(format, v...)
	case logs.LevelTrace:
		l.beelog.Trace(format, v...)
	case logs.LevelWarn:
		l.beelog.Warn(format, v...)
	}
}

func ApplyConfig(lc *LoggerConfig) {
	def.ApplyConfig(lc)
}

//--------------- 不指定id的接口，是否记录依赖于全局的日志级别 begin ---------------
const NO_UID = ""

func TraceNoId(format string, v ...interface{}) {
	if def.config.Level > TraceLevel {
		return
	}

	def.Log(TraceLevel, format, v...)
}

func InfoNoId(format string, v ...interface{}) {
	if def.config.Level > InfoLevel {
		return
	}

	def.Log(InfoLevel, format, v...)
}

func DebugNoId(format string, v ...interface{}) {
	if def.config.Level > DebugLevel {
		return
	}

	def.Log(DebugLevel, format, v...)
}

func WarningNoId(format string, v ...interface{}) {
	if def.config.Level > WarningLevel {
		return
	}
	def.Log(WarningLevel, format, v...)
}

func ErrorNoId(format string, v ...interface{}) {

	if def.config.Level > ErrorLevel {
		return
	}

	def.Log(ErrorLevel, format, v...)

}

func FatalNoId(format string, v ...interface{}) {

	if def.config.Level > FatalLevel {
		return
	}

	def.Log(FatalLevel, format, v...)
}

//--------------- 不指定id的接口，是否记录依赖于全局的日志级别 end ---------------

//--------------- 指定id的接口，是否记录依赖于id和全局的日志级别 begin ---------------
func Trace(id interface{}, format string, v ...interface{}) {
	if IsNeedLog(id, TraceLevel) {
		format = fmt.Sprintf("[%v] %s", id, format)
		def.Log(TraceLevel, format, v...)
	}
}

func Info(id interface{}, format string, v ...interface{}) {
	if IsNeedLog(id, InfoLevel) {
		format = fmt.Sprintf("[%v] %s", id, format)
		def.Log(InfoLevel, format, v...)
	}
}

func Debug(id interface{}, format string, v ...interface{}) {
	if IsNeedLog(id, DebugLevel) {
		format = fmt.Sprintf("[%v] %s", id, format)
		def.Log(DebugLevel, format, v...)
	}
}

func Warning(id interface{}, format string, v ...interface{}) {
	if IsNeedLog(id, WarningLevel) {
		format = fmt.Sprintf("[%v] %s", id, format)
		def.Log(WarningLevel, format, v...)
	}
}

func Error(id interface{}, format string, v ...interface{}) {
	if IsNeedLog(id, ErrorLevel) {
		format = fmt.Sprintf("[%v] %s", id, format)
		def.Log(ErrorLevel, format, v...)
	}
}

func Fatal(id interface{}, format string, v ...interface{}) {
	if IsNeedLog(id, FatalLevel) {
		format = fmt.Sprintf("[%v] %s", id, format)
		def.Log(FatalLevel, format, v...)
	}
}

//--------------- 指定id的接口，是否记录依赖于id和全局的日志级别 end ---------------

//设置全局的日志级别
func SetLogLevel(level LogLevel) {
	def.SetLogLevel(level)
}

func IsNeedLog(id interface{}, level LogLevel) bool {

	if !(DefIdManager.IsIdExist(id)) {
		if def.config.Level > InfoLevel {
			return false
		}
	}

	return true
}

func Close() {
	def.beelog.Close()
}
