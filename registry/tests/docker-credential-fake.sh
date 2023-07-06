#!/usr/bin/env bash

value=$(cat /dev/stdin)

case "${value}" in
  docker.io/*)
    username="user1"
    password="pass1"
    ;;
  example.com/*)
    username="user2"
    password="pass2"
    ;;
  *)
    echo "unknown"
    exit 1
    ;;
esac

echo "{\"Username\": \"$username\", \"Secret\": \"$password\"}"

