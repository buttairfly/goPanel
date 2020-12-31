#!/bin/bash

# set -x

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
BINDIR="${GOPATH}/bin"
PACKAGE="gopanel"
BINARY="${PACKAGE}-arm"

VERSION=`git describe --always --dirty`
DATE=`date -u +%FT%T%z`
ENV='env GOOS=linux GOARCH=arm GOARM=5'

# color codes
RED="\033[0;31m"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
LIGHT_BLUE='\033[1;34m'
NC='\033[0m' # No Color


# commands
function COPY {
  rsync -acE --progress $1 $2
}
function BUILD {
  return $(`${ENV} go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINDIR}/${BINARY} ${PROJECT_DIR}/cmd/${PACKAGE}`)
}
function SSH {
  ssh -t pi@ledpix $@
}

echo -e "deploy service ${GREEN}${BINARY}${NC}: compiled at ${BLUE}${DATE}${NC} with version ${LIGHT_BLUE}${VERSION}${NC}"

if BUILD; then
    echo -e "service build  ${BLUE}${BINARY}${NC}"
    SSH test ! -d ./service && SSH mkdir ./service && echo servive folder created
    if COPY ${BINDIR}/${BINARY} pi@ledpix:~/service/ ; then
        COPY ${PROJECT_DIR}/scripts/service.sh pi@ledpix:~/

        echo -e "service deploy ${BLUE}${BINARY}${NC}"
        SSH test ! -d ./service/config && SSH mkdir ./service/config && echo servive/config folder created

        COPY ${PROJECT_DIR}/config/* pi@ledpix:~/service/config/
        echo -e "service deploy ${BLUE}config folder${NC}"

        COPY ${PROJECT_DIR}/scripts/${PACKAGE}.service pi@ledpix:~/service/
        echo -e "service deploy ${BLUE}${PACKAGE}.service systemd file${NC}"

        SSH sudo rsync -acE --progress ./service/${PACKAGE}.service /etc/systemd/system/${PACKAGE}.service
        echo -e "service deploy ${BLUE}systemd file${NC}"

        echo -e "enable and restart service ${BLUE}${PACKAGE}${NC}"
        SSH "sudo systemctl restart ${PACKAGE}.service"
        SSH "sudo systemctl enable ${PACKAGE}"
        echo -e "enabled and started service ${BLUE}${PACKAGE}${NC}"
        exit 0
    else
        echo -e "service deploy ${RED}failed${NC}"
        exit 1
    fi
else
    echo -e "service build  ${RED}failed${NC}"
    exit 1
fi
