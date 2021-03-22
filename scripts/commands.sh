#!/bin/bash

export BINDIR="${GOPATH}/bin"
export PACKAGE="gopanel"
export BINARY="${PACKAGE}-arm"

function COPY {
  ORIGIN="$1"
  TARGET="$2"
  rsync -acE --progress $ORIGIN $TARGET
}

function BUILD_BACKEND {
  PROJECT_DIR="$1"
  VERSION="$2"
  DATE="$3"
  ENV="$4"
  return $(`${ENV} go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINDIR}/${BINARY} ${PROJECT_DIR}/cmd/${PACKAGE}`)
}

function SSH {
  ssh -t pi@ledpix $@
}

function BUILD_FRONTEND {
  set -e
  PROJECT_DIR="$1"
  TARGET_DIR="$2"
  echo ${PROJECT_DIR}
  cd ${PROJECT_DIR}/ledpanel
  yarn build
  cd -
}
