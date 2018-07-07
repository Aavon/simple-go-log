package log

import (
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"os/exec"
)

type FileLoger struct {
	logger   *log.Logger
	fileName string
	logTime  int64
	fileFd   *os.File
}


func Config(path string,level int) {
	fileLogger.fileName = path
	fileLogger.logger = log.New(fileLogger,NonePrefix,globalFlags)
}

func SetPath(path string) {
	if checkFileLogger() {
		fileLogger.fileName = path
		// restart logger
		fileLogger.logTime = 0
	}
}

func checkFileLogger() bool {
	return fileLogger.logger != nil
}

func FDebugf(format string, args ...interface{}) {
    if checkFileLogger() && globalLevel >= DebugLevel {
        fileLogger.logger.SetPrefix(DebugPrefix)
        fileLogger.logger.Output(2, fmt.Sprintf(format, args...))
    }
}

func FInfof(format string, args ...interface{}) {
    if checkFileLogger() && globalLevel >= InfoLevel {
        fileLogger.logger.SetPrefix(InfoPrefix)
        fileLogger.logger.Output(2, fmt.Sprintf(format, args...))
    }
}

func FWarnf(format string, args ...interface{}) {
    if checkFileLogger() && globalLevel >= WarnLevel {
        fileLogger.logger.SetPrefix(WarnPrefix)
        fileLogger.logger.Output(2, fmt.Sprintf(format, args...))
    }
}

func FErrorf(format string, args ...interface{}) {
    if checkFileLogger() && globalLevel >= ErrorLevel {
        fileLogger.logger.SetPrefix(ErrorPrefix)
        fileLogger.logger.Output(2, fmt.Sprintf(format, args...))
    }
}

func FFatalf(format string, args ...interface{}) {
    if checkFileLogger() && globalLevel >= FatalLevel {
        fileLogger.logger.SetPrefix(FatalPrefix)
        fileLogger.logger.Output(2, fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}

func (me FileLoger) Write(buf []byte) (n int, err error) {
    if fileLogger.fileName == "" {
        return len(buf), nil
    }

	// 1 day compress
    if fileLogger.logTime+3600 * 24 < time.Now().Unix() {
        fileLogger.createLogFile()
        fileLogger.logTime = time.Now().Unix()
    }

    if fileLogger.fileFd == nil {
        return len(buf), nil
    }

    return fileLogger.fileFd.Write(buf)
}

func (me *FileLoger) createLogFile() {
    logdir := "./"
    if index := strings.LastIndex(me.fileName, "/"); index != -1 {
		logdir = me.fileName[0:index] + "/"
		_,err := os.Stat(logdir)
		if !os.IsExist(err) {
			err := os.MkdirAll(logdir, os.ModePerm)
			if err != nil {
				fmt.Println("create log folder failed:",err)
			}
		}
    }

    now := time.Now()
    filename := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d", me.fileName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
    if err := os.Rename(me.fileName, filename); err == nil {
        go func() {
            tarCmd := exec.Command("tar", "-zcf", filename+".tar.gz", filename, "--remove-files")
            tarCmd.Run()

            rmCmd := exec.Command("/bin/sh", "-c", "find "+logdir+` -type f -mtime +2 -exec rm {} \;`)
            rmCmd.Run()
        }()
    }

	// three times try
    for index := 0; index < 3; index++ {
        if fd, err := os.OpenFile(me.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
            me.fileFd.Sync()
            me.fileFd.Close()
			me.fileFd = fd
            break
        } else {
			me.fileFd = nil
			fmt.Println("write log failed:",err)
		}
		time.Sleep(time.Second)
    }
}