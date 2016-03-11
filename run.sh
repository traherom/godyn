#!/bin/bash
DIR=$(readlink -m "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )")
cd "$DIR"

NAME=godyn

if docker ps | grep -q "$NAME"; then
  echo Stopping current instance of $NAME
  docker stop "$NAME" || exit 1
else
  echo $NAME not currently running
fi

if docker ps -a | grep -q "$NAME"; then
  echo Removing old image of $NAME
  docker rm "$NAME" || exit 1
else
  echo Old image of $NAME not found
fi

if [ -f "./build.sh" ]; then
  echo Building $NAME
  ./build.sh || exit 1
else 
  echo "No build script (build.sh) found for component"
fi

echo Creating container for $NAME
docker build -q -t "$NAME" -f Dockerfile-run . || exit 1

EXTRAs=""
if [[ "$1" == "-d" ]]; then
  EXTRAS="-d --restart=always"
fi

echo "Running $NAME ($EXTRAS)"
docker run -it \
  --name $NAME \
  -v "/etc/ssl/certs:/etc/ssl/certs:ro" \
  -e GODYN_SERVICE \
  -e GODYN_USER \
  -e GODYN_PW \
  $EXTRAS \
  $NAME || exit 1

