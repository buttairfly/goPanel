#!/bin/bash

# color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
LIGHT_BLUE='\033[1;34m'
NC='\033[0m' # No Color

BINARY=gopanel
VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINARY} ./src ; then
    echo "build  ${BINARY}"
    if scp ${BINARY} pi@ledpix:~/goPanel ; then
        scp ./config/* pi@ledpix:~/
        echo "deploy ${BINARY}"
        exit 0
    elif scp ${BINARY} pi@ledpix.fritz.box:~/goPanel ; then
        scp ./config/* pi@ledpix.fritz.box:~/config
        echo "deploy ${BINARY}"
        exit 0
    else
        echo "deploy failed" >&2
        exit 1
    fi
else
    echo "build  failed" >&2
    exit 1
fi
