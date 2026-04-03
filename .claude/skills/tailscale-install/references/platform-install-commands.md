<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Platform Install Commands

Exhaustive per-platform Tailscale installation commands. The agent should
use platform detection results from SKILL.md to select the correct section.

## Linux — Quick Install (All Distros)

The official one-liner detects the distro and installs automatically:

```bash
curl -fsSL https://tailscale.com/install.sh | sh
```

This works on Debian, Ubuntu, CentOS, RHEL, Fedora, Amazon Linux, Arch,
and most other distributions. Use the distro-specific sections below when
the one-liner is unavailable or the environment requires explicit repo
configuration.

## Linux — Debian / Ubuntu (apt)

```bash
# Add Tailscale GPG key and repository
curl -fsSL https://pkgs.tailscale.com/stable/ubuntu/$(lsb_release -cs).noarmor.gpg \
  | sudo tee /usr/share/keyrings/tailscale-archive-keyring.gpg >/dev/null
curl -fsSL https://pkgs.tailscale.com/stable/ubuntu/$(lsb_release -cs).tailscale-keyring.list \
  | sudo tee /etc/apt/sources.list.d/tailscale.list

# Install
sudo apt-get update
sudo apt-get install -y tailscale

# Enable and start daemon
sudo systemctl enable --now tailscaled
```

For Debian (not Ubuntu), replace `ubuntu` with `debian` in the URLs above.

## Linux — RHEL / Fedora / CentOS (dnf/yum)

```bash
# Add Tailscale repository
sudo dnf config-manager --add-repo \
  https://pkgs.tailscale.com/stable/fedora/tailscale.repo

# Install
sudo dnf install -y tailscale

# Enable and start daemon
sudo systemctl enable --now tailscaled
```

For CentOS/RHEL, replace `fedora` with `centos` or `rhel` in the repo URL.
For older systems without `dnf`, use `yum` instead.

## Linux — Arch (pacman)

```bash
sudo pacman -S tailscale

# Enable and start daemon
sudo systemctl enable --now tailscaled
```

## Linux — Alpine (apk)

```bash
sudo apk add tailscale

# Start daemon (OpenRC)
sudo rc-update add tailscale
sudo rc-service tailscale start
```

## Linux — openSUSE / SUSE (zypper)

```bash
# Add Tailscale repository
sudo zypper addrepo -g -r \
  https://pkgs.tailscale.com/stable/opensuse/tumbleweed/tailscale.repo

# Install
sudo zypper install tailscale

# Enable and start daemon
sudo systemctl enable --now tailscaled
```

For Leap, replace `tumbleweed` with the Leap version (e.g., `15.5`).

## Linux — Raspberry Pi OS

Use the Debian instructions above. Raspberry Pi OS is Debian-based:

```bash
curl -fsSL https://tailscale.com/install.sh | sh
sudo systemctl enable --now tailscaled
```

## macOS — Homebrew

```bash
brew install --cask tailscale
```

After installation, open Tailscale from Applications to start the system
extension. The CLI is available at:

```bash
# Homebrew install puts CLI in PATH
tailscale version
```

## macOS — Mac App Store

Install "Tailscale" from the Mac App Store. The App Store version uses a
network extension (not kernel extension).

The CLI is NOT in PATH by default. Add an alias:

```bash
alias tailscale="/Applications/Tailscale.app/Contents/MacOS/Tailscale"
```

Or add to shell profile for persistence:

```bash
echo 'alias tailscale="/Applications/Tailscale.app/Contents/MacOS/Tailscale"' \
  >> ~/.zshrc
```

## Windows — WinGet

```powershell
winget install tailscale.tailscale
```

Restart the terminal after installation to pick up the CLI in PATH.

## Windows — Chocolatey

```powershell
choco install tailscale
```

## Windows — MSI Installer

Download the latest MSI from the Tailscale website and run:

```powershell
msiexec /i tailscale-setup.msi /quiet
```

## WSL2 — Do NOT Install Linux Package

Tailscale should run on the **Windows host**, not inside WSL2. The WSL2
virtual machine shares the host network stack, so installing Tailscale
inside WSL2 creates conflicting WireGuard tunnels.

**Correct approach:**

1. Install Tailscale on the Windows host (WinGet, Chocolatey, or MSI)
2. Access the CLI from WSL2 via `tailscale.exe`:

```bash
# From inside WSL2
tailscale.exe version
tailscale.exe status
tailscale.exe ping <peer>
```

1. MagicDNS hostnames resolve automatically from WSL2 when the Windows
   host is connected to the tailnet.

**If `tailscale.exe` is not found in WSL2**, ensure the Windows PATH is
shared with WSL2. Check `/etc/wsl.conf`:

```ini
[interop]
appendWindowsPath = true
```

## Docker

```bash
docker run -d \
  --name tailscale \
  --hostname tailscale-container \
  --cap-add NET_ADMIN \
  --cap-add NET_RAW \
  --device /dev/net/tun:/dev/net/tun \
  -v tailscale-state:/var/lib/tailscale \
  -e TS_AUTHKEY=tskey-auth-<key> \
  -e TS_EXTRA_ARGS=--advertise-tags=tag:container \
  tailscale/tailscale:latest
```

Required capabilities: `NET_ADMIN` and `NET_RAW`. The `/dev/net/tun`
device must be available. Use `TS_AUTHKEY` for unattended authentication.

## Docker Compose

```yaml
services:
  tailscale:
    image: tailscale/tailscale:latest
    hostname: tailscale-service
    cap_add:
      - NET_ADMIN
      - NET_RAW
    devices:
      - /dev/net/tun:/dev/net/tun
    volumes:
      - tailscale-state:/var/lib/tailscale
    environment:
      - TS_AUTHKEY=tskey-auth-<key>
      - TS_EXTRA_ARGS=--advertise-tags=tag:container

volumes:
  tailscale-state:
```
