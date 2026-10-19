# monitor: see what touches a path

Traces file operations under the watched paths. Nothing is blocked. Use it to learn which
programs read a directory before you guard it, or to check what an install script or a new
dependency does on disk.

![monitor](../media/tui/monitor.gif)

```bash
sudo app-listener monitor ~/.ssh
sudo app-listener monitor /var/log --recursive --depth 3
sudo app-listener monitor /path/to/file.txt -e OPEN,READ
```

## Events

| Event | Covers |
|---|---|
| `OPEN`, `READ`, `WRITE`, `MMAP` | Data I/O, whatever the syscall path (io_uring, splice, sendfile, mmap) |
| `DELETE`, `RENAME`, `SYMLINK`, `HARDLINK`, `MKDIR` | Tree changes |
| `ATTR` | chmod, chown, utimes, truncate, setxattr |
| `STAT` | stat, access, readlink |
| `MKNOD` | Device node creation |

A watched path on its own mount (a tmpfs `/tmp`, a bind mount, a btrfs subvolume) is reported
under the watched path. The monitor never stats files outside the watch, so running it while the
daemon is active does not trip the daemon's guards on their protected files.

## Flags

`app-listener monitor PATH...` takes one or more files or directories as arguments.

| Flag | Default | Description |
|------|---------|-------------|
| `-r, --recursive` | `false` | Recurse into subdirectories |
| `-d, --depth <n>` | `0` | Max depth (needs `--recursive`; `0` = unlimited) |
| `-e, --events <list>` | all | Comma-separated subset of the events above |
| `--headless` | `false` | No TUI; log `EVENT\|` lines to stderr |
| `--serve[=host:port]` | off | Mirror the TUI to a browser, see [How it works](how-it-works.md#watching-from-a-browser---serve) |
