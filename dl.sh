#!/bin/bash

TAG=$(curl -sX GET https://api.github.com/repos/iLogtail/config-server-ui/releases/latest \
  | awk '/tag_name/{print $4;exit}' FS='[""]')

TARGET_DIR="router/statics"
mkdir -p "$TARGET_DIR"
rm -rf "$TARGET_DIR"/*

ARCHIVE="config-server-ui-${TAG}.tar.gz"
if ! curl -o "$ARCHIVE" -L "https://github.com/iLogtail/config-server-ui/releases/download/${TAG}/${ARCHIVE}"; then
    echo "❌ failed to download ${ARCHIVE}"
    exit 1
fi

TMP_DIR=$(mktemp -d)
if ! tar -zxf "$ARCHIVE" -C "$TMP_DIR"; then
    echo "❌ failed to untar ${ARCHIVE}"
    exit 2
fi

cp -r "$TMP_DIR"/pub/* "$TARGET_DIR"/

rm -rf "$TMP_DIR" "$ARCHIVE"