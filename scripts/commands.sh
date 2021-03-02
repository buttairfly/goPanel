#!/bin/bash

function COPY {
  rsync -acE --progress $1 $2
}
function BUILD {
  return $(`${ENV} go build -ldflags "-X main.compileDate=${DATE} -X main.versionTag=${VERSION}" -o ${BINDIR}/${BINARY} ${PROJECT_DIR}/cmd/${PACKAGE}`)
}
function SSH {
  ssh -t pi@ledpix $@
}

function BUILD_FRONTEND {

  PROJECT_DIR="$1"
  return $(` pwd && cd $PROJECT_DIR/ledpanel && pwd `)
}
