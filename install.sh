#!/bin/sh

set -e

REPOSITORY="stnaoficial/mimic"
INSTALL_DIR="/usr/local/bin"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
	Linux)
		PLATFORM="linux"
		;;
	Darwin)
		PLATFORM="darwin"
		;;
	*)
		echo "Unsupported operating system"
		exit 1
		;;
esac

case "$ARCH" in
	x86_64)
		ARCHITECTURE="amd64"
		;;
	aarch64|arm64)
		ARCHITECTURE="arm64"
		;;
	*)
		echo "Unsupported architecture"
		exit 1
		;;
esac

VERSION="$(curl -fsSL "https://api.github.com/repos/$REPOSITORY/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4)"

ARCHIVE="mimic-$PLATFORM-$ARCHITECTURE.tar.gz"
URL="https://github.com/$REPOSITORY/releases/download/$VERSION/$ARCHIVE"

TEMP_DIR="$(mktemp -d)"

trap 'rm -rf "$TEMP_DIR"' EXIT

curl -fsSL "$URL" -o "$TEMP_DIR/$ARCHIVE"
tar -xzf "$TEMP_DIR/$ARCHIVE" -C "$TEMP_DIR"

sudo install "$TEMP_DIR/$PLATFORM-$ARCHITECTURE/mimic" "$INSTALL_DIR/mimic"

echo "Mimic $VERSION installed successfully"