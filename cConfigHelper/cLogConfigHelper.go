package cConfigHelper

import (
	"github.com/BurntSushi/toml"
)

type CLogConfig struct {
	MsgChanSpace int64
}

var clc *CLogConfig

func LoadLogConfig(path string) error {
	var logConfig CLogConfig
	_ , err := toml.DecodeFile(path , &logConfig)
	if err != nil {
		return err
	}
	clc = &logConfig
	return nil
}

func GetLogConfig() *CLogConfig {
	return clc
}
