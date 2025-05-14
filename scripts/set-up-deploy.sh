#!/bin/bash
set -e

git checkout master
git pull origin master

echo "Deployment started."

echo "Step 1: Stopping other running docker containers."
docker compose down

echo "Step 2: Starting new docker containers."
docker compose up -d --build nginx app-blue

echo "Deployment finished."