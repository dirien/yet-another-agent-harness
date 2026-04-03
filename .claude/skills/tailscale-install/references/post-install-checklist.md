<!-- Copyright 2025-2026 Richard Shade. Licensed under Apache-2.0. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
# Post-Install Checklist

Verification steps and initial configuration after Tailscale installation.

## Version Verification

```bash
tailscale version
```

Confirm output shows a version string (e.g., `1.78.1`). If the command is
not found, the installation failed or the binary is not in PATH.

For WSL2, use `tailscale.exe version` instead.

## Daemon Status

### Linux (systemd)

```bash
sudo systemctl status tailscaled
```

Expected: `active (running)`. If not running:

```bash
sudo systemctl enable --now tailscaled
```

### Linux (OpenRC — Alpine)

```bash
sudo rc-service tailscale status
```

### macOS

The Tailscale app manages the daemon. Open Tailscale from Applications or
the menu bar. No `systemctl` equivalent.

### Windows

```powershell
Get-Service Tailscale
```

Expected: `Running`. The Tailscale Windows service starts automatically
after installation.

## Initial Authentication

### Interactive (Desktop/Laptop)

```bash
sudo tailscale up
```

This prints an authentication URL. Open it in a browser and sign in with
your identity provider. After authentication, the device appears in the
admin console.

### Headless (Server)

Generate an auth key in the Tailscale admin console under
Settings > Keys > Generate auth key, then:

```bash
sudo tailscale up --authkey=tskey-auth-<key>
```

### Auth Key Options

| Flag | Purpose |
| ---- | ------- |
| `--authkey=tskey-auth-<key>` | Authenticate without browser |
| Reusable key | Same key for multiple devices |
| One-time key | Single device, expires after first use |
| Ephemeral key | Device auto-removed when it goes offline |
| Pre-approved key | Skips admin approval step |

Generate keys in the admin console: Settings > Keys > Generate auth key.

## Common `tailscale up` Flags

| Flag | Purpose | Example |
| ---- | ------- | ------- |
| `--hostname=<name>` | Set MagicDNS hostname | `--hostname=web-prod-1` |
| `--advertise-tags=<tags>` | Apply ACL tags | `--advertise-tags=tag:server` |
| `--advertise-routes=<cidrs>` | Act as subnet router | `--advertise-routes=10.0.0.0/24` |
| `--accept-routes` | Accept routes from subnet routers | `--accept-routes` |
| `--ssh` | Enable Tailscale SSH | `--ssh` |
| `--exit-node=<peer>` | Use a peer as exit node | `--exit-node=exit-us-east` |
| `--advertise-exit-node` | Offer this device as exit node | `--advertise-exit-node` |
| `--operator=<user>` | Allow non-root user to manage | `--operator=$USER` |
| `--shields-up` | Block incoming connections | `--shields-up` |

## Network Verification

After authentication, verify connectivity:

```bash
# Check connection status and peer list
tailscale status

# Test connectivity to a specific peer
tailscale ping <peer-hostname-or-ip>

# Verify MagicDNS resolution
tailscale status | grep <expected-peer>
```

Expected: `tailscale status` shows the device as connected with a
100.x.y.z IP address and lists other peers in the tailnet.

## Common First-Run Issues

### Daemon Not Running

**Symptom:** `tailscale up` fails with "connection refused" or similar.

**Fix:**

```bash
sudo systemctl enable --now tailscaled
```

### Firewall Blocking UDP

**Symptom:** All connections relay through DERP (no direct connections).

**Fix:** Allow outbound UDP port 41641. Tailscale uses this for direct
WireGuard connections. Without it, traffic routes through DERP relay
servers (slower but functional).

### DNS Resolver Conflicts

**Symptom:** MagicDNS names do not resolve after connecting.

**Fix:** Check if another DNS resolver (systemd-resolved, dnsmasq,
NetworkManager) is overriding Tailscale DNS settings:

```bash
tailscale dns status
resolvectl status  # systemd-resolved
```

### Permission Denied

**Symptom:** `tailscale up` requires root but user is non-root.

**Fix:** Use the `--operator` flag during initial setup:

```bash
sudo tailscale up --operator=$USER
```

This allows the specified user to run `tailscale` commands without `sudo`.

### Auth Key Expired or Invalid

**Symptom:** `tailscale up --authkey=...` fails with authentication error.

**Fix:** Generate a new auth key in the admin console. Auth keys have
configurable expiration (default: 90 days). Check that the key has not
been revoked.
