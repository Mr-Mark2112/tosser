package tosser

import (
	"context"
	"errors"
	"final-project/internal/config"
	"final-project/pkg/logging"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"
)

const configFileName = "/etc/tosser/config.yaml"

var (
	errNotModified          = errors.New("not modified")
	processing              = NewProcessingCache()
	processingchan          = make(chan processingItem, 1000)
	processingchanForRemove = make(chan processingRemoveItem, 1000)
	copyError               = make(chan error, 0)
)

type processingItem struct {
	srcFile         string
	fullSrcFilePath string
	params          *config.Config
	size            int64
}

type processingRemoveItem struct {
	dstFile         string
	fullDstFilePath string
	params          *config.Config
	size            int64
}

func HandleSignals(cancel context.CancelFunc) {

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	for {
		sig := <-sigCh
		switch sig {
		case os.Interrupt:
			cancel()
			log.Println("programm stopped")
			return
		}
	}
}

func ScanLoop(ctx context.Context, cfg *config.Config) {
	for i := 1; i <= cfg.MaxCopyThreads; i++ {
		go processItem(cfg)
		go processItemForDelete(cfg)
	}
	for {
		select {
		case <-ctx.Done():
			logging.Info.Printf("programm stopped")
		case <-time.After(time.Duration(cfg.RescanInterval) * time.Second):
			go processScanDir(cfg)
		}
		cfgTmp, err := config.ReloadConfig(configFileName)
		if err != nil {
			if err != errNotModified {
				logging.Debug.Println("readconfig:", err)
			}
		} else {
			logging.Info.Println("rescanning config file")
			cfg = cfgTmp
			logging.InitLogger(cfg)
		}
	}
}

func processScanDir(cfg *config.Config) {
	fullDstDir, err := getAbsPath(cfg.DstDir, "")
	if err != nil {
		logging.Errorln("error culculating destination absolute path:", cfg.DstDir, err)
	}
	fullSrcDir, err := getAbsPath(cfg.SrcDir, "")

	if err != nil {
		logging.Errorln("error culculating source absolute path:", cfg.SrcDir, err)
	}

	if err := os.MkdirAll(fullSrcDir, os.ModeDir); err != nil {
		logging.Errorln("error creating the directory:", fullSrcDir, err)
	}

	if processing.check(fullSrcDir) {
		logging.Debug.Println("directory is already bieng scanned:", fullSrcDir)
	}
	logging.Debug.Println("scanning diractory...", fullSrcDir)
	srcItems, err := ioutil.ReadDir(fullSrcDir)
	if err != nil {
		logging.Errorln(err)
		logging.Debug.Printf("processing directory completed: %s", fullSrcDir)
	}

	dstItems, err := ioutil.ReadDir(fullDstDir)
	if err != nil {
		logging.Errorln(err)
		logging.Debug.Printf("processing directory completed: %s", fullSrcDir)
	}
	processing.add(fullSrcDir)
	processItems(srcItems, fullSrcDir, cfg)
	processItemsForDelete(srcItems, dstItems, fullSrcDir, cfg)
	processing.del(fullSrcDir)
}

func processItems(srcItems []os.FileInfo, fullSrcDir string, params *config.Config) {
	for _, item := range srcItems {
		if !item.Mode().IsRegular() {
			continue
		}
		srcFile := item.Name()
		fullSrcFilePath := filepath.Join(fullSrcDir, srcFile)
		if processing.check(fullSrcFilePath) {
			logging.Debug.Println("the file is already being processed:", fullSrcFilePath)
		}
		//add a file to cache
		processing.add(fullSrcFilePath)
		processingchan <- processingItem{srcFile, fullSrcFilePath, params, item.Size()}
	}
}

func processItemsForDelete(srcItems []os.FileInfo, dstItems []os.FileInfo, fullSrcDir string, params *config.Config) {
	for _, item := range dstItems {
		if !item.Mode().IsRegular() {
			continue
		}
		dstFile := item.Name()
		fullDstFilePath := filepath.Join(fullSrcDir, dstFile)
		processingchanForRemove <- processingRemoveItem{dstFile, fullDstFilePath, params, item.Size()}
	}
}

func getAbsPath(dir, file string) (string, error) {
	filePath := filepath.Join(dir, file)
	abspath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err

	}
	return abspath, nil
}
func processItem(cfg *config.Config) {
	for item := range processingchan {
		fullDstFilePath, err := getAbsPath(item.params.DstDir, item.srcFile)
		if err != nil {
			logging.Errorln("error culculating absolute path:", err)
			continue
		}
		fullDstFileDir := filepath.Dir(fullDstFilePath)
		if _, err := os.Stat(fullDstFilePath); err == nil {
			switch cfg.SkipIfexists {
			case "no":
				log.Printf("file is exists . %s skipIfexists=%s. Replacing file in destination directory.", fullDstFilePath, cfg.SkipIfexists)
				if err := os.Remove(fullDstFilePath); err != nil {
					log.Println("deleting file error:", err)
					continue
				}
			case "yes":
				logging.Debug.Printf("file is already exists '%s', skipping", fullDstFilePath)
				continue
			default:
				log.Printf("file is exists. %s Unknown value: skipIfexists=%s. Skiping file.", fullDstFilePath, cfg.SkipIfexists)
				continue
			}
		}
		if err := os.MkdirAll(fullDstFileDir, os.ModeDir); err != nil {
			logging.Errorln("error creating the directory:", fullDstFileDir, err)
			continue
		}
		logging.Info.Printf("copying item: '%s', size: %s", item.srcFile, convertSize(item.size))
		err = copyFile(item.fullSrcFilePath, fullDstFilePath)
		if err != nil {
			log.Printf("an error occured while copying file: %s , err:%s", item.fullSrcFilePath, err)
		}
		// if err != nil {
		// 	copyError <- err
		// 	return
		// }
		processing.del(item.fullSrcFilePath)
	}
}

func processItemForDelete(cfg *config.Config) {
	for item := range processingchanForRemove {
		checkIfSrcFileExists, err := getAbsPath(cfg.SrcDir, item.dstFile)
		if err != nil {
			logging.Errorln("error culculating absolute path:", err)
			continue
		}
		fileForRM, err := getAbsPath(item.params.DstDir, item.dstFile)
		if err != nil {
			logging.Errorln("error culculating absolute path:", err)
			continue
		}
		if _, err := os.Stat(checkIfSrcFileExists); err != nil {
			if _, err := os.Stat(cfg.SrcDir); os.IsNotExist(err) {
				panic("source directory hase been removed, synchronization stopped")
			} else {
				if err = deleteFile(fileForRM); err != nil {
					fmt.Println("error removing file:", err)
				}
				logging.Info.Printf("removing item: '%s', size: %s", checkIfSrcFileExists, convertSize(item.size))
			}
		}

	}
}

func copyFile(src string, dst string) (err error) {
	sourcefile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dst)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = io.Copy(destfile, sourcefile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if closeErr := destfile.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	sourceinfo, err := os.Stat(src)
	if err == nil {
		err = os.Chmod(dst, sourceinfo.Mode())
	}
	return err
}

func deleteFile(rmFile string) (err error) {
	if err := os.Remove(rmFile); err != nil {
		logging.Errorln(err)
	}
	return err
}

var suffixes [5]string

func convertSize(sizeInBytes int64) string {
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	if sizeInBytes == 0 {
		sizeInBytes = 1
	}
	base := math.Log(float64(sizeInBytes)) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 1)
	getSuffix := suffixes[int(math.Floor(base))]
	res := fmt.Sprintln(strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix))
	return res
}
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
