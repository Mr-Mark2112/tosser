package tosser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Config struct {
	SrcDir         string
	DstDir         string
	MaxCopyThreads int
	RescanInterval int
	LogLevel       string
}

func Test_getAbsPath(t *testing.T) {
	req := require.New(t)
	config := Config{
		SrcDir:         "/root/test_dir_for_project",
		DstDir:         "/root/test_dir_for_project_copy",
		MaxCopyThreads: 4,
		RescanInterval: 3,
		LogLevel:       "INFO",
	}

	t.Run("simple test", func(t *testing.T) {
		fullSrcFilePath := "/root/test_dir_for_project/test_file_1"
		fullDstFilePath := "/root/test_dir_for_project_copy/test_file_1"
		res, _ := getAbsPath(config.SrcDir, "test_file_1")
		res2, _ := getAbsPath(config.DstDir, "test_file_1")
		req.NotNil(res)
		req.Equal(fullSrcFilePath, res)
		req.NotNil(res2)
		req.Equal(fullDstFilePath, res2)
	})
}

func Test_copyFile(t *testing.T) {
	req := require.New(t)
	t.Run("simple test", func(t *testing.T) {
		fullSrcFilePath := "/root/test_dir_for_project/test_file_1.sh"
		fullDstFilePath := "/root/test_dir_for_project_copy/test_file_1.sh"
		err_res := copyFile(fullSrcFilePath, fullDstFilePath)
		req.Nil(err_res)
	})
}

func Test_deleteFile(t *testing.T) {
	req := require.New(t)

	t.Run("simple test", func(t *testing.T) {
		fullDstFilePath := "/root/test_dir_for_project_copy/test_file_1.sh"
		err_res := deleteFile(fullDstFilePath)
		req.Nil(err_res)
	})
}

func Test_convertSize(t *testing.T) {
	req := require.New(t)

	t.Run("simple test", func(t *testing.T) {
		formatted_size := "316.9 KB\n"
		res := convertSize(324521)
		req.Equal(formatted_size, res)
	})

	t.Run("huge size", func(t *testing.T) {
		formatted_size := "93 GB\n"
		res := convertSize(99845462323)
		req.Equal(formatted_size, res)
	})

	t.Run("small size", func(t *testing.T) {
		formatted_size := "33 B\n"
		res := convertSize(33)
		req.Equal(formatted_size, res)
	})
}

func Test_(t *testing.T) {
	req := require.New(t)

	t.Run("simple test", func(t *testing.T) {
		formatted_size := "316.9 KB\n"
		res := convertSize(324521)
		req.Equal(formatted_size, res)
	})

	t.Run("huge size", func(t *testing.T) {
		formatted_size := "93 GB\n"
		res := convertSize(99845462323)
		req.Equal(formatted_size, res)
	})

	t.Run("small size", func(t *testing.T) {
		formatted_size := "33 B\n"
		res := convertSize(33)
		req.Equal(formatted_size, res)
	})
}
