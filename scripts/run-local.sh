#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
BINDIR="${GOPATH}/bin"
PACKAGE="gopanel"
BINARY="${PACKAGE}-$( uname -p )"

VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`

# color codes
RED="\033[0;31m"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
LIGHT_BLUE='\033[1;34m'
NC='\033[0m' # No Color


# commands
function BUILD {
  return $(`go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINDIR}/${BINARY} ${PROJECT_DIR}/cmd/${PACKAGE}`)
}

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD; then
    echo -e "build  ${BLUE}${BINARY}${NC} ${PROJECT_DIR}/${BINARY}"
    ${BINARY} -config ${PROJECT_DIR}/internal/config/testdata/main.composed.print.config.yaml
else
    echo -e "build  ${RED}failed${NC}"
    exit 1
fi
