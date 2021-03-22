#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
SCRIPT_DIR="$PROJECT_DIR/scripts"

VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`

FRONTEND_DIR="$PROJECT_DIR/ledpanel"
FRONTEND_VERSION="$( jq .version ${FRONTEND_DIR}/package.json )"

# color codes
source "$SCRIPT_DIR/color.sh"

# commands
source "$SCRIPT_DIR/commands.sh"

echo -e "${GREEN}frontend${NC}: compiled at ${BLUE}${DATE}${NC} with version ${FRONTEND_VERSION} git ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD_FRONTEND $PROJECT_DIR; then
    exit 0
else
    echo "build frontend failed" >&2
    exit 1
fi
