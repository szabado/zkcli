#!/usr/bin/env bash

set -eux

export GOPATH="$(pwd)/.gobuild"
SRCDIR="${GOPATH}/src/github.com/fJancsoSzabo/zkcli"

[[ -d "${GOPATH}" ]] && rm -rf ${GOPATH}
tar -czf archive.tar.gz .

mkdir -p ${GOPATH}/{src,pkg,bin}
mkdir -p ${SRCDIR}
cp archive.tar.gz  ${SRCDIR}
(
	cd ${SRCDIR}
	tar -xzf archive.tar.gz --strip-components 1
	dep ensure
	go install .
)
