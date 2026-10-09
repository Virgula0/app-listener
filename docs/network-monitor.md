# network-monitor: see what a program talks to

Traces the network operations (TCP, UDP, DNS) of the listed binaries only. Use it to see where a
tool phones home, or which hosts an app needs before you restrict it with
[network-guard](network-guard.md).

```bash
sudo app-listener network-monitor /usr/bin/curl /usr/bin/wget
sudo app-listener network-monitor /usr/bin/curl -e CONNECT,DNS
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `<binary>` | required | Binaries to watch (positional, repeatable) |
| `-e, --events <list>` | all | `CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS` |
| `--headless` | `false` | No TUI; log `NETEVENT\|` lines to stderr |
| `--serve[=host:port]` | off | Mirror the TUI to a browser |

## Events

| Event | Meaning | Hook |
|-------|---------|------|
| `CONNECT` | Outbound connect | `sys_enter_connect` |
| `ACCEPT` | Inbound accepted (TCP) | `kretprobe/inet_csk_accept`, `sys_enter_accept[4]` |
| `SEND` / `RECV` | Data sent / received | `sys_enter_sendto`, `sendmsg` / `recvfrom`, `recvmsg` |
| `CLOSE` | Socket close | `sys_enter_close` |
| `DNS` | Query to port 53 or 853 | `sys_enter_connect` / `sendto` |

## Multicall binaries

A uutils applet is watched as that applet only. A multicall binary that can't be told apart is
watched by its file, with a warning. See [Binary trust](binary-trust.md#multicall-binaries).
