package main

import (
	"context"
	"errors"
	. "final-project/internal/config"
	"final-project/internal/tosser"
	"final-project/pkg/logging"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	errNotModified = errors.New("not modified")
	cfg            *Config
)

func main() {
	config_init := Config{
		MaxCopyThreads: 4,
		RescanInterval: 3600,
		LogLevel:       "INFO",
		SrcDir:         "",
		DstDir:         "",
		SkipIfexists:   "yes",
	}
	yamlConfig, err := yaml.Marshal(&config_init)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	if _, err := os.Stat("/etc/tosser"); os.IsNotExist(err) {
		err := os.Mkdir("/etc/tosser", os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
	configName := "config.yaml"
	if _, err = os.Stat("/etc/tosser/" + configName); os.IsNotExist(err) {
		err = ioutil.WriteFile("/etc/tosser/"+configName, yamlConfig, 0644)
		if err != nil {
			panic("Unable to write data into the file")
		}
	}

	help_msg := "tosser is a file synchronizer between two directories." +
		"\nYou need to specify paths for source directory and destination directory in the config file which path is /etc/tosser/config.yaml," +
		"\n\nUsage:\nrun - starts synchronization between two directories\nhelp - for more information"

	configPath := flag.String("config", "/etc/tosser/config.yaml", "Конфиг для синхронизации директорий")
	flag.Parse() // get the arguments from command line

	arg1 := flag.Arg(0)
	if arg1 == "help" || arg1 == "" {
		fmt.Println(help_msg)
	} else if arg1 != "run" && flag.Arg(0) != "" {
		fmt.Println("Unknown command:", arg1)
		os.Exit(1)
	} else {
		cfg, err = ReloadConfig(*configPath)
		if err != nil {
			fmt.Println(err)
			if err != errNotModified {
				log.Fatalf("Could not load '%s': %s", *configPath, err)
			}
		}
		if err := logging.InitLogger(cfg); err != nil {
			log.Fatalln(err)
		}
		msg := `You need to specify source and destination directories in /erc/tosser/config.yaml, for more information use "tosser help"`
		if cfg.DstDir == "" && cfg.SrcDir == "" && len(os.Args) < 2 {
			fmt.Println(msg)
		} else {
			if arg1 == "run" && cfg.DstDir != "" && cfg.SrcDir != "" {
				ctx, cancel := context.WithCancel(context.Background())
				if info, err := os.Stat(cfg.SrcDir); !os.IsNotExist(err) {
					if !info.IsDir() {
						panic(cfg.SrcDir + " is not a directory!")
					}
				} else {
					panic("directory does not exists: " + cfg.SrcDir)
				}
				go tosser.ScanLoop(ctx, cfg)
				tosser.HandleSignals(cancel)
			} else {
				fmt.Println(msg)
			}
		}
	}
}
