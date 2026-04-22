# mkQR

A fast, flexible QR code generator for the command line.

## Features

- **Multiple data types**: WiFi, URLs, contacts (vCard), OTP/2FA, email, phone, SMS, geographic location
- **Multiple output formats**: PNG, SVG (scalable vector), and Unicode in the terminal
- **Logo embedding**: Composite a PNG/JPEG/GIF image at the QR center; error correction is auto-raised to level H so the code stays scannable
- **Cross-platform**: Linux, macOS, Windows — a single self-contained binary that works offline
- **Script-friendly**: Supports stdin, exit codes, quiet mode, and batch processing
- **Auto-detection**: Recognizes input type from the content itself, including proxy links (vmess/vless/ss/trojan/hysteria/tuic/...)
- **Importable Go library**: Payload formatters (WiFi/vCard/OTP/email/SMS/geo/phone) reusable via `pkg/encoder`

## Installation

### Pre-built Binary (recommended)

Download from the [Releases](https://github.com/Lynthar/mkQR/releases) page. Each release publishes self-contained binaries for Linux (amd64/arm64), macOS (amd64/arm64), and Windows (amd64), plus a `SHA256SUMS` file for integrity verification. Nothing else to install — copy the binary anywhere (including air-gapped machines) and run it.

### `go install`

If you have a Go toolchain (≥ 1.24):

```bash
go install github.com/Lynthar/mkQR/cmd/mkqr@latest
```

This puts `mkqr` in `$(go env GOPATH)/bin`. Note: binaries built this way report `version dev` — only the Makefile wires real version metadata via `-ldflags`.

### From Source

```bash
git clone https://github.com/Lynthar/mkQR.git
cd mkQR
make install      # needs `make`; on Windows use Git Bash / WSL, or substitute `go install ./cmd/mkqr`
```

## Usage

### Basic Usage

```bash
# Generate QR for any text (auto-detect type)
mkqr "https://github.com"
mkqr "vmess://eyJ..."

# Bare URLs without a scheme are automatically prefixed with https://
mkqr github.com

# Read from stdin (great for scripts)
echo "Hello World" | mkqr
cat proxy_link.txt | mkqr

# Save to file
mkqr "https://example.com" -o qr.png

# Show version
mkqr --version
```

### WiFi Network

```bash
mkqr wifi -s "NetworkName" -p "password"
mkqr wifi --ssid "Home WiFi" --password "secret" --encryption WPA
mkqr wifi -s "OpenNetwork" --encryption nopass
mkqr wifi -s "HiddenNetwork" -p "pass" --hidden
```

### Contact Card (vCard)

```bash
mkqr vcard -f "John" --last "Doe" -p "+1234567890" -e "john@example.com"
mkqr vcard --first "Jane" --last "Smith" --org "Acme Inc" --mobile "+1234567890"
```

### Two-Factor Authentication (OTP)

```bash
mkqr otp -s "JBSWY3DPEHPK3PXP" -i "GitHub" -a "user@example.com"
mkqr otp --secret "NBSWY3DPO5XXE3DE" --issuer "AWS" --account "myaccount" --digits 8
```

The secret must be valid base32 (letters A–Z plus digits 2–7; whitespace and hyphens are stripped).

### Email, Phone, SMS

```bash
mkqr email hello@example.com -s "Subject" -b "Message body"
mkqr phone +1234567890
mkqr sms +1234567890 -b "Hello!"
```

### Geographic Location

```bash
mkqr geo --lat 40.7128 --lng -74.0060
mkqr geo --lat 39.9042 --lng 116.4074 --query "Beijing"
```

### Batch Processing

```bash
# Generate QR codes from a file (one per line)
mkqr batch urls.txt -O ./qrcodes/
mkqr batch nodes.txt --output-dir ./out --prefix "node_"

# From stdin
cat links.txt | mkqr batch - -O ./output/
```

### Output Options

```bash
# Terminal display (default)
mkqr "text"

# Save to PNG file
mkqr "text" -o qr.png

# Save to SVG (scales losslessly, ideal for print or web)
mkqr "text" -o qr.svg

# Embed a logo in the center (PNG output only; auto-forces error correction H)
mkqr "https://example.com" -o branded.png --logo mylogo.png

# Invert colors (for dark terminals)
mkqr "text" --invert

# Compact mode (smaller display)
mkqr "text" --small

# Adjust size and error correction
mkqr "text" -o qr.png --size 512 --level H

# Quiet mode (no status messages)
mkqr "text" -q
```

### Output Formats

| Method | Format | Location |
|--------|--------|----------|
| `mkqr "text"` | Unicode blocks (██, ▀, ▄) | Terminal (stdout) |
| `mkqr "text" -o file.png` | PNG image | Specified file path |
| `mkqr "text" -o file.svg` | SVG vector | Specified file path |
| `mkqr "text" -o file.png --logo logo.png` | PNG with centered logo | Specified file path |
| `mkqr batch file.txt -O ./dir/` | PNG images (one per line) | Specified directory |

Notes:

- Default pixel size is **256** (adjustable with `--size`); default error correction is **M** (adjustable with `--level L/M/Q/H`).
- SVG output sets `shape-rendering="crispEdges"` so scanners see hard module edges at any zoom.
- Logo images accept PNG, JPEG, or GIF; the logo is scaled to 20% of the QR edge and given a small white halo. Level H is forced while a logo is in use (a one-line note goes to stderr unless `--quiet`).
- Only `.png` and `.svg` are accepted as `-o` extensions; anything else (e.g. `.jpg`) is rejected up front rather than silently written as PNG.

## Supported Types

| Type | Command | Example |
|------|---------|---------|
| Text | `mkqr text` | `mkqr text "Hello"` |
| URL | `mkqr url` | `mkqr url github.com` |
| WiFi | `mkqr wifi` | `mkqr wifi -s "SSID" -p "pass"` |
| Contact | `mkqr vcard` | `mkqr vcard -f "John" -p "+123"` |
| Email | `mkqr email` | `mkqr email user@example.com` |
| Phone | `mkqr phone` | `mkqr phone +1234567890` |
| SMS | `mkqr sms` | `mkqr sms +123 -b "Hi"` |
| OTP/2FA | `mkqr otp` | `mkqr otp -s "SECRET" -i "App" -a "user"` |
| Location | `mkqr geo` | `mkqr geo --lat 40.71 --lng -74.00` |
| Batch | `mkqr batch` | `mkqr batch file.txt -O ./out/` |

## Integration with Scripts

mkQR is designed to work seamlessly with shell scripts:

```bash
#!/bin/bash
# Example: Generate QR for a proxy node

NODE_LINK="vmess://eyJhZGQiOi..."

# Display in terminal
mkqr "$NODE_LINK"

# Or save to file
mkqr "$NODE_LINK" -o node.png -q
```

## Building

```bash
# Build for current platform
make build

# Cross-compile for all platforms
make cross

# Run tests
make test
```

## Shell Completion

Completion scripts are generated on demand:

```bash
mkqr completion bash       > /etc/bash_completion.d/mkqr   # or source into your rc file
mkqr completion zsh        > "${fpath[1]}/_mkqr"
mkqr completion fish       > ~/.config/fish/completions/mkqr.fish
mkqr completion powershell | Out-String | Invoke-Expression
```

Run `mkqr completion <shell> --help` for shell-specific install hints.

## Use as a Go library

The payload encoders are exported under `github.com/Lynthar/mkQR/pkg/encoder` and can be imported directly if you want to produce the wire-format strings (WiFi, vCard, OTP, email, SMS, geo, phone) without a CLI shell-out.

```go
import "github.com/Lynthar/mkQR/pkg/encoder"

s := (&encoder.WiFi{SSID: "Home", Password: "secret", Encryption: encoder.WPA}).Encode()
// s == "WIFI:T:WPA;S:Home;P:secret;;"
```

Auto-detection is also exposed via `encoder.Detect()` / `encoder.DetectAndDescribe()`.

## Offline Usage

mkQR works completely offline — no network connection required at runtime. All QR code generation is done locally, and the binary is statically linked Go (no shared library dependencies).

### Option 1: Download from Releases (Easiest)

On a computer with internet access, download the appropriate binary and its checksum from the [Releases](https://github.com/Lynthar/mkQR/releases) page. Verify the file against `SHA256SUMS`, then copy it to your offline machine via USB or other removable media. No further steps required — just run the binary.

### Option 2: Build from Source on a Connected Machine

Clone the repo and run `make cross` (see [Building](#building)) to produce binaries for all five target platforms under `build/`. Copy the one that matches the offline machine via USB or similar. The result is functionally identical to a Releases download.

### Option 3: Copy Source with Vendored Dependencies

For strict air-gap scenarios where even a binary copy is impractical and you must compile *on* the offline machine:

```bash
# On a computer with internet access:
git clone https://github.com/Lynthar/mkQR.git
cd mkQR
go mod vendor      # fetches all deps into ./vendor
```

Copy the entire `mkQR/` directory (including `vendor/`) to the offline computer, then:

```bash
go build -mod=vendor -o mkqr ./cmd/mkqr
```

No network traffic at any point during the offline build.

## License

MIT License - see [LICENSE](LICENSE) for details.
