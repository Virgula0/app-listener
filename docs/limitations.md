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

- `gpg`: `gpg -e -r <their key> <file>` or `gpg -c <file>` wraps any file gpg reads (pass needs
  `-e`, so it can't be refused). The catalog keeps the secret keys out of its reach instead: a
  [`[narrow]`](config-reference.md#narrow-path) section lets only `gpg-agent` into
  `~/.gnupg/private-keys-v1.d`, and gpg only writes `openpgp-revocs.d`. What gpg still reads (the
  public keyring, the trust database, the configuration) is what it needs.
- `kubectl`, `oc`: `config view --raw` is refused, but
  `kubectl create configmap x --from-file=<kubeconfig> --dry-run=client -o yaml` prints the file.
- `rclone`: `config show` and `config dump` are refused, but `rclone cat <rclone.conf>` prints it.
  Set a config password (`rclone config encryption set`) to keep it encrypted at rest.

When a tool's whole job is to read the secret it holds, guarding the directory and whitelisting the
tool protects nothing: whoever can run the tool reads the secret through it, and no launch rule
closes that. Such entries are deliberately left out of the catalog rather than offered with a false
sense of protection:

- **`docker`** (`~/.docker`): docker reads `config.json` on every command, so it would disclose the
  registry logins on request. Keep them in the Secret Service keyring (`credsStore` set to a keyring
  helper) and gate that with a [secret scope](config-reference.md#secret-name), or use short-lived
  tokens.
- **`age`** (`~/.config/age`): `age -d` decrypts with the identity it reads, and `age -e -r
  <recipient>` re-encrypts that identity to any key the caller names. A file guard cannot stop age
  from handing out the key it is meant to use.

Python- and Node-launched tools are the opposite case: because a `#!` script runs as its
interpreter, a whitelist line naming the tool never matches, so the directory is locked against
everyone, the tool included. The entry protects the files but the app can't use them while guarded
(see [Script launchers](catalog.md#what-an-entry-cant-cover)).

A browser shows any file it can read, its own profile included: `brave --headless --dump-dom
file://<profile file>` prints it, and `firefox --screenshot` renders it into an image. Headless,
remote-control and logging launches of a browser therefore run without its profile
([the table below](#commands-that-stop-working)).
So does a browser started on a profile outside its guarded tree (`--profile`, `--user-data-dir`,
another `HOME`): that profile's settings and extensions are its caller's
([Browser profiles](binary-trust.md#browser-profiles)).

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
  [runtime section](config-reference.md#runtime): the catalog has one for `gpg-agent` and one
  for the installer's `ssh-agent`. An agent you start elsewhere, and the session bus, can be asked
  by any of your programs (see below).
- **Impersonating a guarded agent.** The programs allowed to create sockets in a runtime directory
  can be run by any of your programs: the agents themselves, and systemd's user manager, which binds
  the GnuPG socket units and takes orders from `systemctl --user` and `systemd-run --user`. With the
  real agent stopped and its socket gone, a program can have one of them listen at the agent's name
  for it. Clients then talk to that listener: it can show its own passphrase prompt, and read what
  `gpg` or `ssh` send it. What `ssh-add` sends (the key) only an agent program may read: see
  [Agents as oracles](#secrets-served-over-ipc).
- **Logins with no daemon running.** A runtime directory created then (by the agent's own unit) is
  guarded once the daemon starts, but may already hold sockets bound by any program. A clean stop
  keeps the sections already there guarded ([details](daemon.md#shutdown-and-crashes)), but a
  logout and a new login while the daemon is down make a new, unguarded directory.
- **Abstract sockets** have no file and are never judged.

## Secrets served over IPC

The guard judges who opens a file or a guarded socket. A program that serves secrets on request
over a channel the guard doesn't see answers whoever asks:

- **The keyring.** The Secret Service (GNOME keyring, KDE Wallet, KeePassXC) is reached over the
  session bus, so any of your programs can read an unlocked collection. `gh`, Copilot CLI, Gemini CLI,
  the Antigravity CLI, Codex (in keyring mode), JetBrains IDEs, Zed, ProtonVPN, `docker` with a `secretservice` helper and `git`'s libsecret helper
  keep tokens there. With gnome-keyring, [secret scopes](config-reference.md#secret-name)
  restrict the items they cover to their whitelists, with the limits below. Elsewhere, prefer file
  storage inside a guarded tree where the app offers it
  (`gh auth login --insecure-storage`), or let KeePassXC serve the Secret Service with "Confirm when
  passwords are retrieved by clients" on, so every read asks you first.
- **Agents as oracles.** A guarded agent socket stops other programs from connecting, but a
  whitelisted client run by malware still asks on its behalf: `ssh-add -L`, `ssh` signing to a host
  of the attacker's choice, `gpg --decrypt`. `git credential fill`, run by anyone, prints what your
  credential helper (`store`, `cache`, libsecret, Git Credential Manager, `gh auth git-credential`)
  hands back. `ssh-add` sends the key it read to whatever socket it is pointed at, so it may reach
  only an agent: a socket inside a guarded tree (the catalog's `ssh-agent`, gcr and `gpg-agent`
  runtime sections, or one in a guarded directory), whose connection only a whitelisted agent
  program (`ssh-agent`, `gpg-agent`, `gcr-ssh-agent`, 1Password) may read or answer. An agent started
  with `eval $(ssh-agent)` listens in `/tmp`, so `ssh-add` is refused it (`DAEMON DENIED op=CONNECT
  comm=ssh-add resource=unguarded-socket`). A listener another program made inside the tree (a
  socket unit systemd binds, an `ssh -L` forward) accepts the connection, but its reads are refused
  (`resource=agent-only-connection`). Its own user lookup is refused the same way as `/tmp`, so for
  a user known only to `sssd` or LDAP `ssh-add` fails (`nscd` is skipped, local users work). Any
  `gpg-agent` counts as an agent, including one a program starts on a home directory of its own,
  and `gpg-agent` hands its keys to any client that asks (see [Whitelisted programs that disclose on request](#whitelisted-programs-that-disclose-on-request)).
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
  has none, unless a [port section](config-reference.md#port-tcpn) guards the port. Even then:
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
  [`[inspectors]`](config-reference.md#inspectors); other tools may break silently (see
  [Troubleshooting](troubleshooting.md#an-app-breaks-and-nothing-is-logged)).
- **Bun apps from desktop menus** still use `/tmp` and are denied their native library
  ([details](installation.md#bun-based-apps-opencode)).
- **Refused launches.** Some commands of whitelisted programs run without their secrets, so they
  fail: see [Commands that stop working](#commands-that-stop-working).

## Commands that stop working

These launches of a whitelisted program still run, but can't open anything the daemon guards (the
journal logs `TRUST DENIED op=ARGV` or `op=LAUNCH`), so they usually fail. The match doesn't look
at what a flag points to: `ssh -F /tmp/test` is refused like `ssh -F ~/.ssh/id_ed25519`, since the
file named and the file opened can differ (a symlink, a relative path) by the time the program
opens it. The full list is in [Binary trust](binary-trust.md#launch-rules).

| Program | What stops working | Why it is refused | Do instead |
|---|---|---|---|
| `ssh` | `-F <any file>`, in any option cluster (`-4F/tmp/c`), also in a remote command (`ssh host ls -F`) and through `GIT_SSH_COMMAND="ssh -F ..."` | `ssh` prints each line of a config it can't parse, a private key included | Put the settings in `~/.ssh/config`, or pass them as `-o Option=value` |
| `sshd` | Any run by a user other than root (`sshd -t -f <file>`) | `sshd -t -f` prints the file it names line by line | Run it as root |
| `ssh-add` | Talking to an agent outside a guarded tree (`eval $(ssh-agent)` in `/tmp`) | It would hand the key it read to that socket | Use the catalog's `ssh-agent` |
| Chrome, Chromium, Brave | `--headless`, `--dump-dom`, `--screenshot`, `--print-to-pdf` | A headless browser prints or renders any file it reads | Run automation on a separate profile (`--user-data-dir=/tmp/p`) |
| Firefox | `--headless`, `--screenshot`, `MOZ_HEADLESS=1` | Same | Run automation on a separate profile (`firefox --profile /tmp/p`) |
| Browsers | `--remote-debugging-port`/`-pipe`, and for Firefox `--marionette`, `--start-debugger-server`, `--jsdebugger` | Whoever connects drives the browser and reads its cookies | Run Selenium, Playwright or Puppeteer on a separate profile |
| Chrome, Chromium, Brave | `--allow-file-access-from-files`, `--disable-web-security` | A local page reads any file the browser can read and sends it out | Serve the page over `http://localhost`; test on a separate profile |
| Firefox, Chrome, Chromium, Brave | Any profile outside the guarded tree: `--profile`, `--user-data-dir`, `-P` naming one elsewhere, another `HOME` or `XDG_CONFIG_HOME` | The profile's `user.js` or extensions would run in the whitelisted browser | Nothing to change: a test or automation profile runs, only without the real profile's secrets |
| Browsers | `--log-net-log`, `--net-log-capture-mode`, `--ssl-key-log-file`, `SSLKEYLOGFILE`, Firefox's `MOZ_LOG` / `-MOZ_LOG` | The log holds cookies or the TLS keys | Debug on a separate profile |
| Browsers, Electron apps | `--proxy-server`, `--proxy-pac-url`, `--ignore-certificate-errors*` | The app would send its cookies and tokens to the caller's server | Use the desktop's proxy settings |
| Browsers, Electron apps | `--load-extension`, `--disable-extensions-except`, `--pack-extension`, `--remote-debugging-*`, `--inspect*` | They run code of the caller's choosing in the app | Install extensions from the store; debug on a separate profile |
| Node, Electron, Bun, JVM apps | Code-loading variables and flags (`NODE_OPTIONS=--require`, `ELECTRON_RUN_AS_NODE`, `BUN_OPTIONS`, `JAVA_TOOL_OPTIONS=-javaagent`, ...) | Same | Leave them unset for whitelisted apps |
| `gh` | `auth token`, `auth status -t` / `--show-token`, `alias set`, `alias import` | They print the token, or make `gh` run a command of the caller's choosing | `gh auth status` without `-t` |
| `aws` | `configure export-credentials`, `configure get` | They print the keys | Read `~/.aws` with an editor through `edit-protected` |
| `gpg` | `--export-secret-keys`, `--export-secret-subkeys`, `--export-secret-ssh-key`, `--enarmor`, `--store`, in every abbreviation | They print the secret keys, or wrap any file `gpg` reads | None today |
| `kubectl`, `oc` | `config view --raw` | It prints the client keys and tokens | `config view` (redacted) |
| `rclone` | `config show`, `config dump` | They print every remote's tokens | `rclone config redacted` |
| `bitcoin-cli`, `litecoin-cli`, `lightning-cli` | Key dumps and every send or payment command | They print keys or spend funds | None today |
| JetBrains IDEs, `pi`, `opencode`, `codex` | `installPlugins`; `pi -e` / `--extension`; `OPENCODE_CONFIG*`; `OPENAI_BASE_URL` | They load code, or send the token to a server of the caller's choosing | Install plugins from the IDE; keep these in the app's own config |
| Every whitelisted program | `SSL_CERT_FILE`, `SSL_CERT_DIR`, `NODE_EXTRA_CA_CERTS`, `NODE_TLS_REJECT_UNAUTHORIZED` set to anything but the system bundle | The program would trust a certificate of the caller's choosing | Add the corporate bundle with [`allow_env`](config-reference.md#allow_env) |

## Project status

This is a vibe-coding experiment, coded mainly with DeepSeek V4 Flash and Claude Code. It is not
meant for production systems.
