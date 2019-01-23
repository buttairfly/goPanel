#!/bin/bash

if env GOOS=linux GOARCH=arm GOARM=5 go build ./src/gopanel.go ; then
    if scp gopanel pi@ledpix:~/goPanel ; then
        exit 0
    elif scp gopanel pi@ledpix.fritz.box:~/goPanel ; then
        exit 0
    else
        echo "deploy failed" >&2
        exit 1
    fi
else
    echo "build failed" >&2
    exit 1
fi
