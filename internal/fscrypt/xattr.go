// Xattr preservation for in-place migrations: copying between plaintext and
// encrypted copies would drop SELinux labels, file capabilities and user.*
// xattrs. Best-effort: failed reads/writes are logged and skipped, not fatal.
package fscrypt

import (
	"bytes"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// copyXattrs copies every xattr of src onto dst; ENOTSUP means nothing to copy.
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

// listXattrNames returns path's xattr names, growing the buffer on ERANGE
// so concurrent additions cannot truncate the list.
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

// splitXattrNames decodes the NUL-separated listxattr(2) name list.
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
