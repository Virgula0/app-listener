package daemonconfig

import (
	"fmt"
	"strconv"
	"strings"
)

// MaxInodes bounds for [settings] max_inodes. Below the floor a typical config's populate fills the
// map, and every inode it can't hold is guarded by the ancestor walk alone (degraded coverage); the
// ceiling keeps the preallocated map from taking hundreds of MiB of kernel memory.
const (
	MaxInodesFloor   = 1 << 15
	MaxInodesCeiling = 1 << 22
)

// applySetting parses one `key: value` line of a [settings] block. Unknown or repeated keys fail:
// a misspelled setting must not silently keep the default.
func applySetting(cfg *Config, line string, lineNo int) error {
	key, value, ok := strings.Cut(line, ":")
	key, value = strings.TrimSpace(key), strings.TrimSpace(value)
	if !ok || key == "" || value == "" {
		return fmt.Errorf("daemon config line %d: a [settings] line is `key: value`, got %q", lineNo, line)
	}
	if key != "max_inodes" {
		return fmt.Errorf("daemon config line %d: unknown setting %q", lineNo, key)
	}
	if cfg.MaxInodes != 0 {
		return fmt.Errorf("daemon config line %d: duplicate max_inodes", lineNo)
	}
	n, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fmt.Errorf("daemon config line %d: max_inodes %q is not a number", lineNo, value)
	}
	if n < MaxInodesFloor || n > MaxInodesCeiling {
		return fmt.Errorf("daemon config line %d: max_inodes %d is outside %d..%d", lineNo, n,
			MaxInodesFloor, MaxInodesCeiling)
	}
	cfg.MaxInodes = uint32(n)
	return nil
}
