#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

if [ -z "$BUILD_ENV" ]; then
  # Make the build environment and run ourselves from there
  echo Creating go build container
  docker build \
    -q \
    -t=gobuildenv \
    -f Dockerfile-build . || exit 1

  echo Executing build
  docker run \
    -e EXUID=`id -u` \
    -e SKIP_DEP_INSTALL \
    -v "$DIR:/mainsrc/src/github.com/traherom/godyn:rw" \
    -v "$DIR/vendor:/extsrc:rw" \
    gobuildenv || exit 1

  exit
fi

# Actual build actions
if [ -z "$SKIP_DEP_INSTALL" ]; then
  echo Ensuring all packages are installed
  go get || exit 1
else
  echo Skipping dependency installation \(SKIP_DEP_INSTALL set\)
fi

echo Compiling
go build || exit 1

echo Changing files to be user-accessible
chown -hR $EXUID:$EXUID . "$DIR" || exit 1

