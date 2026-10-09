# Limitations

app-listener answers one question: *which executable* may touch this file or this socket. Knowing
where that answer stops is part of using it well.

## Whitelisted programs that disclose on request

A whitelisted program that hands out its secret when asked still does. With `/usr/bin/aws`
whitelisted on `~/.aws`, `aws configure export-credentials` prints the access key and secret to
whoever runs it; many apps have a similar export, debug or "show token" command. The guard
authorizes which executable opens the file, not why. Any process that can run the whitelisted
binary with arguments of its choosing gets whatever that binary prints.

- Whitelist only what needs the resource.
- Prefer apps that keep secrets in a running agent.
- Never whitelist a general tool: a shell, an interpreter (`node`, `python`, `java`), `git`, `cp`,
  `curl`, ... reads the resource on behalf of any script. The installer and the daemon warn about
  every resource that whitelists one.

## Interpreters and JITs

A whitelisted interpreter runs whatever script its caller passes, and a JIT inside a whitelisted
app runs code from anonymous memory, which is never judged. The [launch scan](binary-trust.md#launch-scan)
catches the common ways of injecting code into Node, Electron, Chromium and JVM apps through
environment or flags; it doesn't make `node` itself safe to whitelist.

## Root

Root can stop the daemon or load its own BPF programs. app-listener
protects against processes running **as you** (a malicious package, a compromised dependency, an
info-stealer), not against an attacker who already has root.

## Binaries you can write

Identity is the inode. If a whitelisted binary is writable by a user, that user controls what the
whitelisted identity does. The daemon limits who may replace such binaries
([Binary trust](binary-trust.md)), but root-owned binaries are the strong case.

## Event masks are not confinement

A binary restricted to `READ` that execs another whitelisted binary escapes its mask. Masks are a
least-privilege hint.

## Already-running processes

- Processes started before a guard are judged once from `/proc`; code injected into them earlier
  can't be seen. Restart whitelisted apps after starting `network-guard -w`.
- Multicall applets started before a guard have no attested applet and are denied until they exec
  again ([details](binary-trust.md#multicall-binaries)).
- A process started before its binary was replaced keeps its access.

## Side effects to expect with the daemon

- **Backups and snapshots.** Raw block-device reads and btrfs snapshot/send are blocked on any
  filesystem that hosts a guarded resource, so `snapper`, `timeshift`, `btrbk`, `snap-pac` and raw
  imaging stop working there ([details](daemon.md#filesystem-wide-gates)).
- **Process inspection.** Programs outside a resource's whitelist can't read `/proc/<pid>` of a
  process that read protected content. Screen-sharing portals are allowed by default through
  [`[inspectors]`](daemon-config.md#process-inspectors); other tools may break silently (see
  [Troubleshooting](troubleshooting.md#an-app-breaks-and-nothing-is-logged)).
- **Bun apps from desktop menus** still use `/tmp` and are denied their native library
  ([details](installation.md#bun-based-apps-opencode)).

## Project status

This is a vibe-coding experiment, coded mainly with DeepSeek V4 Flash and Claude Code. It is not
meant for production systems.
