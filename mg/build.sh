#!/bin/bash

# Simple script for building mg tool. Includes injection
# of provenance information.

GITHASH=$(git rev-parse --short HEAD)
GITTAG=$(git tag -l --points-at HEAD)
GITBRANCH=$(git rev-parse --abbrev-ref HEAD)

echo "${GITHASH} ${GITTAG}"

go build -ldflags \
	"-X 'metagraf/pkg/mgver.GitHash=${GITHASH}' \
	-X 'metagraf/pkg/mgver.GitTag=${GITTAG}' \
	-X 'metagraf/pkg/mgver.GitBranch=${GITBRANCH}'" 
