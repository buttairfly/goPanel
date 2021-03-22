#!/bin/bash

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
SCRIPT_DIR="$PROJECT_DIR/scripts"

VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`
ENV='env GOOS=linux GOARCH=arm GOARM=5'

# color codes
source "$SCRIPT_DIR/color.sh"

# commands
source "$SCRIPT_DIR/commands.sh"

echo -e "${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD_BACKEND "$PROJECT_DIR" "$VERSION" "$DATE" "$ENV"; then
    echo -e "build  ${BLUE}${BINARY}${NC}"

    echo -e "stop service ${BLUE}${PACKAGE}${NC}, when available"
    SSH "sudo systemctl stop ${PACKAGE}.service"

    if SSH test "$( ps a | grep ${BINARY} | wc -l )" -ne "1" ; then
      echo -e "program ${BINARY} ${RED}already running${NC}"
      exit 1
    fi

    if source "$SCRIPT_DIR/build-frontend.sh" ; then
      echo -e "frontend ${RED}compile error${NC}"
      exit 1
    fi

    if COPY ${BINDIR}/${BINARY} pi@ledpix:~ ; then
        echo -e "deploy ${BLUE}${BINARY}${NC}"

        SSH test ! -d ./config && SSH mkdir ./config && echo "~/config folder created"

        SSH test ! -d ./ledpanel/build && SSH mkdir -p ./ledpanel/build && echo "~/ledpanel/build folder created"

        COPY ${PROJECT_DIR}/config/ pi@ledpix:~/config
        echo -e "deploy ${BLUE}config folder${NC}"

        COPY ${PROJECT_DIR}/config/ pi@ledpix:~/config
        echo -e "deploy ${BLUE}config folder${NC}"

        echo -e "start ${GREEN}${BINARY}${NC} ${LIGHT_BLUE}${VERSION}${NC}"
        SSH ./${BINARY}

        exit 0
    else
        echo -e "deploy ${RED}failed${NC}"
        exit 1
    fi
else
    echo -e "build  ${RED}failed${NC}"
    exit 1
fi
