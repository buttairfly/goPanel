#!/bin/bash

BINARY=gopanel
VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`
echo "${BINARY}: compiled at ${DATE} with version ${VERSION}"

if env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINARY} ./src ; then
    echo "build  ${BINARY}"
    if scp ${BINARY} pi@ledpix:~/goPanel ; then
        echo "deploy ${BINARY}"
        exit 0
    elif scp ${BINARY} pi@ledpix.fritz.box:~/goPanel ; then
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
