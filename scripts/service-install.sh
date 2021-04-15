#!/bin/bash

# set -x

PROJECT_DIR="$(dirname "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
SCRIPT_DIR="$PROJECT_DIR/scripts"
PACKAGE="gopanel"

# color codes
source "$SCRIPT_DIR/color.sh"

# commands
source "$SCRIPT_DIR/commands.sh"

echo -e "service install"

if SSH test ! -d ./service/config ; then
    SSH mkdir -p ./service/config
fi

if source "$SCRIPT_DIR/deploy.sh"; then
    SSH sudo rsync -acE --progress ./service/${PACKAGE}.service /etc/systemd/system/${PACKAGE}.service
    echo -e "service deploy ${BLUE}systemd file${NC}"

    echo -e "enable and restart service ${BLUE}${PACKAGE}${NC}"
    SSH "sudo systemctl restart ${PACKAGE}.service"
    SSH "sudo systemctl enable ${PACKAGE}"
    echo -e "enabled and started service ${BLUE}${PACKAGE}${NC}"
    exit 0
else
    echo -e "service deploy  ${RED}failed${NC}"
    exit 1
fi
