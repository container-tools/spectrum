package builder

import "io"

type Options struct {
	PullInsecure  bool
	PushInsecure  bool
	PullConfigDir string
	PushConfigDir string
	Base          string
	Target        string
	Annotations   map[string]string
	Stdout        io.Writer
	Stderr        io.Writer
	Recursive     bool
}
