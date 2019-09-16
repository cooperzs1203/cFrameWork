package cFrameWork

import (
	"cFrameWork/cFrameWork/cConfigHelper"
	"cFrameWork/cFrameWork/cLogHelper"
	"cFrameWork/cFrameWork/cNatsHelper"
	"errors"
)

type cFrameWork struct {
	Mod string
}

var defaultCFrameWork = &cFrameWork{}


func InitServe(mod string) error {
	if mod == "" {
		return errors.New("mod can not be empty")
	}

	defaultCFrameWork.Mod = mod
	return nil
}

func InitNatsHelper(configFilePath string) error {
	if configFilePath == "" {
		return errors.New("configuration file path can not be empty")
	}

	// load cNatsHelper configuration file
	err := cConfigHelper.LoadNatsConfig(configFilePath)
	if err != nil {
		return err
	}

	// load cNatsHelper
	err = cNatsHelper.InitAndLoadHelper(defaultCFrameWork.Mod , cConfigHelper.GetNatsConfig().GetAddr())
	if err != nil {
		return err
	}

	return nil
}

func InitLogHelper(configFilePath string) error {
	if configFilePath == "" {
		return errors.New("configuration file path can not be empty")
	}

	// load cLogHelper configuration file
	err := cConfigHelper.LoadLogConfig(configFilePath)
	if err != nil {
		return err
	}

	// load cLogHelper
	cLogHelper.InitLogSystem(cConfigHelper.GetLogConfig().MsgChanSpace)

	return nil
}