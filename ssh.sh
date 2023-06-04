#!/bin/sh

set -x

PRIVATE_KEY=$1
REMOTE_HOST=$2

if [ -z "${REMOTE_HOST}" ]; then
  echo "usage: $0 private-key remote-host"
fi

tar czf - -C ./ setup.sh main.go go.mod | ssh -i "${PRIVATE_KEY}" ec2-user@"${REMOTE_HOST}" 'tar zxvf - -C ~/'
ssh -i "${PRIVATE_KEY}" ec2-user@"${REMOTE_HOST}"
