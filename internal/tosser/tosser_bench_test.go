package tosser

import (
	"final-project/internal/config"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

type Config_test struct {
	SrcDir         string `yaml:"src_dir"`
	DstDir         string `yaml:"dst_dir"`
	MaxCopyThreads int    `yaml:"max_copy_threads"`
	RescanInterval int    `yaml:"rescaninterval"`
	LogLevel       string `yaml:"loglevel"`
}

type processingItemTest struct {
	srcFile         string
	fullSrcFilePath string
	params          *config.Config
	size            int64
}

func BenchmarkProcessItem(b *testing.B) {
	config_test := Config_test{
		MaxCopyThreads: 4,
		RescanInterval: 3,
		LogLevel:       "INFO",
		SrcDir:         "/tmp",
		DstDir:         "/tmp/test_copy",
	}
	yamlConfig, err := yaml.Marshal(&config_test)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	fileName := "test_config.yaml"
	_ = ioutil.WriteFile(fileName, yamlConfig, 0644)
	var processingchanTest = make(chan processingItemTest, 1000)
	cfg, _ := config.ReloadConfig("test_config.yaml")

	fullSrcDir, _ := getAbsPath(config_test.SrcDir, "")

	srcItems, _ := ioutil.ReadDir(fullSrcDir)
	go func() {
		for _, item := range srcItems {
			srcFile := item.Name()
			fullSrcFilePath := filepath.Join(fullSrcDir, srcFile)
			processingchanTest <- processingItemTest{srcFile, fullSrcFilePath, cfg, item.Size()}
		}
		close(processingchanTest)
	}()

	for i := 0; i < b.N; i++ {
		for item := range processingchanTest {
			fullDstFilePath, _ := getAbsPath(item.params.DstDir, item.srcFile)
			copyFile(item.fullSrcFilePath, fullDstFilePath)
		}
	}

}
