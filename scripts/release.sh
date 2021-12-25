#!/usr/bin/env bash

name=$1
targets_dir="$PWD/$2"
release_dir="$PWD/$3"

echo "start: 解压 $targets_dir 到 $release_dir."
mkdir -p "$release_dir"
cd "$targets_dir" || exit

targets=$(ls "$targets_dir")
for target in $targets
do
    echo "解压中 $target..."
    cd "$target" && zip -q "$name-$target.zip" ./* && mv "$name-$target.zip" "$release_dir" && cd ..
done

echo "解压完成."