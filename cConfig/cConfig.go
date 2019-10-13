/**
* @Author: Cooper
* @Date: 2019/10/13 19:59
 */

package cConfig

import (
	"cFrameWork/cLog"
	"github.com/BurntSushi/toml-master"
	"net"
)

func TomlDecodeFile(filePath string , data interface{}) error {
	_ , err := toml.DecodeFile(filePath , data)
	if err != nil {
		return err
	}
	return nil
}

func CFrameWorkConfigFileLoad(filePath string) (CFWConfig , error) {
	var cfwc CFWConfig
	err := TomlDecodeFile(filePath , &cfwc)
	return cfwc , err
}

type CFWConfig struct {
	Mod string
	LogConfig CLogConfig
	NatsConfig CNatsConfig
}

type CLogConfig struct {
	LogSumUp bool
	LogConfigs []cLog.LogConfig
}

type CNatsConfig struct {
	Host string
	Port string
}

func (cnc *CNatsConfig) GetAddress() string {
	return net.JoinHostPort(cnc.Host , cnc.Port)
}