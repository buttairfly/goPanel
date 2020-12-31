#!/bin/bash

if ssh -t pi@ledpix $@; then
    exit 0
else
    echo "ssh failed" >&2
    exit 1
fi
