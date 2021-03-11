#!/usr/bin/env bash
# edited from deno_install
set -e

if ! command -v unzip >/dev/null; then
	echo "Error: unzip is required to install ipgw." 1>&2
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

bin_dir="$HOME/neucn/bin"
target_path="$bin_dir/ipgw"

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$target_path.zip" "$download_url"
unzip -d "$bin_dir" -o "$target_path.zip"
chmod +x "$target_path"
rm "$target_path.zip"

echo "ipgw was installed successfully to $target_path"
if command -v ipgw >/dev/null; then
	echo "Run 'ipgw --help' to get started"
else
	case $SHELL in
	/bin/zsh) shell_profile=".zshrc" ;;
	*) shell_profile=".bash_profile" ;;
	esac
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export PATH=\"\$HOME/neucn/bin:\$PATH\""
	echo "Run '$target_path --help' to get started"
fi