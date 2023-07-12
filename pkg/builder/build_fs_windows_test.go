//go:build windows
// +build windows

package builder

import (
	"archive/tar"
	"os"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/sys/windows"

	"github.com/stretchr/testify/assert"
)

func TestTarPermissionsSingleEntry(t *testing.T) {
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
	assert.Equal(t, windows.Getuid(), header.Uid)
	assert.Equal(t, windows.Getgid(), header.Gid)
}

func TestTarDirAllEntriesRecursive(t *testing.T) {
	tmpDir1, err := os.MkdirTemp("", "camel-k-dir1-*")
	assert.NoError(t, err)
	tmpDir2, err := os.MkdirTemp(tmpDir1, "camel-k-dir2-*")
	assert.NoError(t, err)
	_, err = os.MkdirTemp(tmpDir1, "camel-k-dir3-*")
	assert.NoError(t, err)
	tmpFile1, err := os.CreateTemp(tmpDir2, "camel-k-file-*.txt")
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

	header, err := tr.Next()
	assert.NotNil(t, header)
	assert.Nil(t, err)
	matched, _ := regexp.MatchString("/path/to/target/", header.Name)
	assert.True(t, matched)
	assert.Equal(t, windows.Getuid(), header.Uid)
	assert.Equal(t, windows.Getgid(), header.Gid)

	header, err = tr.Next()
	assert.NotNil(t, header)
	assert.Nil(t, err)
	matched, _ = regexp.MatchString("/path/to/target/camel-k-dir2-.*", header.Name)
	assert.True(t, matched)
	assert.Equal(t, windows.Getuid(), header.Uid)
	assert.Equal(t, windows.Getgid(), header.Gid)

	header, err = tr.Next()
	assert.NotNil(t, header)
	assert.Nil(t, err)
	matched, _ = regexp.MatchString("/path/to/target/camel-k-dir2-.*/camel-k-file-.*\\.txt", header.Name)
	assert.True(t, matched)
	assert.Equal(t, windows.Getuid(), header.Uid)
	assert.Equal(t, windows.Getgid(), header.Gid)

	header, err = tr.Next()
	assert.NotNil(t, header)
	assert.Nil(t, err)
	matched, _ = regexp.MatchString("/path/to/target/camel-k-dir3-.*/", header.Name)
	assert.True(t, matched)
	assert.Equal(t, windows.Getuid(), header.Uid)
	assert.Equal(t, windows.Getgid(), header.Gid)

	header, err = tr.Next()
	assert.Nil(t, header)
	assert.NotNil(t, err)
	assert.Equal(t, "EOF", err.Error())
}
