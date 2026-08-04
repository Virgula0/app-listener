package install

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed all:daemon-samples
var sampleFS embed.FS

// SampleFiles lists the bundled daemon-samples files (service units, the
// pacman reload hook and the sample config) shipped inside the binary.
func SampleFiles() ([]string, error) {
	entries, err := fs.ReadDir(sampleFS, "daemon-samples")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// SampleContent returns the embedded content of one daemon-samples file.
func SampleContent(name string) ([]byte, error) {
	data, err := sampleFS.ReadFile("daemon-samples/" + name)
	if err != nil {
		return nil, fmt.Errorf("embedded sample %s: %w", name, err)
	}
	return data, nil
}
