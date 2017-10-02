#!/bin/bash

if ssh gopanel pi@ledpix:~ ; then
    exit 0
elif ssh gopanel pi@ledpix.fritz.box:~ ; then
    exit 0
else
    echo "ssh failed" >&2
    exit 1
fi
