#!/bin/bash

# set -x

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
SCRIPT_DIR="$PROJECT_DIR/scripts"
BINDIR="${GOPATH}/bin"
PACKAGE="gopanel"
BINARY="${PACKAGE}-arm"

ENV='env GOOS=linux GOARCH=arm GOARM=5'

# color codes
source "$SCRIPT_DIR/color.sh"

# commands
source "$SCRIPT_DIR/commands.sh"

echo -e "deploy service ${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD_BACKEND "$PROJECT_DIR" "$VERSION" "$DATE" "$ENV"; then
    echo -e "service build  ${BLUE}${BINARY}${NC}"
    if SSH test ! -d ./service ; then
      echo -e "servive folder ${RED}not available${NC}"
      exit 1
    fi
    if COPY ${BINDIR}/${BINARY} pi@ledpix:~/service/ ; then
        COPY ${PROJECT_DIR}/scripts/service.sh pi@ledpix:~/

        echo -e "service deploy ${BLUE}${BINARY}${NC}"
        if SSH test ! -d ./service/config ; then
          echo -e "service config folder ${RED}not available${NC} run service install beforehand"
          exit 1
        fi

        COPY ${PROJECT_DIR}/config/* pi@ledpix:~/service/config/
        echo -e "service deploy ${BLUE}config folder${NC}"

        echo -e "restart service ${BLUE}${PACKAGE}${NC}"
        SSH "sudo systemctl restart ${PACKAGE}.service"
        echo -e "restarted service ${BLUE}${PACKAGE}${NC}"
        exit 0
    else
        echo -e "service deploy ${RED}failed${NC}"
        exit 1
    fi
else
    echo -e "service build  ${RED}failed${NC}"
    exit 1
fi
