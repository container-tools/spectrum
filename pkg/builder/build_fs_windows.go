//go:build windows
// +build windows

package builder

import (
	"archive/tar"
	"io/fs"

	"golang.org/x/sys/windows"
)

func prepareHeader(tp, name string, fi fs.FileInfo) *tar.Header {
	// prepare the tar header
	header := new(tar.Header)
	header.Name = name
	header.Size = fi.Size()
	header.Mode = int64(fi.Mode().Perm())
	fileSys := fi.Sys()
	if fileSys != nil {
		header.Uid = windows.Getuid()
		header.Gid = windows.Getgid()
	} else {
		StepLogger.Printf("Warning: could not read UID/GID. Assuming default (root) permissions.")
	}
	header.ModTime = fi.ModTime()

	return header
}
