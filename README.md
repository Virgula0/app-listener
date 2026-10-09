# app-listener

On a typical Linux desktop, every program you run as your user can read your SSH keys, your
browser cookies, your cloud credentials and your AI agents' tokens. That is exactly what
info-stealers and malicious packages (an `npm` postinstall script, a compromised PyPI dependency,
a fake VS Code extension) rely on: they don't need root, only your files.

The same holds for the coding agents and AI CLIs now running in your terminal (Claude Code, Gemini
CLI, opencode, Cursor, Copilot). They run commands and read files on your behalf, with your
permissions. A prompt injection hidden in a README, an issue or a web page, or simply a wrong
step, is enough for one of them to read `~/.aws/credentials`, your browser profile or another
app's tokens.

app-listener lets you decide **which programs may open which directories**, and the Linux kernel
enforces it. `ssh` can read `~/.ssh`; a script that runs `cat ~/.ssh/id_ed25519` gets
`Operation not permitted`. It can also simply show you who touches a directory or what a program
connects to, so you know what to protect.

```console
$ cat ~/.ssh/id_ed25519
cat: /home/alice/.ssh/id_ed25519: Operation not permitted
$ ssh -T git@github.com
Hi alice! You've successfully authenticated, but GitHub does not provide shell access.
```

> [!NOTE]
> A vibe-coding experiment, coded mainly with DeepSeek V4 Flash and Claude Code. Not for
> production systems.

## What makes it useful

- **Each app keeps to its own data.** Every protected directory has its own whitelist: Claude
  Code's config is readable by Claude Code, Discord's tokens by Discord, `~/.ssh` by the SSH tools.
  An agent, a CLI or a dependency that reaches outside its own directory is denied, even though it
  runs as the same user as everything else.
- **Mandatory access control, kept simple.** Linux file permissions (DAC) only ask *which user*
  is opening a file, so everything you run passes. app-listener adds a mandatory layer that also
  asks *which program*, and the file's owner can't override it. The policy is a plain list of
  directories and the programs allowed in each: no AppArmor or SELinux profiles to write, and no
  AppArmor or SELinux needed at all.
- **Set up automatically.** The binary installs with one command in seconds. The wizard then finds
  the sensitive directories on your system and writes the whole policy for you.
- **Programs are identified by the file they run from, not by name.** Rules key on the
  executable's inode. Renaming a binary to `ssh` or faking its process name gets it nothing.
- **Enforced in the kernel with eBPF LSM.** The check sits behind every way of reaching a file
  (plain `open`, io_uring, hard links, renames, raw disk reads, `/proc` tricks), and the test
  suite carries [a corpus of bypass attempts](integrationtests/exploits/README.md) to prove it.
- **Encrypted when unguarded.** In daemon mode, protected directories are encrypted with fscrypt
  and unlocked only after the guard is attached. When the daemon is stopped, they are locked.
- **Knows where your secrets are.** The installer ships a catalog of about 70 sensitive locations
  (SSH, GnuPG, cloud CLIs, Claude Code and other AI agents, browsers, password managers, messaging
  apps, crypto wallets) with a ready-made whitelist for each.
- **Fails closed.** If the kernel can't enforce (BPF-LSM not active, a required hook missing),
  it refuses to start rather than pretending to protect you.

## Modes

One binary, one subcommand per job. The monitors only observe; the guards and the daemon enforce.

### monitor: see what touches a directory

```bash
sudo app-listener monitor -w ~/.ssh
```

Lists every open, read, write, rename, delete and metadata change under a path, with the program
that did it. Nothing is blocked. Good for finding out what a new tool reads before deciding what
to protect.

![monitor](media/monitor.gif)

[Reference](docs/monitor.md)

### guard: block access to a directory

```bash
sudo app-listener guard ~/secret -w /usr/bin/cat     # only cat may open it
sudo app-listener guard ~/secret -b /usr/bin/rm      # everyone but rm
```

Enforces a whitelist or blacklist on a path for as long as it runs, with a live view of what was
allowed and blocked. Useful for testing a whitelist before making it permanent.

![guard](media/guard.gif)

[Reference](docs/guard.md)

### network-monitor: see what a program talks to

```bash
sudo app-listener network-monitor /usr/bin/curl
```

Shows the connections, DNS queries and traffic of the programs you name.

[Reference](docs/network-monitor.md)

### network-guard: decide who gets network access

```bash
sudo app-listener network-guard -b /usr/bin/curl                          # cut one tool off
sudo app-listener network-guard -w /usr/lib/firefox/firefox --auto-infra   # only Firefox and system services
```

A whitelisted program can't be hijacked to get through: starting it with `LD_PRELOAD` or making
it load a planted library costs it network access, and other programs can't attach a debugger to
it.

[Reference](docs/network-guard.md)

### daemon: permanent protection

```bash
sudo app-listener install
```

The mode most people want. An interactive wizard finds the sensitive directories on your system,
lets you pick which to protect, encrypts them and installs a systemd service that guards them from
boot. It follows package upgrades and self-updating apps (Discord, VS Code, Steam) without manual
steps.

[Reference](docs/daemon.md)

### edit-protected: change a protected file

```bash
sudo app-listener edit-protected
```

Once a directory is protected, your editor can't open it either. `edit-protected` opens it in a
built-in editor. With the daemon running and an edit password set at install time, it can also
give another program temporary access (`--forward`) or change the daemon's config, without
stopping anything.

[Reference](docs/edit-protected.md)

## Quick start

**1. Check your system.** You need Linux with kernel 5.17 or newer and BPF-LSM active. Arch Linux
and Ubuntu 24.04 are fully supported; on Ubuntu, BPF-LSM must be switched on with a kernel boot
parameter. Both `x86_64` and `aarch64` are released. [Details](docs/compatibility.md)

**2. Install the binary.** The installer runs the compatibility check, verifies the release
signature and puts `app-listener` in `/usr/local/sbin`. It does not touch your files.

```bash
curl -fsSL https://raw.githubusercontent.com/Virgula0/app-listener/main/scripts/install.sh | sudo bash
```

**3. Try a mode** from the list above, or **protect your directories** with the daemon:

```bash
sudo app-listener install
```

**4. Verify.**

```bash
systemctl is-active app-listener-daemon          # active
cat ~/.ssh/id_ed25519                            # Operation not permitted
sudo journalctl -u app-listener-daemon -f        # see what gets denied
```

To undo everything: `sudo app-listener uninstall` decrypts the directories and removes the service.
To update: `sudo app-listener update`.

Prefer building from source? See [Installation](docs/installation.md#build-from-source).

## What it does not do

It decides *which program* may read a file, not *why*. A whitelisted tool that prints its secret
on request (`gh auth token`, `ssh-add -L`) still does, so never whitelist shells, interpreters or
general tools like `git` or `cp`. It does not defend against an attacker who already has root.
[More](docs/limitations.md)

## Documentation

Full reference, every flag and every config option: **[docs/](docs/README.md)**

- [Compatibility](docs/compatibility.md) and [Installation](docs/installation.md)
- [How it works](docs/how-it-works.md)
- Modes: [monitor](docs/monitor.md), [guard](docs/guard.md),
  [network-monitor](docs/network-monitor.md), [network-guard](docs/network-guard.md),
  [daemon](docs/daemon.md)
- Daemon: [configuration](docs/daemon-config.md), [binary trust](docs/binary-trust.md),
  [catalog and self-updating apps](docs/catalog.md), [edit-protected](docs/edit-protected.md)
- [Limitations](docs/limitations.md), [Troubleshooting](docs/troubleshooting.md),
  [Development](docs/development.md)
