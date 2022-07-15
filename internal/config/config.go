package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	configModtime  int64
	errNotModified = errors.New("not modified")
)

type Config struct {
	SrcDir         string `yaml:"src_dir"`
	DstDir         string `yaml:"dst_dir"`
	MaxCopyThreads int    `yaml:"max_copy_threads"`
	RescanInterval int    `yaml:"rescaninterval"`
	LogLevel       string `yaml:"loglevel"`
	SkipIfexists   string `yaml:"skipIfExists"`
}

func readConfig(ConfigName string) (x *Config, err error) {
	var file []byte
	if file, err = ioutil.ReadFile(ConfigName); err != nil {
		return nil, err
	}
	x = new(Config)
	if err = yaml.Unmarshal(file, x); err != nil {
		return nil, err
	}
	if x.LogLevel == "" {
		x.LogLevel = "Debug"
	}
	return x, nil
}

func ReloadConfig(configName string) (cfg *Config, err error) {
	info, err := os.Stat(configName)
	if err != nil {
		return nil, err
	}
	if configModtime != info.ModTime().UnixNano() {
		configModtime = info.ModTime().UnixNano()
		cfg, err = readConfig(configName)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return nil, errNotModified
}
