#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
BINARY="gopanel"

VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`
ENV='env GOOS=linux GOARCH=arm GOARM=5'

# color codes
RED="\033[0;31m"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
LIGHT_BLUE='\033[1;34m'
NC='\033[0m' # No Color


# commands
function COPY {
  rsync -acE --progress $1 $2
}
function BUILD {
  return $(`${ENV} go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${PROJECT_DIR}/${BINARY} ${PROJECT_DIR}/cmd/${BINARY}`)
}

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD; then
    echo -e "build  ${BLUE}${BINARY}${NC}"
    if COPY ${PROJECT_DIR}/${BINARY} pi@ledpix:~ ; then
        echo -e "deploy ${BLUE}${BINARY}${NC}"
        COPY ${PROJECT_DIR}/config/ pi@ledpix:~/config
        echo -e "deploy ${BLUE}config folder${NC}"
        exit 0
    else
        echo -e "deploy ${RED}failed${NC}"
        exit 1
    fi
else
    echo -e "build  ${RED}failed${NC}"
    exit 1
fi
