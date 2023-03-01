#!/bin/bash

export SLUG=ghcr.io/awakari/router
export VERSION=latest
docker tag awakari/router "${SLUG}":"${VERSION}"
docker push "${SLUG}":"${VERSION}"
