---
name: tailscale-install
description: >
  Install and configure Tailscale across platforms. Detects OS, distro,
  and environment (including WSL2 and containers). Verifies existing
  installations, performs platform-appropriate install, and guides initial
  connection. Use when setting up Tailscale on a new machine, onboarding
  a server to a tailnet, or verifying an existing install.
compatibility: >
  Requires root/admin access for installation. Works on Linux (apt, dnf,
  yum, pacman, apk, zypper), macOS (Homebrew), Windows (WinGet,
  Chocolatey). Detects WSL2 and container environments.
---
<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

# Tailscale Install

Install Tailscale and connect a device to a tailnet. Handles platform
detection, installation, verification, and initial authentication.

## What this skill does

1. Detects the platform (OS, distro, WSL2, container).
2. Checks for an existing Tailscale installation.
3. Installs Tailscale using the appropriate package manager.
4. Verifies the installation and daemon status.
5. Guides initial connection to the tailnet.

## Platform detection

Detect the environment before selecting an install method. Check WSL2
**first** to avoid installing the Linux package inside a WSL2 instance.

```text
1. Check WSL2 (but not containers running on WSL2)
   ├─ grep -qi microsoft /proc/version 2>/dev/null
   ├─ AND [ ! -f /.dockerenv ] (not a container)
   ├─ WSL2 → go to "WSL2 special case"
   └─ Not WSL2 (or container) → continue

2. Check OS
   ├─ uname -s → Linux → detect distro
   ├─ uname -s → Darwin → macOS
   └─ OS is Windows → PowerShell commands

3. Detect Linux distro
   ├─ . /etc/os-release && echo $ID
   ├─ debian, ubuntu → apt
   ├─ fedora, rhel, centos, amzn → dnf/yum
   ├─ arch, manjaro → pacman
   ├─ alpine → apk
   ├─ opensuse*, sles → zypper
   └─ other → curl one-liner fallback
```

## Existing install check

Before installing, check if Tailscale is already present:

```bash
tailscale version 2>/dev/null
```

If installed:

- Report the current version.
- Check daemon status: `sudo systemctl status tailscaled` (Linux).
- Check connection: `tailscale status`.
- Stop here unless the user requests a reinstall or upgrade.

If not installed, proceed to installation.

## WSL2 special case

WSL2 shares the Windows host network stack. Installing Tailscale inside
WSL2 creates conflicting WireGuard tunnels. **Do not install the Linux
package in WSL2.**

Detection (exclude containers running on WSL2 hosts):

```bash
grep -qi microsoft /proc/version 2>/dev/null && [ ! -f /.dockerenv ] && echo "WSL2 detected"
```

When WSL2 is detected:

1. Check for Windows host Tailscale:

   ```bash
   tailscale.exe version 2>/dev/null
   ```

2. If found — report version and verify with `tailscale.exe status`.
3. If not found — instruct the user to install on the Windows host using
   WinGet, Chocolatey, or MSI (see `references/platform-install-commands.md`).
4. Verify `tailscale.exe` is accessible from WSL2. If not, check that
   `appendWindowsPath = true` in `/etc/wsl.conf` under `[interop]`.

## Installation

### Linux (most distros)

```bash
curl -fsSL https://tailscale.com/install.sh | sh
```

### macOS

```bash
brew install --cask tailscale
```

### Windows

```powershell
winget install tailscale.tailscale
```

For distro-specific repo setup (apt, dnf, pacman, apk, zypper), Docker,
Chocolatey, and MSI, see `references/platform-install-commands.md`.

## Post-install verification

After installation, verify three things:

### 1. Binary is available

```bash
tailscale version
```

If not found, check PATH or restart the terminal session.

### 2. Daemon is running

```bash
# Linux (systemd)
sudo systemctl enable --now tailscaled
sudo systemctl status tailscaled

# macOS — open Tailscale from Applications
# Windows — check: Get-Service Tailscale
```

### 3. Connection status

```bash
tailscale status
```

If the device is not yet authenticated, proceed to initial connection.

## Initial connection

### Interactive (desktop/laptop)

```bash
sudo tailscale up
```

This prints an authentication URL. Open it in a browser to sign in.

### Headless (server)

Generate an auth key in the admin console (Settings > Keys), then:

```bash
sudo tailscale up --authkey=tskey-auth-<key>
```

For tags, subnet routes, exit nodes, and all `tailscale up` flags, see
`references/post-install-checklist.md`.

## Error handling

Do **not** silently skip a failed installation. Report the error with
specific remediation steps.

- **Permission denied** — prefix with `sudo` or run as administrator.
- **Package not found** — Tailscale repo not configured. Use the `curl`
  one-liner or add the repo manually (see references).
- **Network error** — check internet connectivity and DNS resolution.
- **Daemon not starting** — check logs with
  `sudo journalctl -u tailscaled --no-pager -n 50`. Verify no other VPN
  is bound to the same port.
- **Auth failure** — expired auth key (regenerate in admin console), SSO
  misconfiguration, or device limit reached.

## References

- `references/platform-install-commands.md` — exhaustive per-platform
  install commands
- `references/post-install-checklist.md` — verification steps, auth key
  options, common flags, and first-run troubleshooting
