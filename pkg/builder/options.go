package builder

import "io"

type Options struct {
	PullInsecure  bool
	PushInsecure  bool
	PullConfigDir string
	PushConfigDir string
	Base          string
	Target        string
	Stdout        io.Writer
	Stderr        io.Writer
}
