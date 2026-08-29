# mkQR

[![license](https://img.shields.io/github/license/Lynthar/mkQR)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/Lynthar/mkQR/ci.yml?branch=main&label=CI)](https://github.com/Lynthar/mkQR/actions/workflows/ci.yml)
[![release](https://img.shields.io/github/v/release/Lynthar/mkQR)](https://github.com/Lynthar/mkQR/releases)

离线命令行二维码生成器：WiFi / vCard / OTP / 日历等十种结构化载荷各有子命令，输出 PNG、SVG 与终端字符

[English](README.md) | 简体中文

多数二维码工具是「给一个字符串、画出来」，麻烦的部分留给你：WiFi 的载荷格式对不对、
vCard 怎么折行、`otpauth://` 的参数顺序。写错了，扫码的一端就什么都不会做。

我把这些各做成一个子命令，你写 `-s MyNet -p secret`，不用自己拼
`WIFI:S:MyNet;T:WPA;P:secret;;`。单个静态 Go 二进制，没有运行时依赖，不联网。

## 安装

从 [Releases](https://github.com/Lynthar/mkQR/releases) 取二进制——`mkqr-linux-amd64`、
`mkqr-linux-arm64`、`mkqr-darwin-amd64`、`mkqr-darwin-arm64`、`mkqr-windows-amd64.exe`，
同批带 `SHA256SUMS`：

```bash
curl -LO https://github.com/Lynthar/mkQR/releases/latest/download/mkqr-darwin-arm64
chmod +x mkqr-darwin-arm64 && sudo mv mkqr-darwin-arm64 /usr/local/bin/mkqr
```

或者自己编译，需要 Go 1.24.7 以上：

```bash
go install github.com/Lynthar/mkQR/cmd/mkqr@latest
```

## 用法

```bash
mkqr "https://github.com"
echo "Hello" | mkqr
```

不带子命令时它自己识别类型——裸域名会补上 `https://`，另外识别十种代理 URL 前缀。
带子命令时，载荷由它替你拼好：

| 子命令 | 生成的载荷 |
|---|---|
| `wifi` | `WIFI:` 网络凭据——SSID、密码、加密方式、是否隐藏 |
| `vcard` | vCard 名片——姓名、电话、邮箱、单位、地址 |
| `otp` | `otpauth://`，给 TOTP / HOTP 验证器 App 用 |
| `email` · `phone` · `sms` | 预填好字段的 `mailto:`、`tel:`、`sms:` |
| `geo` | `geo:` 坐标 |
| `event` | 带起止时间的日历事件 |
| `url` · `text` | 原样 |

```bash
mkqr wifi -s "MyNet" -p "pass" --encryption WPA
mkqr vcard -f "John" --last "Doe" -p "+1234567890" -e "john@example.com"
mkqr otp -s "JBSWY3DPEHPK3PXP" -i "GitHub" -a "user@example.com"
mkqr event -s "Sync" --start 2026-05-01T10:00:00Z --end 2026-05-01T11:00:00Z
```

输出是 PNG、SVG，或直接在终端里用 Unicode 方块画出来，由 `-o` 决定。批量模式读一个
每行一条载荷的文件，输出编号好的 PNG：

```bash
mkqr "hello" -o qr.svg
mkqr "hello" -o qr.png --logo logo.png --fg blue
mkqr batch urls.txt -O ./out/ --prefix "node_"
```

全局旗标：`-o/--output`、`--size`（256）、`-l/--level`（L/M/Q/H，默认 M）、`--logo`、
`--fg`、`--bg`、`--invert`、`--small`、`-q/--quiet`。颜色名支持九个（`black` `white`
`red` `green` `blue` `yellow` `cyan` `magenta` `transparent`），也可以直接给十六进制值；
logo 默认占二维码宽度的 20%，超过 35% 直接拒绝。没有配置文件，没有环境变量。

## 能力边界

- **只生成，不解码。**
- **批量模式只出 PNG**，文件名是自动生成的——那种情况下 `-o` 会被忽略，所以批量出不了 SVG。
- **logo 只对 PNG 有效**，SVG 和终端输出会直接报错。
- **`--invert` 只影响终端输出。** 输出到文件请用 `--fg` 和 `--bg`。
- **不校验对比度。** 挑两个相近的颜色，就会生成一个谁也扫不出来的二维码。
- **载荷正确性是按规范验的，不是按扫码端验的。** vCard 折行、`otpauth://` 参数顺序、
  日历时区表示，各家扫码端都有自己的方言。

## 许可证

MIT —— 见 [LICENSE](LICENSE)。Copyright (c) 2025 Lynthar。
