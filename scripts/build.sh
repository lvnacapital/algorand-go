#!/usr/bin/env bash
set -e
set -u
set -o pipefail

GIT=git
TR=tr
RMDIR='rm -rf'
AWK=awk
DIRNAME=dirname
LS=ls
CAT=cat
WHICH=which
UNAME='uname -s'
GO=go
GOX=gox
PLATFORMS='darwin/amd64 linux/amd64 windows/amd64'
BUILDDIR=${PWD}/build
SCRIPTDIR=$(${DIRNAME} ${0})
BINARY=algorand
USER=lvnacapital
PACKAGE=github.com/${USER}/${BINARY}

while getopts ':b:p:o:' opt; do
  case ${opt} in
    b)
      BUILDDIR=$OPTARG
      ;;
    p)
      PACKAGE=$OPTARG
      ;;
    o)
      PLATFORMS=$OPTARG
      ;;
    \?)
      echo "Usage: $(basename $0) -b /path/to/build" >&2
      exit 1
      ;;
    :)
      echo "Invalid option: -$OPTARG requires an argument" 1>&2
      exit 1
      ;;
  esac
done
shift $(($OPTIND - 1))

if [ "$(${UNAME})" == "Darwin" ]; then
  SHA256=sha256      
elif [ "$(expr substr $(${UNAME}) 1 5)" == "Linux" ]; then
  SHA256=sha256sum
elif [ "$(expr substr $(${UNAME}) 1 10)" == "MINGW32_NT" ]; then
  SHA256=sha256sum
elif [ "$(expr substr $(${UNAME}) 1 10)" == "MINGW64_NT" ]; then
  SHA256=sha256sum
fi

getCurrCommit() {
  echo `${GIT} rev-parse --short HEAD | ${TR} -d "[ \r\n\']"`
}

getCurrTag() {
  echo `${GIT} describe --always --tags --abbrev=0 | ${TR} -d "[v\r\n]"`
}

[ -e "${BUILDDIR}" ] && \
  echo "Cleaning up old builds..." && \
  ${RMDIR} "${BUILDDIR}"

echo "Building '${BINARY}'..."
${GOX} -ldflags="-s -X ${PACKAGE}/cmd.version=$(getCurrTag)
  -X ${PACKAGE}/cmd.commit=$(getCurrCommit)" \
  -osarch "${PLATFORMS}" -output="${BUILDDIR}/{{.OS}}/{{.Arch}}/${BINARY}"

echo 'Generating SHA256 hashes...'
for os in $(${LS} ${BUILDDIR}); do
  for arch in $(${LS} ${BUILDDIR}/${os}); do
    for file in $(${LS} ${BUILDDIR}/${os}/${arch}); do
      ${CAT} "${BUILDDIR}/${os}/${arch}/${file}" | ${SHA256} | ${AWK} "{ print \$1 \"  ${file}\" }" >> "${BUILDDIR}/${os}/${arch}/${file}.sha256"
    done
  done
done
