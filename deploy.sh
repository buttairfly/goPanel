#!/bin/bash

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
  return $(`${ENV} go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINARY} ./src`)
}

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD; then
    echo -e "build  ${GREEN}${BINARY}${NC}"
    if COPY ./${BINARY} pi@ledpix:~/${BINARY} ; then
        COPY ./config/ pi@ledpix:~/config
        echo -e "deploy ${GREEN}${BINARY}${NC}"
        exit 0
    elif COPY ./${BINARY} pi@ledpix.fritz.box:~/${BINARY} ; then
        COPY ./config/ pi@ledpix.fritz.box:~/config
        echo -e "deploy ${GREEN}${BINARY}${NC}"
        exit 0
    else
        echo -e "deploy ${RED}failed${NC}"
        exit 1
    fi
else
    echo -e "build  ${RED}failed${NC}"
    exit 1
fi
