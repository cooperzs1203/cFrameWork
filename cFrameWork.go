/**
* @Author: Cooper
* @Date: 2019/10/13 19:53
 */

package cFrameWork

import (
	"cFrameWork/cConfig"
	"cFrameWork/cLog"
	"cFrameWork/cNats"
	"errors"
)

func LoadAndServe(configFilePath string) error {
	if configFilePath == "" {
		return errors.New("we need configuration file for cFrameWork")
	}

	config , err := cConfig.CFrameWorkConfigFileLoad(configFilePath)
	if err != nil {
		return err
	}

	err = loadLogServe(config.LogConfig)
	if err != nil {
		return err
	}

	err = loadNatsServe(config.Mod , config.NatsConfig.GetAddress())
	if err != nil {
		return err
	}

	return nil
}

func loadLogServe(config cConfig.CLogConfig) error {
	cLog.LoadLoggers(config.LogSumUp , config.LogConfigs)
	return nil
}

func loadNatsServe(mod , addr string) error {
	return cNats.InitAndLoadHelper(mod , addr)
}