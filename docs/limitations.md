# Limitations

app-listener answers one question: *which executable* may touch this file or this socket. Knowing
where that answer stops is part of using it well.

## Whitelisted programs that disclose on request

A whitelisted program that hands out its secret when asked still does. With `/usr/bin/aws`
whitelisted on `~/.aws`, `aws configure export-credentials` prints the access key and secret to
whoever runs it; many apps have a similar export, debug or "show token" command. The guard
authorizes which executable opens the file, not why. Any process that can run the whitelisted
binary with arguments of its choosing gets whatever that binary prints.

[Launch rules](binary-trust.md#launch-rules) refuse the secret to the known disclosure commands
(`gh auth token`, `aws configure export-credentials`, `bitcoin-cli dumpprivkey`), but they are a
tripwire: a CLI has more ways to print a secret than any list names.

Programs that read a configuration file can disclose too: `ssh` and `sshd` print each line of a
config they can't parse, a key included. The catalog's rules refuse `ssh -F` and any `sshd` not run
by root, the routes to a config of the caller's choosing.

Some whitelisted CLIs also process any file they are given, so a rule only stops the one-command
dump, never the copy:

- `gpg`: the rules refuse `--export-secret-keys`, `--export-secret-subkeys`,
  `--export-secret-ssh-key`, `--enarmor` and `--store` in every abbreviation gpg accepts, but
  `gpg -e -r <their key> <file>` or `gpg -c <file>` still wraps any file gpg reads, the keys in
  `~/.gnupg/private-keys-v1.d` included (they stay protected by their passphrase, if they have one).
  pass needs `-e`, so it can't be refused.
- `kubectl`, `oc`: `config view --raw` is refused, but
  `kubectl create configmap x --from-file=<kubeconfig> --dry-run=client -o yaml` prints the file.
- `rclone`: `config show` and `config dump` are refused, but `rclone cat <rclone.conf>` prints it.
  Set a config password (`rclone config encryption set`) to keep it encrypted at rest.

A browser or Electron app sends its cookies and tokens wherever its TLS stack trusts. The launch
scan refuses the flags that redirect it (`--ignore-certificate-errors`, `--proxy-server`), but not a
CA added to `~/.pki/nssdb` (Chromium trusts that store) combined with a desktop or `https_proxy`
proxy setting: both are your own to change.

- Whitelist only what needs the resource.
- Prefer apps that keep secrets in a running agent.
- Never whitelist a general tool: a shell, an interpreter (`node`, `python`, `java`), `git`, `cp`,
  `curl`, a file reader or encoder (`cat`, `dd`, `head`, `base64`, `sed`, `grep`, `openssl`, `nc`),
  an editor, ... reads the resource on behalf of any script. The installer and the daemon warn about
  every resource that whitelists one.

## Unix sockets

A unix socket bound inside a guarded tree (an SSH agent or ControlMaster socket in `~/.ssh`,
`~/.1password/agent.sock`, Core Lightning's `lightning-rpc`) is guarded like a file: connecting,
sending a datagram, or sending and receiving on a connection handed to another program counts as
`CONNECT`, judged by the program doing it. A daemon.conf line with an explicit event list needs
`CONNECT` for a program that talks to such a socket. What stays outside:

- **Sockets outside the tree.** Agents in `/run/user/<uid>` are guarded only through a
  [runtime section](daemon-config.md#runtime-sections): the catalog has one for `gpg-agent` and one
  for the installer's `ssh-agent`. An agent you start elsewhere, and the session bus, can be asked
  by any of your programs (see below).
- **Impersonating a guarded agent.** The programs allowed to create sockets in a runtime directory
  can be run by any of your programs: the agents themselves, and systemd's user manager, which binds
  the GnuPG socket units and takes orders from `systemctl --user` and `systemd-run --user`. With the
  real agent stopped and its socket gone, a program can have one of them listen at the agent's name
  for it. Clients then talk to that listener: it gets no key, but can show its own passphrase
  prompt.
- **Logins with no daemon running.** A runtime directory created then (by the agent's own unit) is
  guarded once the daemon starts, but may already hold sockets bound by any program.
- **Abstract sockets** have no file and are never judged.

## Secrets served over IPC

The guard judges who opens a file or a guarded socket. A program that serves secrets on request
over a channel the guard doesn't see answers whoever asks:

- **The keyring.** The Secret Service (GNOME keyring, KDE Wallet, KeePassXC) is reached over the
  session bus, so any of your programs can read an unlocked collection. `gh`, Copilot CLI, Gemini CLI,
  the Antigravity CLI, Codex (in keyring mode), JetBrains IDEs, Zed, ProtonVPN, `docker` with a `secretservice` helper and `git`'s libsecret helper
  keep tokens there. With gnome-keyring, [secret scopes](daemon-config.md#secret-service-sections)
  restrict the items they cover to their whitelists, with the limits below. Elsewhere, prefer file
  storage inside a guarded tree where the app offers it
  (`gh auth login --insecure-storage`), or let KeePassXC serve the Secret Service with "Confirm when
  passwords are retrieved by clients" on, so every read asks you first.
- **Agents as oracles.** A guarded agent socket stops other programs from connecting, but a
  whitelisted client run by malware still asks on its behalf: `ssh-add -L`, `ssh` signing to a host
  of the attacker's choice, `gpg --decrypt`. `git credential fill`, run by anyone, prints what your
  credential helper (`store`, `cache`, libsecret, Git Credential Manager, `gh auth git-credential`)
  hands back. `ssh-add` may reach only an agent whose socket is inside a guarded tree (the catalog's
  `ssh-agent`, gcr and `gpg-agent` runtime sections, or one in a guarded directory): pointed at
  any other socket it would send the key it read there. An agent started with `eval $(ssh-agent)`
  listens in `/tmp`, so `ssh-add` is refused it (`DAEMON DENIED op=CONNECT comm=ssh-add
  resource=unguarded-socket`). Its own user lookup is refused the same way, so for a user known only
  to `sssd` or LDAP `ssh-add` fails (`nscd` is skipped, local users work).
- **Secret scopes.** The gate covers gnome-keyring's Secret Service only:
  - items no scope covers stay open to every program, as do item labels and attributes;
  - a whitelisted program run by malware still asks on its behalf (`gh auth token`);
  - programs started before the daemon (re)started, or before it first installed its bus policy,
    are refused every scoped item until restarted;
  - the bus is trusted for who sent each call: it must be a root-placed `dbus-broker` or
    `dbus-daemon`, but one you start yourself with your own configuration counts, and so does one
    any of your programs can debug (`kernel.yama.ptrace_scope` 0). After a broker upgrade and a
    daemon restart, a bus still running the old binary can't answer the relay until you log in again;
  - the keyring must be started by its user unit (`gnome-keyring-daemon.service`): one started some
    other way (D-Bus activation outside systemd, a PAM `auto_start`) can't own the Secret Service
    name, so secrets are unavailable;
  - `dbus-monitor` and `busctl monitor` stop working for your user;
  - gnome-keyring's control and pkcs11 sockets in `/run/user/<uid>/keyring` bypass the bus: guard
    them with the catalog's runtime section, which a browser loading the keyring's PKCS#11 module
    then can't use.
- **Local servers.** A service on a loopback port (`bitcoind` RPC on 8332, `ollama` on 11434, `adb`
  on 5037, `opencode serve`) answers any local client that has its credential, or any client when it
  has none, unless a [port section](daemon-config.md#port-sections) guards the port. Even then:
  - a whitelisted client run by malware still asks on its behalf, as with agents;
  - connections opened before the daemon started are not judged;
  - UDP, other hosts and other network namespaces are out of scope;
  - a socket-activated service is listened on by systemd, which needs `MKNOD` itself;
  - a port the service picks at random (the opencode TUI's own server) can't be named;
  - user-mode network stacks (`pasta`, `slirp4netns`) connect as themselves, so a rootless
    container reaching the host's service needs them whitelisted, which admits every container.
- **The display.** Under X11, any client can read the clipboard, other windows' contents and
  keystrokes. Wayland compositors isolate clients, but the clipboard is still shared. Neither is in
  scope.

## Interpreters and JITs

A whitelisted interpreter runs whatever script its caller passes, and a JIT inside a whitelisted
app runs code from anonymous memory, which is never judged. The [launch scan](binary-trust.md#launch-scan)
catches the common ways of injecting code into Node, Electron, Chromium and JVM apps through
environment or flags; it doesn't make `node` itself safe to whitelist.

Some apps load code, or settings that redirect their token, from the directory they start in, which
no launch rule can see:

- `opencode` runs the plugins of the project it starts in (`.opencode/plugins`), in its own process,
  and takes provider base URLs from the project's `opencode.json`;
- `codex -c openai_base_url=...` (any `-c` key) is logged, not refused;
- a JetBrains IDE asks before trusting a project, and pi before loading a project's extensions:
  keep those prompts.

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
  can't be seen. Restart whitelisted apps after starting `network guard -w`.
- Multicall applets started before a guard have no attested applet and are denied until they exec
  again ([details](binary-trust.md#multicall-binaries)).
- A process started before its binary was replaced keeps its access.
- A unix socket connection to a guarded socket made before the guard started (a daemon restart
  counts, a reload does not) is judged by nobody: whoever holds it later can keep using it.

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
- **Refused launches.** A launch that matches a rule runs without its secrets, so the app fails:
  `ssh -F <config>` and `ssh host ls -F`, `git` over ssh with `GIT_SSH_COMMAND="ssh -F ..."`, a
  browser or VS Code started with `--proxy-server`, `ssh-add` to an agent outside a guarded tree.
  Use the user config (`~/.ssh/config`) and the desktop's proxy settings instead.

## Project status

This is a vibe-coding experiment, coded mainly with DeepSeek V4 Flash and Claude Code. It is not
meant for production systems.
