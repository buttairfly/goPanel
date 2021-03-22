#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
SCRIPT_DIR="$PROJECT_DIR/scripts"
BINDIR="${GOPATH}/bin"
PACKAGE="gopanel"
BINARY="${PACKAGE}-$( uname -p )"

VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`
ENV=''

# CONFIG_PATH="${PROJECT_DIR}/internal/config/testdata/main.composed.print.config.yaml"
CONFIG_PATH="${PROJECT_DIR}/config/main.composed.config.yaml"

# color codes
source "$SCRIPT_DIR/color.sh"

# commands
source "$SCRIPT_DIR/commands.sh"

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD_BACKEND "$PROJECT_DIR" "$VERSION" "$DATE"; then
    echo -e "build  ${BLUE}${BINARY}${NC} ${PROJECT_DIR}/${BINARY}"
    ${BINARY} -config ${CONFIG_PATH}
else
    echo -e "build  ${RED}failed${NC}"
    exit 1
fi
