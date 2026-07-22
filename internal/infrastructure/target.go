package ebpf

type Target struct {
	Path  string // original path (for display)
	Dir   string // parent dir (for files) or the dir itself
	File  string // basename (empty for directories)
	IsDir bool
}
