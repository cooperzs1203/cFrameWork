package cLogHelper

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

const (
	Trace   = 1 << iota
	Debug   = 1 << iota
	System  = 1 << iota
	Error   = 1 << iota
)

var logPrefix = map[int]string {
	Trace 	: "Trace",
	Debug 	: "Debug",
	System  : "System",
	Error 	: "Error",
}

func InitLogSystem(msgChanSpace int64) {
	err := initLogger(msgChanSpace)
	// if init logger error , stop
	if err != nil {
		panic(err)
	}
	go defaultLogger.printInfo()
	go defaultLogger.monitorTime()
}

func GetGPID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

func GetCallerInfo(level int) (string , int , bool) {
	pc , _ , line , ok := runtime.Caller(level)
	if !ok {
		return "" , -1 , ok
	}
	callFun := runtime.FuncForPC(pc)
	return callFun.Name() , line , true
}

func GetLogTime() string {
	return time.Now().Format("2006-01-02 15:04:05.999")
}

func baseLog(model string , format string , infos...interface{}) {
	logInfo := fmt.Sprintf("[%s]" , GetLogTime())
	if function , line , ok := GetCallerInfo(4); ok {
		logInfo = fmt.Sprintf("%s|[%s:%d]" , logInfo , function , line)
	}
	logInfo = fmt.Sprintf("%s|[%s]|[%s]|" , logInfo , GetGPID() , model)
	infosLog := fmt.Sprintf(format , infos...)
	logInfo = fmt.Sprintf("%s%s" , logInfo , infosLog)
	defaultLogger.infoChan <- logInfo
}

func midLog(model string, format string, infos ...interface{}) {
	baseLog(model , format , infos...)
}

func LogDebug(format string, infos ...interface{}) {
	midLog(logPrefix[Debug] , format , infos...)
}

func LogTrace(format string, infos ...interface{}) {
	midLog(logPrefix[Trace] , format , infos...)
}

func LogSystem(format string, infos ...interface{}) {
	midLog(logPrefix[System] , format , infos...)
}

func LogError(format string, infos ...interface{}) {
	midLog(logPrefix[Error] , format , infos...)
}
