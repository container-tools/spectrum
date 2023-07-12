package builder

import (
	"archive/tar"
	"fmt"
	"os"
	"strings"
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

func TestTarSingleEntry(t *testing.T) {
	var tmpFile1 *os.File
	var err error
	if tmpFile1, err = os.CreateTemp("", "camel-k-*.txt"); err != nil {
		t.Error(err)
	}

	assert.Nil(t, tmpFile1.Close())
	assert.Nil(t, os.WriteFile(tmpFile1.Name(), []byte(`
	This is for simple testing
	`), 0o400))

	tarFileName, err := tarPackage(tmpFile1.Name(), "/path/to/target", false)
	assert.NoError(t, err)
	r, err := os.Open(tarFileName)
	assert.NoError(t, err)
	tr := tar.NewReader(r)
	assert.NotNil(t, tr)
	header, err := tr.Next()
	assert.NotNil(t, tr)
	assert.True(t, strings.HasPrefix(header.Name, "/path/to/target/camel-k-"))
	assert.True(t, strings.HasSuffix(header.Name, ".txt"))
	// Assuming you're testing with a non-root user!
	assert.NotEqual(t, 0, header.Uid)
	assert.NotEqual(t, 0, header.Gid)
}

func TestTarDir(t *testing.T) {
	var tmpDir string
	var tmpFile1 *os.File
	var err error
	if tmpDir, err = os.MkdirTemp("", "dir-*"); err != nil {
		t.Error(err)
	}
	if tmpFile1, err = os.CreateTemp(tmpDir, "camel-k-*.txt"); err != nil {
		t.Error(err)
	}

	assert.Nil(t, tmpFile1.Close())
	assert.Nil(t, os.WriteFile(tmpFile1.Name(), []byte(`
	This is for simple testing
	`), 0o400))

	tarFileName, err := tarPackage(tmpDir, "/path/to/target", false)
	assert.NoError(t, err)
	r, err := os.Open(tarFileName)
	assert.NoError(t, err)
	tr := tar.NewReader(r)
	assert.NotNil(t, tr)
	header, err := tr.Next()
	assert.NotNil(t, tr)
	assert.True(t, strings.HasPrefix(header.Name, "/path/to/target/camel-k-"))
	assert.True(t, strings.HasSuffix(header.Name, ".txt"))
}

func TestTarDirRecursive(t *testing.T) {
	tmpDir1, err := os.MkdirTemp("", "dir1-*")
	assert.NoError(t, err)
	tmpDir2, err := os.MkdirTemp(tmpDir1, "dir2-*")
	assert.NoError(t, err)
	tmpFile1, err := os.CreateTemp(tmpDir2, "camel-k-*.txt")
	assert.NoError(t, err)

	assert.Nil(t, tmpFile1.Close())
	assert.Nil(t, os.WriteFile(tmpFile1.Name(), []byte(`
	This is for simple testing
	`), 0o400))

	tarFileName, err := tarPackage(tmpDir1, "/path/to/target", true)
	assert.NoError(t, err)
	r, err := os.Open(tarFileName)
	assert.NoError(t, err)
	tr := tar.NewReader(r)
	assert.NotNil(t, tr)
	// The first headers are the directory
	header, err := tr.Next()
	assert.NotNil(t, tr)
	header, err = tr.Next()
	assert.NotNil(t, tr)

	// Now the file
	header, err = tr.Next()
	assert.NotNil(t, tr)
	assert.True(t, strings.HasPrefix(header.Name, "/path/to/target/dir2"))
	assert.True(t, strings.HasSuffix(header.Name, ".txt"))
}
