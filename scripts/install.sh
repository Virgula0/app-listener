#!/usr/bin/env bash
#
# install.sh — one-line installer for app-listener.
#
#   curl -fsSL https://raw.githubusercontent.com/Virgula0/app-listener/main/scripts/install.sh | sudo bash
#   curl -fsSL https://raw.githubusercontent.com/Virgula0/app-listener/main/scripts/install.sh | sudo bash -s -- --channel prerelease
#
# What it does, in order:
#   1. runs scripts/check-compatibility.sh — a static host check. If the host
#      cannot run the guard/daemon modes the installer aborts and touches
#      nothing.
#   2. downloads the latest release of the selected --channel from GitHub
#      (release [default] | prerelease — same channel model as
#      `app-listener update`).
#   3. verifies the release the same way `app-listener update` does: the
#      Ed25519 signature of the sha256 checksum against the embedded release
#      public key, then the checksum against the downloaded binary, then the
#      GitHub-provided asset digest.
#   4. installs the verified binary at /usr/local/sbin/app-listener and the
#      /usr/local/bin/app-listener PATH symlink (atomic replace).
#
# It does NOT install or start the protective systemd daemon. Run
#   sudo app-listener install
# afterwards — that step is deliberately manual.
#
# Environment overrides (testing only):
#   APP_LISTENER_REF   git ref for the fetched check-compatibility.sh (default: main)

set -euo pipefail

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[0;33m'; BOLD='\033[1m'; NC='\033[0m'

info() { printf "  ${GREEN}[*]${NC} %s\n" "$*"; }
warn() { printf "  ${YELLOW}[!]${NC} %s\n" "$*"; }
die()  { printf "\n  ${RED}${BOLD}installation aborted:${NC} %s\n" "$*" >&2; exit 1; }

##############################################################################
# Constants
##############################################################################

REPO="Virgula0/app-listener"
API_BASE="https://api.github.com"
REF="${APP_LISTENER_REF:-main}"
RAW_BASE="https://raw.githubusercontent.com/${REPO}/${REF}"

RELEASE_ASSET_BINARY="app-listener"
RELEASE_ASSET_CHECKSUM="app-listener.sha256"
RELEASE_ASSET_SIGNATURE="app-listener.sha256.sig"

INSTALL_PATH="/usr/local/sbin/app-listener"   # systemd.InstallBinaryPath
SYMLINK_PATH="/usr/local/bin/app-listener"    # systemd.BinSymlinkPath

# Release signing public key — must match certificates/app-listener-release.pub
# in the repository. The matching private key lives only in a GitHub Actions
# secret; a release whose checksum is not signed by this key is rejected.
RELEASE_PUBKEY="-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAhxyOi12wFMeWgaqKjUI1YoGMM6ly8Mr+wM6xFDQxu8k=
-----END PUBLIC KEY-----"

##############################################################################
# Arguments
##############################################################################

CHANNEL="release"

usage() {
	cat <<EOF
app-listener installer

Usage: install.sh [--channel release|prerelease]

  --channel   release    latest stable (non-pre-release) GitHub release [default]
              prerelease latest pre-YYYYMMDD-<sha> build
  -h, --help  show this help
EOF
}

while [ $# -gt 0 ]; do
	case "$1" in
	--channel)
		[ $# -ge 2 ] || die "--channel needs a value (release | prerelease)"
		CHANNEL="$2"; shift 2 ;;
	--channel=*)
		CHANNEL="${1#*=}"; shift ;;
	-h | --help)
		usage; exit 0 ;;
	*)
		die "unknown argument: $1 (try --help)" ;;
	esac
done

case "$CHANNEL" in
release | stable)          WANT_PRERELEASE="false" ;;
prerelease | pre-release)  WANT_PRERELEASE="true" ;;
*)                         die "invalid --channel '$CHANNEL' (use: release | prerelease)" ;;
esac

##############################################################################
# Preflight
##############################################################################

printf "${BOLD}app-listener installer${NC}  (channel: %s, repo: %s)\n\n" "$CHANNEL" "$REPO"

[ "$(id -u)" -eq 0 ] || die "run as root: pipe into 'sudo bash', or 'sudo $0'"

missing=""
for tool in curl openssl jq sha256sum od mktemp; do
	command -v "$tool" >/dev/null 2>&1 || missing="$missing $tool"
done
if [ -n "$missing" ]; then
	die "missing required tool(s):$missing
     Debian/Ubuntu: apt install curl jq openssl coreutils
     Arch:          pacman -S curl jq openssl coreutils"
fi

if [ "$(uname -m)" != "x86_64" ]; then
	die "the published release is linux/amd64 only; this host is $(uname -m)"
fi

TMP="$(mktemp -d "${TMPDIR:-/tmp}/app-listener-install.XXXXXX")"
trap 'rm -rf "$TMP"' EXIT

##############################################################################
# 1. Host compatibility
##############################################################################

printf "${BOLD}[1/4] Host compatibility check${NC}\n"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" 2>/dev/null && pwd || true)"
if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/check-compatibility.sh" ]; then
	COMPAT="$SCRIPT_DIR/check-compatibility.sh"
	info "using $COMPAT"
else
	COMPAT="$TMP/check-compatibility.sh"
	info "fetching check-compatibility.sh from $REPO@$REF"
	curl -fsSL --proto '=https' --tlsv1.2 -o "$COMPAT" "$RAW_BASE/scripts/check-compatibility.sh" \
		|| die "could not download check-compatibility.sh"
fi

if ! bash "$COMPAT"; then
	die "the host cannot run app-listener — fix the blockers listed above and re-run. Nothing was installed."
fi

##############################################################################
# 2. Resolve and download the latest release
##############################################################################

printf "\n${BOLD}[2/4] Downloading the latest release (%s channel)${NC}\n" "$CHANNEL"

releases_json="$(curl -fsSL --proto '=https' --tlsv1.2 \
	-H 'Accept: application/vnd.github+json' \
	-H 'User-Agent: app-listener-install' \
	"$API_BASE/repos/$REPO/releases?per_page=100")" \
	|| die "could not list releases of $REPO"

line="$(printf '%s' "$releases_json" | jq -r --argjson pre "$WANT_PRERELEASE" '
	([ .[] | select(.prerelease == $pre) ] | sort_by(.published_at) | reverse | .[0]) as $r
	| if $r == null then "" else
		[ $r.tag_name,
		  ( [ $r.assets[] | select(.name=="app-listener")            | .browser_download_url ] | first // "" ),
		  ( [ $r.assets[] | select(.name=="app-listener")            | .digest ]              | first // "" ),
		  ( [ $r.assets[] | select(.name=="app-listener.sha256")     | .browser_download_url ] | first // "" ),
		  ( [ $r.assets[] | select(.name=="app-listener.sha256.sig") | .browser_download_url ] | first // "" )
		] | @tsv
	  end
')" || die "could not parse the GitHub API response"

[ -n "$line" ] || die "no '$CHANNEL' release published for $REPO.$( [ "$CHANNEL" = release ] && printf ' Try: --channel prerelease' )"

IFS=$'\t' read -r TAG BIN_URL BIN_DIGEST SHA_URL SIG_URL <<<"$line" || true

[ -n "${TAG:-}" ]     || die "could not read the release tag from the GitHub API response"
[ -n "${BIN_URL:-}" ] || die "release $TAG has no '$RELEASE_ASSET_BINARY' asset"
[ -n "${SHA_URL:-}" ] || die "release $TAG has no '$RELEASE_ASSET_CHECKSUM' asset"
[ -n "${SIG_URL:-}" ] || die "release $TAG has no '$RELEASE_ASSET_SIGNATURE' asset"

info "latest release on the $CHANNEL channel: $TAG"

download() {
	curl -fsSL --proto '=https' --tlsv1.2 --retry 3 --retry-delay 1 \
		-H 'User-Agent: app-listener-install' -o "$2" "$1" \
		|| die "download failed: $1"
}

info "fetching release assets (~45 MiB) ..."
download "$BIN_URL" "$TMP/$RELEASE_ASSET_BINARY"
download "$SHA_URL" "$TMP/$RELEASE_ASSET_CHECKSUM"
download "$SIG_URL" "$TMP/$RELEASE_ASSET_SIGNATURE"
info "downloaded binary, checksum and signature"

##############################################################################
# 3. Verify (Ed25519 signature -> checksum -> asset digest -> sanity)
##############################################################################

printf "\n${BOLD}[3/4] Verifying the release${NC}\n"

printf '%s\n' "$RELEASE_PUBKEY" > "$TMP/release.pub"

openssl pkeyutl -verify -rawin -pubin -inkey "$TMP/release.pub" \
	-in "$TMP/$RELEASE_ASSET_CHECKSUM" -sigfile "$TMP/$RELEASE_ASSET_SIGNATURE" >/dev/null 2>&1 \
	|| die "Ed25519 signature check FAILED — the checksum is not signed by the app-listener release key. The download is rejected."
info "Ed25519 signature of the checksum: OK"

expected_sha="$(awk 'NR==1{print tolower($1)}' "$TMP/$RELEASE_ASSET_CHECKSUM")"
case "$expected_sha" in
*[!0-9a-f]* | "") die "malformed checksum file" ;;
esac
[ "${#expected_sha}" -eq 64 ] || die "malformed checksum file (expected a 64-hex sha256)"

actual_sha="$(sha256sum "$TMP/$RELEASE_ASSET_BINARY" | awk '{print tolower($1)}')"
[ "$actual_sha" = "$expected_sha" ] \
	|| die "sha256 mismatch — checksum says $expected_sha, download is $actual_sha. The download is rejected."
info "sha256 checksum vs binary: OK"

if [ -n "${BIN_DIGEST:-}" ]; then
	want_digest="$(printf '%s' "${BIN_DIGEST#sha256:}" | tr 'A-F' 'a-f')"
	[ "$want_digest" = "$expected_sha" ] \
		|| die "GitHub asset digest ($BIN_DIGEST) disagrees with the signed checksum ($expected_sha). The download is rejected."
	info "GitHub asset digest: OK"
else
	warn "release $TAG exposes no asset digest — skipping that cross-check (signature + checksum still enforced)"
fi

magic="$(od -An -tx1 -N4 "$TMP/$RELEASE_ASSET_BINARY" | tr -d ' \n')"
[ "$magic" = "7f454c46" ] || die "the downloaded file is not an ELF executable"

chmod 0700 "$TMP/$RELEASE_ASSET_BINARY"
if ! ver_out="$("$TMP/$RELEASE_ASSET_BINARY" --version 2>&1)"; then
	die "the downloaded binary does not run: $ver_out"
fi
printf '%s' "$ver_out" | grep -qF "$TAG" \
	|| die "the downloaded binary reports '$ver_out', expected tag '$TAG' — refusing to install"
info "binary runs and reports $TAG"

##############################################################################
# 4. Install the binary (no systemd, no config)
##############################################################################

printf "\n${BOLD}[4/4] Installing the binary${NC}\n"

install_dir="$(dirname "$INSTALL_PATH")"
mkdir -p "$install_dir"
staged="$(mktemp "$install_dir/.app-listener-new-XXXXXX")"
cat "$TMP/$RELEASE_ASSET_BINARY" > "$staged"
chmod 0700 "$staged"
mv -f "$staged" "$INSTALL_PATH"        # atomic replace — safe even if a daemon is running
info "installed $INSTALL_PATH"

link_dir="$(dirname "$SYMLINK_PATH")"
if [ ! -e "$SYMLINK_PATH" ] && [ ! -L "$SYMLINK_PATH" ]; then
	if [ -d "$link_dir" ]; then
		ln -s "$INSTALL_PATH" "$SYMLINK_PATH" && info "symlinked $SYMLINK_PATH -> $INSTALL_PATH"
	else
		warn "$link_dir does not exist — skipping the PATH symlink; invoke $INSTALL_PATH directly"
	fi
elif [ -L "$SYMLINK_PATH" ] && [ "$(readlink "$SYMLINK_PATH")" = "$INSTALL_PATH" ]; then
	info "PATH symlink $SYMLINK_PATH already correct"
else
	warn "$SYMLINK_PATH already exists and is not our symlink — leaving it untouched; invoke $INSTALL_PATH directly"
fi

if command -v systemctl >/dev/null 2>&1 && systemctl is-active --quiet app-listener-daemon 2>/dev/null; then
	warn "the app-listener-daemon service is running the previous binary — restart it when ready:"
	warn "  sudo systemctl restart app-listener-daemon"
fi

##############################################################################

printf "\n  ${GREEN}${BOLD}app-listener %s is installed.${NC}\n\n" "$TAG"
printf "  The protective daemon is NOT installed automatically. To protect\n"
printf "  directories with fscrypt + the systemd daemon, run now:\n\n"
printf "      ${BOLD}sudo app-listener install${NC}\n\n"
printf "  Other modes work immediately, e.g.:\n\n"
printf "      sudo app-listener monitor -w /tmp\n"
printf "      sudo app-listener guard /secret -w /usr/bin/cat\n\n"
