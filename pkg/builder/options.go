package builder

type Options struct {
	PullInsecure  bool
	PushInsecure  bool
	PullConfigDir string
	PushConfigDir string
	Base          string
	Target        string
}
