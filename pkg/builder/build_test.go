package builder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	// windows
	windowsPath := "C:\\USER\\HOME\\"
	targetPath := "."
	paths := fmt.Sprintf("%s:%s", windowsPath, targetPath)
	localPath, targetPath, err := getPaths(paths, "windows")

	assert.NoError(t, err)
	assert.Equal(t, windowsPath, localPath)
	assert.Equal(t, ".", targetPath)

	// errors on linux
	_, _, err = getPaths(paths, "linux")
	assert.Error(t, err)

	// windows relative path
	windowsRelPath := "src\\main\\resources"
	paths = fmt.Sprintf("%s:%s", windowsRelPath, targetPath)
	localPath, targetPath, err = getPaths(paths, "windows")

	assert.NoError(t, err)
	assert.Equal(t, windowsRelPath, localPath)
	assert.Equal(t, ".", targetPath)

	// linux
	linuxPath := "/home/user/name/dir"
	paths = fmt.Sprintf("%s:%s", linuxPath, targetPath)
	localPath, targetPath, err = getPaths(paths, "linux")

	assert.NoError(t, err)
	assert.Equal(t, linuxPath, localPath)
	assert.Equal(t, ".", targetPath)
}
