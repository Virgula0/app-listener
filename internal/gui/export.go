//go:build gui

package gui

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/safeio"
)

var exportFormats = []string{"CSV", "JSON"}

// encodeTable renders a table as CSV (header row first) or as a JSON array of objects keyed by
// column title.
func encodeTable(format string, headers []string, rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	if format == "JSON" {
		objs := make([]map[string]string, len(rows))
		for i, r := range rows {
			objs[i] = make(map[string]string, len(headers))
			for c, h := range headers {
				objs[i][h] = r[c]
			}
		}
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		err := enc.Encode(objs)
		return buf.Bytes(), err
	}
	w := csv.NewWriter(&buf)
	if err := w.Write(headers); err != nil {
		return nil, err
	}
	if err := w.WriteAll(rows); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// invokingIDs is the uid/gid sudo or pkexec ran us for, (-1, -1) when none: exports are handed to
// that user, not left root-owned in their home.
func invokingIDs() (uid, gid int) {
	u, ok := invokingUID()
	if !ok {
		return -1, -1
	}
	if g, err := strconv.Atoi(os.Getenv("SUDO_GID")); err == nil && g > 0 {
		return int(u), g
	}
	if usr, err := user.LookupId(strconv.FormatUint(uint64(u), 10)); err == nil {
		if g, err := strconv.Atoi(usr.Gid); err == nil {
			return int(u), g
		}
	}
	return -1, -1
}

// writeUserFile writes data to the absolute path as root without following a symlink at any
// component (the folder is usually the invoking user's, who could plant one to redirect the write),
// replacing an existing file atomically.
func writeUserFile(path string, data []byte) error {
	if !filepath.IsAbs(path) {
		return errors.New("choose an absolute folder")
	}
	dir, name := filepath.Split(filepath.Clean(path))
	if name == "" || name == "." || name == ".." {
		return errors.New("choose a file name")
	}
	root, err := unix.Open("/", unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return err
	}
	var comps []string
	for _, c := range strings.Split(filepath.Clean(dir), "/") {
		if c != "" {
			comps = append(comps, c)
		}
	}
	dirFD, err := safeio.DescendNoFollow(root, comps, false, 0)
	defer unix.Close(dirFD)
	if err != nil {
		return err
	}
	uid, gid := invokingIDs()
	return safeio.AtomicWriteAt(dirFD, name, "."+name+".tmp", data, 0o644, uid, gid)
}

// exportTable asks for a format, folder and file name, then writes the table's rows as shown
// (filters applied).
func (s *shell) exportTable(baseName string, t *dataTable) {
	headers, rows := t.texts()
	format := newRadioGroup(exportFormats, nil)
	format.Horizontal = true
	format.Required = true
	format.Selected = exportFormats[0]
	folder := widget.NewEntry()
	folder.SetText(invokingHome())
	name := widget.NewEntry()
	stem := baseName + "-" + time.Now().Format("20060102-150405")
	name.SetText(stem + ".csv")
	format.OnChanged = func(f string) {
		ext := "." + strings.ToLower(f)
		base := strings.TrimSuffix(strings.TrimSuffix(name.Text, ".csv"), ".json")
		name.SetText(base + ext)
	}
	browse := newButton("Folder…", icon("folder"), func() {
		s.pickPath(pickFolder, "Export folder", folder.Text, folder.SetText)
	})
	form := widget.NewForm(
		widget.NewFormItem("Format", format),
		widget.NewFormItem("Folder", container.NewBorder(nil, nil, nil, browse, folder)),
		widget.NewFormItem("File name", name),
	)
	note := mutedText(fmt.Sprintf("%d rows, as currently filtered. The file is owned by you, not root.", len(rows)))
	d := dialog.NewCustomConfirm("Export table", "Export", "Cancel", scrollable(container.NewVBox(form, note)),
		func(ok bool) {
			if !ok {
				return
			}
			data, err := encodeTable(format.Selected, headers, rows)
			if err != nil {
				showErrorDialog(err, s.win)
				return
			}
			path := filepath.Join(strings.TrimSpace(folder.Text), strings.TrimSpace(name.Text))
			s.confirmOverwrite(path, func() {
				if err := writeUserFile(path, data); err != nil {
					showErrorDialog(fmt.Errorf("export to %s: %w", path, err), s.win)
					return
				}
				showInformation("Exported", fmt.Sprintf("Wrote %d rows to %s.", len(rows), path), s.win)
			})
		}, s.win)
	s.showSized(d, fyne.NewSize(640, 300))
}

// confirmOverwrite runs write at once, or after the user agrees to replace an existing path.
func (s *shell) confirmOverwrite(path string, write func()) {
	if _, err := os.Lstat(path); err != nil {
		write()
		return
	}
	showConfirm("Replace file?", path+" already exists. Replace it?", func(ok bool) {
		if ok {
			write()
		}
	}, s.win)
}
