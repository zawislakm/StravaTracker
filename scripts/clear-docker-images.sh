#!/bin/bash
set -e

echo "Removing all old Docker images without tags."

docker images --filter "dangling=true" -q | xargs -r docker rmi -f

echo "Successfully removed old Docker images."
