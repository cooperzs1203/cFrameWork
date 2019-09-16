package cLogHelper

import (
	"log"
	"os"
	"time"
)

const (
	_LOG_DIR_ = "./log"
	_FILE_SUFFIX_ = ".txt"
)

var defaultLogger *cLogger

type cLogger struct {
	logDir string
	todayFile string
	fileHandler *os.File
	infoChan chan string
}

func initLogger(msgChanSpace int64) error {
	log.SetFlags(0)
	log.SetPrefix("")

	defaultLogger = &cLogger{
		logDir:_LOG_DIR_,
		todayFile:"",
		fileHandler: nil,
		infoChan:    make(chan string , msgChanSpace),
	}
	err := defaultLogger.checkAndCreateDir()
	if err != nil { return err }

	return nil
}

func (cLog *cLogger) printInfo() {
	for {
		if cLog.fileHandler == nil { continue }
		info := <- cLog.infoChan
		log.Println(info)
	}
}

func (cLog *cLogger) monitorTime() {
	for {
		<- time.After(time.Duration(10) * time.Second)
		todayFile := cLog.logDir + "/" + time.Now().Format("20060102") + _FILE_SUFFIX_

		// if cLog.todayFile is equal to todayFile , that means still today
		if cLog.todayFile == todayFile {
			LogSystem("fileHandler don't need change")
			continue
		}

		err := cLog.checkAndCreateFile(todayFile)
		if err != nil {
			LogSystem("check file %s error : %s" , todayFile , err.Error())
			continue
		}

		todayFileHandler , err := os.OpenFile(todayFile , os.O_WRONLY|os.O_CREATE|os.O_APPEND , 0644)
		if err != nil {
			LogSystem("Open %s error : %s" , todayFile , err.Error())
			continue
		}

		// close the last file handler
		if cLog.fileHandler != nil {
			err = cLog.fileHandler.Close()
			if err != nil {
				LogSystem("cLog.fileHandler Close error : %s" , err.Error())
				continue
			}
		}

		cLog.fileHandler = todayFileHandler
		cLog.todayFile = todayFile
		log.SetOutput(cLog.fileHandler)
	}
}

// check file exist , if not than create it
func (cLog *cLogger) checkAndCreateFile(file string) error {
	if fileExist := cLog.fileExist(file); !fileExist {
		err := cLog.createFile(file)
		if err != nil {
			return err
		}
		LogSystem("Create %s success." , file)
	} else {
		LogSystem("%s exist." , file)
	}
	return nil
}

// check ./log dir exist , if not than create it
func (cLog *cLogger) checkAndCreateDir() error {
	if dirExist := cLog.dirExist(cLog.logDir); !dirExist {
		err := cLog.createDir(cLog.logDir)
		if err != nil {
			return err
		}
		LogSystem("Create %s success." , cLog.logDir)
	} else {
		LogSystem("%s exist." , cLog.logDir)
	}
	return nil
}

func (cLog *cLogger) createDir(path string) error {
	err := os.Mkdir(path , os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (cLog *cLogger) dirExist(path string) bool {
	fileInfo , err := os.Stat(path)
	if err == nil {
		return fileInfo.IsDir()
	}

	return false
}

func (cLog *cLogger) fileExist(file string) bool {
	fileInfo , err := os.Stat(file)
	if err == nil {
		return !fileInfo.IsDir()
	}

	return false
}

func (cLog *cLogger) createFile(file string) error {
	_ , err := os.Create(file)
	if err != nil {
		return err
	}

	return nil
}