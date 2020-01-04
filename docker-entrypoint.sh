#!/bin/sh
set -e
./barometer init --use-env-vars -vv
./barometer serve --port=$PORT -vv
