/**
* @Author: Cooper
* @Date: 2019/10/13 19:58
 */

package cLog

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	_LOGS_DIR_ = "./logs"
	_LOG_FILE_SUFFIX_ = "txt"
	_B_  = 1
	_KB_ = 1024*_B_
	_MB_ = 1024*_KB_
	_GB_ = 1024*_MB_
	_TB_ = 1024*_GB_
	_CLOGGER_SCHAN_SPACE_ = 5000
)

var (
	logSumUp bool
	logs map[int]*cLogger
)

func init() {
	err := createLogsDirIfNotExists()
	if err != nil {
		panic(err)
	}
	logs = make(map[int]*cLogger)
}

func createLogsDirIfNotExists() error {
	fileInfo , err := os.Stat(_LOGS_DIR_)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) || !fileInfo.IsDir() {
		return os.Mkdir(_LOGS_DIR_ , os.ModePerm)
	}

	return nil
}

type LogConfig struct {
	LogLevel int
	LogMsgSpace int64
	LogBufferSpace int64
	FileMaxSize int64
	FileMaxReserveDays int
}

func LoadLoggers(lsu bool, configs []LogConfig) {
	logSumUp = lsu
	for _ , config := range configs {
		logger := &cLogger{
			logLevel:config.LogLevel,
			fileHandler:nil,
			filePath:_LOGS_DIR_,
			fileDate:getNowDateFormatString(),
			fileDateCount:1,
			fileMaxSize:config.FileMaxSize,
			maxReserveDay:config.FileMaxReserveDays,
			sChan:make(chan []byte , config.LogMsgSpace),
			buffer:make([]byte , 0 , config.LogBufferSpace),
		}
		logger.initializes()
		logs[logger.logLevel] = logger
	}
}

func getNowDateFormatString() string {
	return time.Now().Format("20060102")
}

func findDesignatedPrefixFileFromDir(filePrefix string) []string {
	var logFiles []string
	_ = filepath.Walk(_LOGS_DIR_ , func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name() , filePrefix) {
			logFiles = append(logFiles , info.Name())
		}
		return nil
	})

	return logFiles
}

