#!/usr/bin/env bash
set -e

# Use the correct SHA256 library
SHA256=$(which sha256 || which sha256sum)

# For versioning
getCurrCommit() {
  echo `git rev-parse --short HEAD | tr -d "[ \r\n\']"`
}

# For versioning
getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

# Remove any previous builds that may have failed
[ -e "./build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build"

# Build 'algorand'
echo "Building 'algorand'..."
gox -ldflags="-s -X github.com/lvnacapital/algorand/cmd.version=$(getCurrTag)
  -X github.com/lvnacapital/algorand/cmd.commit=$(getCurrCommit)" \
  -osarch "darwin/amd64 linux/amd64 windows/amd64" -output="./build/{{.OS}}/{{.Arch}}/algorand"

# look through each os/arch/file and generate an SHA256 for each
echo "Generating SHA256 hashes..."
for os in $(ls ./build); do
  for arch in $(ls ./build/${os}); do
    for file in $(ls ./build/${os}/${arch}); do
      cat "./build/${os}/${arch}/${file}" | ${SHA256} | awk '{print $1}' >> "./build/${os}/${arch}/${file}.sha256"
    done
  done
done