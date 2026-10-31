package guard

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// bunTrailer ends a `bun build --compile` exe's module graph, which Linux builds keep in the
// `.bun` ELF section: u64 length, graph (Offsets.byte_count bytes), Offsets, trailer.
const bunTrailer = "\n---- Bun! ----\n"

// bunOffsetsSize is Bun's StandaloneModuleGraph.Offsets: byte_count u64, modules_ptr and
// compile_exec_argv_ptr {u32,u32}, entry_point_id u32, flags u32 (the last word).
const bunOffsetsSize = 32

// bunDisableAutoloadBunfig is Offsets.flags' disable_autoload_bunfig bit.
const bunDisableAutoloadBunfig = 1 << 1

// bunAutoloadsBunfig reports whether the Bun image f holds runs ./bunfig.toml's `preload` from its
// working directory at startup. Only a compiled exe whose graph validates and sets
// disable_autoload_bunfig doesn't: the bun CLI, another layout or an unreadable one count (fail
// closed: the cost is a refused secret after the app opens a file of that name).
func bunAutoloadsBunfig(f *os.File) bool {
	r, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return true
	}
	defer r.Close()
	ef, err := elf.NewFile(r)
	if err != nil {
		return true
	}
	sec := ef.Section(".bun")
	if sec == nil || sec.Type == elf.SHT_NOBITS {
		return true
	}
	flags, ok := bunGraphFlags(sec, sec.Size)
	return !ok || flags&bunDisableAutoloadBunfig == 0
}

// bunGraphFlags reads Offsets.flags from a `.bun` section of size bytes; ok only when its length
// prefix, byte_count and trailer agree.
func bunGraphFlags(sec io.ReaderAt, size uint64) (flags uint32, ok bool) {
	const tailSize uint64 = bunOffsetsSize + uint64(len(bunTrailer))
	if size < 8+tailSize {
		return 0, false
	}
	var head [8]byte
	tail := make([]byte, tailSize)
	if _, err := sec.ReadAt(head[:], 0); err != nil {
		return 0, false
	}
	if _, err := sec.ReadAt(tail, int64(size-tailSize)); err != nil { //nolint:gosec // size < 2^63: an ELF section
		return 0, false
	}
	if string(tail[bunOffsetsSize:]) != bunTrailer {
		return 0, false
	}
	length := binary.LittleEndian.Uint64(head[:])
	byteCount := binary.LittleEndian.Uint64(tail[:8])
	if length != size-8 || byteCount != length-tailSize {
		return 0, false
	}
	return binary.LittleEndian.Uint32(tail[bunOffsetsSize-4 : bunOffsetsSize]), true
}
