#!/bin/bash

if env GOOS=linux GOARCH=arm GOARM=5 go build ./src/gopanel/gopanel.go ; then
    scp gopanel pi@ledpix.fritz.box:~
    exit 0
else
    echo "build failed" >&2
    exit 1
fi
