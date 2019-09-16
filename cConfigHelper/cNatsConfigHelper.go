package cConfigHelper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"net"
)

type CNatsConfig struct {
	Host string
	Port int64
}

func (cnConfig *CNatsConfig) GetAddr() string {
	addr := net.JoinHostPort(cnConfig.Host , fmt.Sprintf("%d" , cnConfig.Port))
	return "nats://" + addr
}

var cnc *CNatsConfig

func LoadNatsConfig(path string) error {
	var natsConfig CNatsConfig
	_ , err := toml.DecodeFile(path , &natsConfig)
	if err != nil {
		return err
	}
	cnc = &natsConfig
	return nil
}

func GetNatsConfig() *CNatsConfig {
	return cnc
}
