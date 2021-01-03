#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
BINDIR="${GOPATH}/bin"
PACKAGE="gopanel"
BINARY="${PACKAGE}-arm"

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
  return $(`${ENV} go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINDIR}/${BINARY} ${PROJECT_DIR}/cmd/${PACKAGE}`)
}
function SSH {
  ssh -t pi@ledpix $@
}

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD; then
    echo -e "build  ${BLUE}${BINARY}${NC}"

    echo -e "stop service ${BLUE}${PACKAGE}${NC}, when available"
    SSH "sudo systemctl stop ${PACKAGE}.service"

    if COPY ${BINDIR}/${BINARY} pi@ledpix:~ ; then
        echo -e "deploy ${BLUE}${BINARY}${NC}"

        SSH test ! -d ./config && SSH mkdir ./config && echo "~/config folder created"

        COPY ${PROJECT_DIR}/config/ pi@ledpix:~/config
        echo -e "deploy ${BLUE}config folder${NC}"


        echo -e "start ${GREEN}${BINARY}${NC} ${LIGHT_BLUE}${VERSION}${NC}"
        SSH ./${BINARY}

        exit 0
    else
        echo -e "deploy ${RED}failed${NC}"
        exit 1
    fi
else
    echo -e "build  ${RED}failed${NC}"
    exit 1
fi
