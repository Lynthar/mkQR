# mkQR

[![license](https://img.shields.io/github/license/Lynthar/mkQR)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/Lynthar/mkQR/ci.yml?branch=main&label=CI)](https://github.com/Lynthar/mkQR/actions/workflows/ci.yml)
[![release](https://img.shields.io/github/v/release/Lynthar/mkQR)](https://github.com/Lynthar/mkQR/releases)

Offline CLI QR generator with typed subcommands for WiFi, vCard, OTP, geo and calendar payloads. PNG, SVG, terminal

English | [简体中文](README.zh-CN.md)

Most QR tools take a string and draw it, which leaves the hard part to you: the
WiFi payload format, the vCard folding, the `otpauth://` parameter order. Get one
of those wrong and the scanner at the other end does nothing useful.

I gave each of them its own subcommand, so you pass `-s MyNet -p secret` instead
of hand-assembling `WIFI:S:MyNet;T:WPA;P:secret;;`. One static Go binary, no
runtime dependencies, no network access.

## Install

Grab a binary from [Releases](https://github.com/Lynthar/mkQR/releases) —
`mkqr-linux-amd64`, `mkqr-linux-arm64`, `mkqr-darwin-amd64`,
`mkqr-darwin-arm64` or `mkqr-windows-amd64.exe`, with `SHA256SUMS` alongside:

```bash
curl -LO https://github.com/Lynthar/mkQR/releases/latest/download/mkqr-darwin-arm64
chmod +x mkqr-darwin-arm64 && sudo mv mkqr-darwin-arm64 /usr/local/bin/mkqr
```

Or build it — needs Go 1.24.7 or newer:

```bash
go install github.com/Lynthar/mkQR/cmd/mkqr@latest
```

## Usage

```bash
mkqr "https://github.com"
echo "Hello" | mkqr
```

Without a subcommand it identifies the type itself — bare domains get `https://`
added, and ten proxy URL schemes are recognised. With one, the payload is built
for you:

| Subcommand | Builds |
|---|---|
| `wifi` | `WIFI:` network credentials — SSID, password, encryption, hidden |
| `vcard` | vCard contact — name, phone, email, organisation, address |
| `otp` | `otpauth://` for TOTP and HOTP authenticator apps |
| `email` · `phone` · `sms` | `mailto:`, `tel:` and `sms:` with prefilled fields |
| `geo` | `geo:` coordinates |
| `event` | Calendar event with start and end |
| `url` · `text` | Exactly what you typed |

```bash
mkqr wifi -s "MyNet" -p "pass" --encryption WPA
mkqr vcard -f "John" --last "Doe" -p "+1234567890" -e "john@example.com"
mkqr otp -s "JBSWY3DPEHPK3PXP" -i "GitHub" -a "user@example.com"
mkqr event -s "Sync" --start 2026-05-01T10:00:00Z --end 2026-05-01T11:00:00Z
```

Output is PNG, SVG, or Unicode blocks in the terminal, decided by `-o`. Batch
mode reads one payload per line and writes numbered PNGs:

```bash
mkqr "hello" -o qr.svg
mkqr "hello" -o qr.png --logo logo.png --fg blue
mkqr batch urls.txt -O ./out/ --prefix "node_"
```

Global flags: `-o/--output`, `--size` (256), `-l/--level` (L/M/Q/H, default M),
`--logo`, `--fg`, `--bg`, `--invert`, `--small`, `-q/--quiet`. Nine colour names
are accepted (`black`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`,
`magenta`, `transparent`), or any hex value; a logo defaults to 20% of the QR
width and is refused above 35%. There's no config file and no environment
variables.

## Limitations

- **Generate only.** It doesn't decode QR codes.
- **Batch mode writes PNG only**, with generated filenames — `-o` is ignored
  there, so you can't batch out SVG.
- **Logos only work for PNG.** SVG and terminal output refuse them outright.
- **`--invert` only affects terminal output.** For files, use `--fg` and `--bg`.
- **Contrast isn't checked.** Pick two similar colours and you'll get a QR code
  that nothing can scan.
- **Payload correctness is verified against the specs, not against scanners.**
  vCard folding, `otpauth://` parameter order and calendar time zones all have
  dialects in the wild.

## License

MIT — see [LICENSE](LICENSE). Copyright (c) 2025 Lynthar.
