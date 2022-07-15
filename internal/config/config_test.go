package config

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type Config_test struct {
	SrcDir         string `yaml:"src_dir"`
	DstDir         string `yaml:"dst_dir"`
	MaxCopyThreads int    `yaml:"max_copy_threads"`
	RescanInterval int    `yaml:"rescaninterval"`
	LogLevel       string `yaml:"loglevel"`
}

func Test_readConfig(t *testing.T) {
	req := require.New(t)
	config_test := Config_test{
		MaxCopyThreads: 4,
		RescanInterval: 3,
		LogLevel:       "INFO",
		SrcDir:         "/root/test_dir_for_project",
		DstDir:         "/root/test_dir_for_project_copy",
	}
	yamlConfig, err := yaml.Marshal(&config_test)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	fileName := "test_config.yaml"
	err = ioutil.WriteFile(fileName, yamlConfig, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	t.Run("simple test", func(t *testing.T) {
		res, err := readConfig("test_config.yaml")
		req.NotNil(res)
		req.Nil(err)
		req.Equal(config_test.DstDir, res.DstDir)
		req.Equal(config_test.SrcDir, res.SrcDir)
		req.Equal(config_test.LogLevel, res.LogLevel)
		req.Equal(config_test.RescanInterval, res.RescanInterval)
		req.Equal(config_test.MaxCopyThreads, res.MaxCopyThreads)

	})
}

func TestReloadConfig(t *testing.T) {
	req := require.New(t)
	config_test := Config_test{
		MaxCopyThreads: 4,
		RescanInterval: 3,
		LogLevel:       "INFO",
		SrcDir:         "/root/test_dir_for_project",
		DstDir:         "/root/test_dir_for_project_copy",
	}
	yamlConfig, err := yaml.Marshal(&config_test)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	fileName := "test_config.yaml"
	err = ioutil.WriteFile(fileName, yamlConfig, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	t.Run("simple test", func(t *testing.T) {
		res, err := readConfig("test_config.yaml")
		req.NotNil(res)
		req.Nil(err)
		req.Equal(config_test.DstDir, res.DstDir)
		req.Equal(config_test.SrcDir, res.SrcDir)
		req.Equal(config_test.LogLevel, res.LogLevel)
		req.Equal(config_test.RescanInterval, res.RescanInterval)
		req.Equal(config_test.MaxCopyThreads, res.MaxCopyThreads)

	})
}
