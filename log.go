package log

import (
	"os"
	"log"
)

const (
	LOG_TO_FILE = true
	DEF_LOG_FILE = "./log/output.log"
)

const (
    PanicLevel int = iota
    FatalLevel
    ErrorLevel
    WarnLevel
    InfoLevel
    DebugLevel
)

const (
    PanicPrefix = "[Panic]"
    FatalPrefix = "[Fatal]"
    ErrorPrefix = "[Error]"
    WarnPrefix  = "[Warn]"
    InfoPrefix  = "[Info]"
	DebugPrefix = "[Debug]"
	NonePrefix  = ""
)

var (
	// the local logger
	logger *log.Logger

	// log to file
	fileLogger FileLoger

	// log level
	globalLevel int

	// log file path
	globalFile  string

	// lof flags 
	globalFlags int
)


// config
func init(){
	// default [debug]
	globalLevel = DebugLevel
	// log file path 
	globalFile  = DEF_LOG_FILE
	// flags
	globalFlags = log.Lmicroseconds | log.Lshortfile

	logger = log.New(os.Stdout,NonePrefix,globalFlags)

	if LOG_TO_FILE {
		Config(DEF_LOG_FILE,DebugLevel)
	}
}

// Log makes use of github.com/go-log/log.Log
func Log(v ...interface{}) {
	if checkLogger() {
		logger.SetPrefix(NonePrefix)
		logger.Println(v...)
	}
}

func Logf(format string, v ...interface{}) {
	if checkLogger() {
		logger.SetPrefix(NonePrefix)
		logger.Printf(format + "\n", v...)
	}
}

//Log info with log title 
func LogI(v ...interface{}) {
	if checkLogger() && globalLevel >= InfoLevel {
		logger.SetPrefix(InfoPrefix)
		title := v[0]
		logger.Println(title,":", v[1:])
	}
}

//Log debug with log title 
func LogD(v ...interface{}) {
	if checkLogger() && globalLevel >= DebugLevel {
		logger.SetPrefix(DebugPrefix)
		title := v[0]
		logger.Println(title,":", v[1:])
	}
}

//Log error with log title 
func LogE(v ...interface{}) {
	if checkLogger() && globalLevel >= ErrorLevel {
		logger.SetPrefix(ErrorPrefix)
		title := v[0]
		logger.Println(title,":", v[1:])
	}
}

// Fatal logs with Log and then exits with os.Exit(1)
func Fatal(v ...interface{}) {
	if checkLogger() && globalLevel >= FatalLevel {
		logger.SetPrefix(FatalPrefix)
		title := v[0]
		logger.Println(title,":", v[1:])
	}
	os.Exit(1)
}

// Fatalf logs with Logf and then exits with os.Exit(1)
func Fatalf(format string, v ...interface{}) {
	if checkLogger() && globalLevel >= FatalLevel {
		logger.SetPrefix(FatalPrefix)
		logger.Printf(format + "\n",v)
	}
	os.Exit(1)
}

func checkLogger() bool {
	return logger != nil
}

// SetLogger sets the local logger
func SetLogger(l *log.Logger) {
	logger = l
}

func SetLevel(level int) {
	globalLevel = level
}

func SetFlags(flags int) {
	globalFlags = flags
	if logger != nil {
		logger.SetFlags(flags)
	}
	if fileLogger.logger != nil {
		fileLogger.logger.SetFlags(flags)
	} 
}
