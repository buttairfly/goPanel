#!/bin/bash

if env GOOS=linux GOARCH=arm GOARM=5 go build ./src/gopanel/gopanel.go ; then
    if scp gopanel pi@ledpix:~ ; then
        exit 0
    elif scp gopanel pi@ledpix.fritz.box:~ ; then
        exit 0
    else
        echo "deploy failed" >&2
        exit 1
    fi
else
    echo "build failed" >&2
    exit 1
fi
