#!/usr/bin/env bash
# edited from deno_install
set -e

if ! command -v unzip >/dev/null; then
	echo "错误: 需要unzip来安装ipgw。" 1>&2
	exit 1
fi

if [ "$OS" = "Windows_NT" ]; then
	target="windows-amd64"
else
	case $(uname -sm) in
	"Darwin x86_64") target="darwin-amd64" ;;
	"Darwin arm64") target="darwin-arm64" ;;
	"FreeBSD x86_64") target="freebsd-amd64" ;;
	"FreeBSD i386") target="freebsd-386" ;;
	"Linux x86_64") target="linux-amd64" ;;
	"Linux arm") target="linux-arm" ;;
	"Linux i386") target="linux-386" ;;
	"Linux mips64") target="linux-mips64" ;;
	"Linux mips64le") target="linux-mips64le" ;;
	esac
fi

download_url="https://github.com/neucn/ipgw/releases/latest/download/ipgw-${target}.zip"

bin_dir="/usr/local/bin"
target_path="$bin_dir/ipgw"

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$target_path.zip" "$download_url"
unzip -d "$bin_dir" -o "$target_path.zip"
chmod +x "$target_path"
rm "$target_path.zip"

echo "ipgw成功安装到了 $target_path"
if command -v ipgw >/dev/null; then
	echo "运行 'ipgw --help' 来开始吧！"
else
	case $SHELL in
	/bin/zsh) shell_profile=".zshrc" ;;
	*) shell_profile=".bash_profile" ;;
	esac
	echo "请手动添加路径到 \$HOME/$shell_profile (或类似)"
	echo "  export PATH=\"$bin_dir:\$PATH\""
	echo "运行 '$target_path --help' 来开始吧！"
fi