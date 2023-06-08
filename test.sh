#!/bin/sh

PRIVATE_KEY=$1
REMOTE_HOST=$2
DAX_ADDRESS=$3

if [ -z "${PRIVATE_KEY}" ]; then
  usage
  exit 1
fi

if [ -z "${REMOTE_HOST}" ]; then
  usage
  exit 1
fi

if [ -z "${DAX_ADDRESS}" ]; then
  usage
  exit 1
fi

set -x
GOOS=linux GOARCH=amd64 go build -o daxtest .
tar czf - -C ./ daxtest | ssh -i "${PRIVATE_KEY}" ec2-user@"${REMOTE_HOST}" "tar zxvf - -C ~/ && ./daxtest -ep ${DAX_ADDRESS}"

usage() {
  echo "usage: $0 private-key remote-host dax-address"
}