#!/bin/bash
# mkQR 综合测试脚本

MKQR="/home/user/mkQR/build/mkqr"
TEST_DIR="/tmp/mkqr-tests"
PASS=0
FAIL=0

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 测试函数
test_case() {
    local name="$1"
    local expected_exit="$2"
    shift 2
    local cmd="$@"

    echo -n "  $name... "

    output=$(eval "$cmd" 2>&1)
    actual_exit=$?

    if [ "$actual_exit" -eq "$expected_exit" ]; then
        echo -e "${GREEN}PASS${NC}"
        ((PASS++))
        return 0
    else
        echo -e "${RED}FAIL${NC} (expected exit $expected_exit, got $actual_exit)"
        echo "    Command: $cmd"
        echo "    Output: ${output:0:100}..."
        ((FAIL++))
        return 1
    fi
}

# 设置测试环境
setup() {
    rm -rf "$TEST_DIR"
    mkdir -p "$TEST_DIR"
    cd "$TEST_DIR"
}

# 清理
cleanup() {
    rm -rf "$TEST_DIR"
}

echo "=========================================="
echo "        mkQR 综合测试套件"
echo "=========================================="
echo ""

setup

###########################################
echo -e "${YELLOW}[1] 基本子命令测试${NC}"
###########################################

test_case "text 子命令" 0 "$MKQR text 'Hello World' -q -o test_text.png"
test_case "url 子命令" 0 "$MKQR url 'github.com' -q -o test_url.png"
test_case "wifi 子命令" 0 "$MKQR wifi -s 'TestNet' -p 'pass123' -q -o test_wifi.png"
test_case "email 子命令" 0 "$MKQR email 'test@example.com' -q -o test_email.png"
test_case "phone 子命令" 0 "$MKQR phone '+8613800138000' -q -o test_phone.png"
test_case "sms 子命令" 0 "$MKQR sms '+123' -b 'Hello' -q -o test_sms.png"
test_case "geo 子命令" 0 "$MKQR geo --lat 39.9 --lng 116.4 -q -o test_geo.png"
test_case "otp 子命令" 0 "$MKQR otp -s 'JBSWY3DPEHPK3PXP' -i 'GitHub' -a 'user' -q -o test_otp.png"
test_case "vcard 子命令" 0 "$MKQR vcard -f 'John' --last 'Doe' -p '+123' -q -o test_vcard.png"

echo ""
###########################################
echo -e "${YELLOW}[2] 自动检测测试${NC}"
###########################################

test_case "自动检测 URL (https)" 0 "$MKQR 'https://github.com' -q -o auto_url.png"
test_case "自动检测 URL (无协议)" 0 "$MKQR 'github.com' -q -o auto_url2.png"
test_case "自动检测 vmess 协议" 0 "$MKQR 'vmess://eyJhZGQiOiIxLjEuMS4xIn0=' -q -o auto_vmess.png"
test_case "自动检测 vless 协议" 0 "$MKQR 'vless://uuid@host:443' -q -o auto_vless.png"
test_case "自动检测 ss 协议" 0 "$MKQR 'ss://YWVzLTI1Ni1nY206cGFzc0AxLjEuMS4xOjEwODA=' -q -o auto_ss.png"
test_case "自动检测纯文本" 0 "$MKQR 'Hello World' -q -o auto_text.png"

echo ""
###########################################
echo -e "${YELLOW}[3] 边界条件测试${NC}"
###########################################

# 空输入
test_case "空字符串输入" 1 "$MKQR '' -q"

# 超长输入 (QR 码有容量限制，约 3000 字符)
LONG_TEXT=$(python3 -c "print('A' * 2000)")
test_case "长文本 (2000字符)" 0 "$MKQR '$LONG_TEXT' -q -o long.png"

VERY_LONG=$(python3 -c "print('A' * 5000)")
test_case "超长文本 (5000字符，应失败)" 1 "$MKQR '$VERY_LONG' -q -o verylong.png"

# 特殊字符
test_case "特殊字符 (引号)" 0 "$MKQR 'Hello \"World\"' -q -o special1.png"
test_case "特殊字符 (反斜杠)" 0 "$MKQR 'path\\\\to\\\\file' -q -o special2.png"
test_case "Unicode 中文" 0 "$MKQR '你好世界' -q -o unicode_cn.png"
test_case "Unicode 日文" 0 "$MKQR 'こんにちは' -q -o unicode_jp.png"
test_case "Unicode emoji" 0 "$MKQR '🎉🔥💯' -q -o unicode_emoji.png"
test_case "换行符" 0 "$MKQR 'Line1\nLine2' -q -o newline.png"

# WiFi 特殊情况
test_case "WiFi 无密码" 0 "$MKQR wifi -s 'OpenNet' -e nopass -q -o wifi_open.png"
test_case "WiFi 隐藏网络" 0 "$MKQR wifi -s 'HiddenNet' -p 'pass' --hidden -q -o wifi_hidden.png"
test_case "WiFi SSID含特殊字符" 0 "$MKQR wifi -s 'My;Network:Name' -p 'pass' -q -o wifi_special.png"

echo ""
###########################################
echo -e "${YELLOW}[4] 错误处理测试${NC}"
###########################################

test_case "缺少必需参数 (wifi无ssid)" 1 "$MKQR wifi -p 'pass' -q"
test_case "缺少必需参数 (otp无secret)" 1 "$MKQR otp -i 'GitHub' -a 'user' -q"
test_case "缺少必需参数 (geo无lat)" 1 "$MKQR geo --lng 116.4 -q"
test_case "无效的纠错级别" 1 "$MKQR 'test' -l X -q"
test_case "url缺少参数" 1 "$MKQR url -q"
test_case "email缺少参数" 1 "$MKQR email -q"

echo ""
###########################################
echo -e "${YELLOW}[5] 管道输入测试${NC}"
###########################################

test_case "echo 管道输入" 0 "echo 'Hello from pipe' | $MKQR -q -o pipe1.png"
test_case "cat 管道输入" 0 "echo 'https://example.com' | $MKQR -q -o pipe2.png"
test_case "多行管道 (只取第一行)" 0 "printf 'Line1\nLine2\nLine3' | $MKQR -q -o pipe_multi.png"

echo ""
###########################################
echo -e "${YELLOW}[6] 批量处理测试${NC}"
###########################################

# 创建测试文件
echo -e "https://github.com\nhttps://google.com\nhttps://example.com" > batch_input.txt
test_case "批量处理 (3个URL)" 0 "$MKQR batch batch_input.txt -O ./batch_out/ -q"

# 检查生成的文件数量
BATCH_COUNT=$(ls ./batch_out/*.png 2>/dev/null | wc -l)
if [ "$BATCH_COUNT" -eq 3 ]; then
    echo -e "  批量文件数量检查... ${GREEN}PASS${NC} (3 files)"
    ((PASS++))
else
    echo -e "  批量文件数量检查... ${RED}FAIL${NC} (expected 3, got $BATCH_COUNT)"
    ((FAIL++))
fi

# 空行和注释
echo -e "# This is a comment\nhttps://a.com\n\nhttps://b.com\n# Another comment" > batch_comments.txt
test_case "批量处理 (跳过注释和空行)" 0 "$MKQR batch batch_comments.txt -O ./batch_out2/ -q"

BATCH_COUNT2=$(ls ./batch_out2/*.png 2>/dev/null | wc -l)
if [ "$BATCH_COUNT2" -eq 2 ]; then
    echo -e "  批量跳过注释检查... ${GREEN}PASS${NC} (2 files)"
    ((PASS++))
else
    echo -e "  批量跳过注释检查... ${RED}FAIL${NC} (expected 2, got $BATCH_COUNT2)"
    ((FAIL++))
fi

# stdin 批量
test_case "批量从stdin读取" 0 "echo -e 'url1\nurl2' | $MKQR batch - -O ./batch_stdin/ -q"

echo ""
###########################################
echo -e "${YELLOW}[7] 输出选项测试${NC}"
###########################################

test_case "终端输出 (默认)" 0 "$MKQR 'test' -q --small"
test_case "终端输出反色" 0 "$MKQR 'test' -q --invert --small"
test_case "PNG 输出" 0 "$MKQR 'test' -q -o output.png && file output.png | grep -q PNG"
test_case "自定义尺寸 (512)" 0 "$MKQR 'test' -q -o size512.png --size 512"
test_case "纠错级别 L" 0 "$MKQR 'test' -q -o level_l.png -l L"
test_case "纠错级别 H" 0 "$MKQR 'test' -q -o level_h.png -l H"

# 检查文件大小差异 (H 级别应该比 L 大)
SIZE_L=$(stat -c%s level_l.png 2>/dev/null || stat -f%z level_l.png 2>/dev/null)
SIZE_H=$(stat -c%s level_h.png 2>/dev/null || stat -f%z level_h.png 2>/dev/null)
if [ "$SIZE_H" -gt "$SIZE_L" ]; then
    echo -e "  纠错级别影响大小... ${GREEN}PASS${NC} (L=$SIZE_L, H=$SIZE_H)"
    ((PASS++))
else
    echo -e "  纠错级别影响大小... ${YELLOW}WARN${NC} (L=$SIZE_L, H=$SIZE_H) 可能相同"
fi

echo ""
###########################################
echo -e "${YELLOW}[8] 版本和帮助测试${NC}"
###########################################

test_case "显示版本" 0 "$MKQR --version"
test_case "显示帮助" 0 "$MKQR --help"
test_case "子命令帮助 (wifi)" 0 "$MKQR wifi --help"
test_case "子命令帮助 (otp)" 0 "$MKQR otp --help"

echo ""
###########################################
echo -e "${YELLOW}[9] 代理协议专项测试${NC}"
###########################################

test_case "vmess 协议" 0 "$MKQR 'vmess://eyJ2IjoiMiIsInBzIjoibm9kZSIsImFkZCI6IjEuMS4xLjEiLCJwb3J0Ijo0NDN9' -q -o vmess.png"
test_case "vless 协议" 0 "$MKQR 'vless://uuid@example.com:443?type=tcp#node' -q -o vless.png"
test_case "trojan 协议" 0 "$MKQR 'trojan://password@example.com:443#node' -q -o trojan.png"
test_case "ss 协议" 0 "$MKQR 'ss://YWVzLTEyOC1nY206dGVzdA==@1.1.1.1:8388#node' -q -o ss.png"
test_case "hysteria2 协议" 0 "$MKQR 'hysteria2://auth@example.com:443' -q -o hy2.png"

echo ""
echo "=========================================="
echo "              测试结果汇总"
echo "=========================================="
echo ""
echo -e "  通过: ${GREEN}$PASS${NC}"
echo -e "  失败: ${RED}$FAIL${NC}"
echo ""

TOTAL=$((PASS + FAIL))
RATE=$((PASS * 100 / TOTAL))
echo "  通过率: $RATE%"
echo ""

if [ "$FAIL" -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}有 $FAIL 个测试失败${NC}"
    exit 1
fi
