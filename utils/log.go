package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/YouChenJun/hscan/libs"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = logrus.New()

func InitLog(cfg *libs.Cfg) {
	mwr := io.MultiWriter(os.Stdout)
	logDir := libs.LOGDIR
	if cfg.Env.LogFile == "" {
		if !FolderExists(logDir) {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create log dir:%v \n", logDir)
				os.Exit(1)
			}
		}
		tmpFile, err := os.CreateTemp(logDir, "hscan-*.log")
		if err == nil {
			cfg.Env.LogFile = tmpFile.Name()
		} else {
			tmpFile, _ := os.CreateTemp("/tmp/", "hscan-*.log")
			cfg.Env.LogFile = tmpFile.Name()
		}
	}

	logDir = filepath.Dir(cfg.Env.LogFile)
	if !FolderExists(logDir) {
		if err := os.MkdirAll(logDir, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create log dir: %v\n", logDir)
			os.Exit(1)
		}
	}

	f, err := os.OpenFile(cfg.Env.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open log file: %v\n", cfg.Env.LogFile)
		fmt.Fprintf(os.Stderr, "You can use the root user to run this program\n")
		os.Exit(1)
	} else {
		mwr = io.MultiWriter(os.Stdout, f)
	}
	logger = &logrus.Logger{
		Out:   mwr,
		Level: logrus.InfoLevel,
		Hooks: make(logrus.LevelHooks),
		Formatter: &prefixed.TextFormatter{
			ForceColors:     true,
			ForceFormatting: true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}

	if cfg.Mics.Debug == true {
		logger.SetLevel(logrus.DebugLevel)
	}
}

// PrintLine print seperate line
func PrintLine() {
	dash := color.HiWhiteString("-")
	fmt.Println(strings.Repeat(dash, 40))
}

// GoodF print good message
func GoodF(format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiBlueString("▶▶"))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// PrefixF print good message
func PrefixF(symbol string, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiGreenString(symbol))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// BannerF print info message
func BannerF(format string, data string) {
	banner := fmt.Sprintf("%v%v%v ", color.WhiteString("["), color.BlueString(format), color.WhiteString("]"))
	fmt.Printf("%v%v\n", banner, color.HiGreenString(data))
}

// BlockF print info message
func BlockF(name string, data string) {
	prefix := fmt.Sprintf("%v ", color.HiGreenString("💬 %v ", name))
	message := fmt.Sprintf("%v%v", prefix, data)
	logger.Info(message)
}

// TSPrintF print info message
func TSPrintF(format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiBlueString(" ▶ "))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// BadBlockF print info message
func BadBlockF(format string, args ...interface{}) {
	prefix := fmt.Sprintf("%v ", color.HiRedString(" [!] "))
	message := fmt.Sprintf("%v%v", prefix, fmt.Sprintf(format, args...))
	logger.Info(message)
}

// InforF print info message
func InforF(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

// Infor print info message
func Infor(args ...interface{}) {
	logger.Info(args...)
}

// ErrorF print good message
func ErrorF(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

// Error print good message
func Error(args ...interface{}) {
	logger.Error(args...)
}

// WarnF print good message
func WarnF(format string, args ...interface{}) {
	logger.Warning(fmt.Sprintf(format, args...))
}

// Warn print good message
func Warn(args ...interface{}) {
	logger.Warning(args...)
}

// TraceF print good message
func TraceF(format string, args ...interface{}) {
	logger.Trace(fmt.Sprintf(format, args...))
}

// Trace print good message
func Trace(args ...interface{}) {
	logger.Trace(args...)
}

// DebugF print debug message
func DebugF(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

// Debug print debug message
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
