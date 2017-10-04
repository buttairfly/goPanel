#!/bin/bash

if ssh pi@ledpix ; then
    exit 0
elif ssh pi@ledpix.fritz.box ; then
    exit 0
else
    echo "ssh failed" >&2
    exit 1
fi
