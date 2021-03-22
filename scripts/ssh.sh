#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
SCRIPT_DIR="$PROJECT_DIR/scripts"

# commands
source "$SCRIPT_DIR/commands.sh"

if SSH; then
    exit 0
else
    echo "ssh failed" >&2
    exit 1
fi
