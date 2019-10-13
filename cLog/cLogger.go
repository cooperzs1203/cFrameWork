/**
* @Author: Cooper
* @Date: 2019/10/13 19:58
 */

package cLog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Trace = 1 << iota
	Debug
	System
	Error
	Total
)

var logLevels = map[int]string{
	Trace 	: 	"TRACE",
	Debug 	: 	"DEBUG",
	System 	: 	"SYSTEM",
	Error 	: 	"ERROR",
	Total	:	"TOTAL",
}

type cLogger struct {
	logLevel int

	fileHandler *os.File
	filePath string
	fileDate string
	fileDateCount int
	fileMaxSize int64
	maxReserveDay int

	sChan chan []byte
	buffer []byte
}

func (cl *cLogger) initializes() {
	preFix := fmt.Sprintf("%s_%s_" , logLevels[cl.logLevel] , cl.fileDate)
	files := findDesignatedPrefixFileFromDir(preFix)
	cl.fileDateCount = len(files)

	if overSize := cl.checkFileOverSize(); overSize {
		cl.fileDateCount++ // if over max size , we start from next
	}

	err := cl.createAndUseLogFile()
	if err != nil {
		return
	}

	cl.monitorLogInfoAndFlushToBuffer()
	cl.monitorDateChange()
}

func (cl *cLogger) getCurrentFileName() string {
	return fmt.Sprintf("%s_%s_%d.%s" , logLevels[cl.logLevel] , cl.fileDate , cl.fileDateCount , _LOG_FILE_SUFFIX_)
}

func (cl *cLogger) getCurrentFilePath() string {
	return filepath.Join(_LOGS_DIR_ , cl.getCurrentFileName())
}

func (cl *cLogger) checkFileOverSize() bool {
	fileInfo , err := os.Stat(cl.getCurrentFilePath())
	if err != nil { // when we get error from open current file path , we start from next file
		return true
	}

	log.Printf("%s  --- %d" , cl.getCurrentFileName() , fileInfo.Size())

	if fileInfo.Size() + int64(len(cl.buffer)) >= cl.fileMaxSize {
		return true
	}

	return false
}

func (cl *cLogger) createAndUseLogFile() error {
	file , err := os.OpenFile(cl.getCurrentFilePath() , os.O_CREATE|os.O_APPEND|os.O_WRONLY , 0755)
	if err != nil {
		return err
	}

	if cl.fileHandler != nil {
		_ = cl.fileHandler.Close()
	}

	cl.fileHandler = file
	return nil
}

func (cl *cLogger) monitorLogInfoAndFlushToBuffer() {
	go func() {
		for {
			s , ok := <- cl.sChan
			if !ok {
				break
			}
			cl.flushLogInfoToBuffer(s)
		}
	}()
}

func (cl *cLogger) flushLogInfoToBuffer(s []byte) {
	if s == nil || len(s) == 0 {
		return
	}

	s = []byte(string(s)+"\n")
	if len(s)+len(cl.buffer) > cap(cl.buffer) {
		cl.flushBufferToLogFile()
		cl.buffer = make([]byte , 0 , _MB_)
	}
	cl.buffer = append(cl.buffer , s...)
}

func (cl *cLogger) flushBufferToLogFile() {
	if overSize := cl.checkFileOverSize(); overSize {
		cl.fileDateCount++ // if over max size , we start from next

		err := cl.createAndUseLogFile()
		if err != nil {
			return
		}
	}

	log.Printf("Write to %s -- %d -- %d" , cl.getCurrentFileName()  , len(cl.buffer) , cap(cl.buffer))

	_ , _ = cl.fileHandler.Write(cl.buffer)
}

func (cl *cLogger) monitorDateChange() {
	go func() {
		for {
			nowTime := time.Now()
			todayEndTime := nowTime.Add(time.Duration(24) * time.Hour).Format("2006-01-02") + " 00:00:01" // reserved 1 second space
			loc , _ := time.LoadLocation("Local")
			te , err := time.ParseInLocation("2006-01-02 15:04:05" , todayEndTime , loc)
			if err != nil {
				continue
			}
			countSec := te.Sub(nowTime).Seconds()

			<- time.After(time.Second * time.Duration(countSec))

			cl.fileDate = getNowDateFormatString()
			cl.fileDateCount = 1

			err = cl.createAndUseLogFile()
			if err != nil {
				defer cl.close()
				break
			}

			cl.deleteOverMaxReserveDayFiles()
		}
	}()
}

func (cl *cLogger) deleteOverMaxReserveDayFiles() {
	go func() {
		files := findDesignatedPrefixFileFromDir("")

		for _ , fileName := range files {
			fileDateStr := strings.Split(fileName , "_")[1]
			fileDate , err := time.Parse("20060102" , fileDateStr)
			if err != nil {
				continue
			}
			over , err := cl.checkDateOverMaxReserveDayOrNot(fileDate)
			if err != nil {
				continue
			}

			if over {
				filePath := filepath.Join(_LOGS_DIR_ , fileName)
				_ = os.Remove(filePath)
			}
		}

	}()
}

func (cl *cLogger) checkDateOverMaxReserveDayOrNot(date time.Time) (bool , error) {
	currentDate , err := time.Parse("20060102" , cl.fileDate)
	if err != nil {
		return false , err
	}

	gap := currentDate.Sub(date)
	over := gap.Hours()/24.0 > float64(cl.maxReserveDay)
	return over , nil
}

func (cl *cLogger) close() {
	_ = cl.fileHandler.Close()
	close(cl.sChan)
	cl.buffer = nil
}
