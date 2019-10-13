/**
* @Author: Cooper
* @Date: 2019/10/13 19:58
 */

package cLog

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"time"
)

func getGPID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

func getCallerInfo(level int) (string , int , bool) {
	_ , file , line , ok := runtime.Caller(level)
	if !ok {
		return "" , -1 , ok
	}
	_ , fileName := path.Split(file)
	return fileName , line , true
}

func baseLog(model string , format string , infos...interface{}) string {
	logInfo := fmt.Sprintf("[%s]" , time.Now().Format("2006-01-02 15:04:05.999"))
	if function , line , ok := getCallerInfo(3); ok {
		logInfo = fmt.Sprintf("%s|[%s:%d]" , logInfo , function , line)
	}
	logInfo = fmt.Sprintf("%s|[%s]|[%s]|" , logInfo , getGPID() , model)
	infosLog := fmt.Sprintf(format , infos...)
	logInfo = fmt.Sprintf("%s%s" , logInfo , infosLog)
	return logInfo
}

func LogTrace(format string, infos ...interface{}) {
	s := baseLog(logLevels[Trace] , format , infos...)
	logs[Trace].sChan <- []byte(s)
	if logSumUp && logs[Total] != nil {
		logs[Total].sChan <- []byte(s)
	}
}

func LogDebug(format string, infos ...interface{}) {
	s := baseLog(logLevels[Debug] , format , infos...)
	logs[Debug].sChan <- []byte(s)
	if logSumUp && logs[Total] != nil {
		logs[Total].sChan <- []byte(s)
	}
}

func LogSystem(format string, infos ...interface{}) {
	s := baseLog(logLevels[System] , format , infos...)
	logs[System].sChan <- []byte(s)
	if logSumUp && logs[Total] != nil {
		logs[Total].sChan <- []byte(s)
	}
}

func LogError(format string, infos ...interface{}) {
	s := baseLog(logLevels[Error] , format , infos...)
	logs[Error].sChan <- []byte(s)
	if logSumUp && logs[Total] != nil {
		logs[Total].sChan <- []byte(s)
	}
}
