// Extended-attribute preservation for in-place migrations: moving content
// into an encrypted (or back into a plaintext) copy would silently drop
// SELinux labels, file capabilities (e.g. cap_net_raw on ping) and user.*
// xattrs. The copy is best-effort: each attribute that cannot be read or
// re-applied is logged and skipped, never fatal to the migration.
package fscrypt

import (
	"bytes"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// copyXattrs copies every extended attribute of src onto dst. Filesystems
// without xattr support are treated as "nothing to copy".
func copyXattrs(src, dst string) error {
	names, err := listXattrNames(src)
	if err != nil {
		if errors.Is(err, unix.ENOTSUP) {
			return nil
		}
		return fmt.Errorf("list xattrs of %s: %w", src, err)
	}
	for _, name := range names {
		size, err := unix.Lgetxattr(src, name, nil)
		if err != nil {
			log.Warnf("reading xattr %q of %s: %v (skipped)", name, src, err)
			continue
		}
		value := make([]byte, size)
		n, err := unix.Lgetxattr(src, name, value)
		if err != nil {
			log.Warnf("reading xattr %q of %s: %v (skipped)", name, src, err)
			continue
		}
		if err := unix.Lsetxattr(dst, name, value[:n], 0); err != nil {
			log.Warnf("setting xattr %q on %s: %v (skipped)", name, dst, err)
		}
	}
	return nil
}

// listXattrNames returns the extended attribute names of path. The buffer
// grows on ERANGE so a concurrent attribute addition between the sizing
// call and the read cannot truncate the list.
func listXattrNames(path string) ([]string, error) {
	size := 256
	for {
		buf := make([]byte, size)
		n, err := unix.Llistxattr(path, buf)
		if err == nil {
			return splitXattrNames(buf[:n]), nil
		}
		if !errors.Is(err, unix.ERANGE) {
			return nil, err
		}
		size *= 2
	}
}

// splitXattrNames decodes the NUL-separated name list returned by
// listxattr(2).
func splitXattrNames(buf []byte) []string {
	var names []string
	for len(buf) > 0 {
		i := bytes.IndexByte(buf, 0)
		if i < 0 {
			break
		}
		names = append(names, string(buf[:i]))
		buf = buf[i+1:]
	}
	return names
}
