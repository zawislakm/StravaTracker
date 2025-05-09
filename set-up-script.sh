#!/bin/bash
set -e

# Pull changes from repository
git fetch
git checkout zero-downtime-deploy
git pull origin zero-downtime-deploy

#!/bin/bash

set -e

NGINX_CONF="./nginx/nginx.conf"
NGINX_CONTAINER="nginx"

echo "👉 Step 1: Build and run app-green"
docker compose up -d --build app-green

echo "⏳ Waiting for app-green to start..."
sleep 5  # Można tu dodać curl http://app-green:8080/health

echo "🔁 Step 2: Switching NGINX to app-green"
sed -i 's/server app-blue:8080;/server app-green:8080;/' "$NGINX_CONF"
docker compose restart $NGINX_CONTAINER

echo "🔨 Step 3: Build and run new app-blue"
docker compose up -d --build app-blue

echo "⏳ Waiting for app-blue to start..."
sleep 5  # Można dodać curl http://app-blue:8080/health

echo "🔁 Step 4: Switching NGINX back to app-blue"
sed -i 's/server app-green:8080;/server app-blue:8080;/' "$NGINX_CONF"
docker compose restart $NGINX_CONTAINER

echo "🧹 Step 5: Stopping and removing app-green"
docker stop app-green
docker rm app-green

echo "✅ Blue-Green deploy complete, app-blue is live"

