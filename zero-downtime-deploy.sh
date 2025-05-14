#!/bin/bash
set -e

#git checkout master
#git pull origin master

CONFIG_DIR="./nginx/conf.d"
DEFAULT_CONF="$CONFIG_DIR/default.conf"
NGINX_CONTAINER="nginx-proxy"
APP_BLUE="app-blue"
APP_GREEN="app-green"

wait_for_healthy() {
  local container="$1"
  echo "Waiting for $container to become healthy..."

  for i in {1..15}; do
    status=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "unknown")
    echo "Try $i/15: status = $status"

    if [[ "$status" == "healthy" ]]; then
      echo "$container is healthy"
      return 0
    fi

    sleep 3
  done

  echo "$container did not become healthy in time. Aborting."
  exit 1
}

echo "Zero downtime deployment started."

echo "Step 1: New version started up on container: app-green."
docker compose up -d --build $APP_GREEN
wait_for_healthy $APP_GREEN

echo "Step 2: NGINX proxy switch to container: app-green."
ln -nsf nginx-green.conf "$DEFAULT_CONF"
docker exec $NGINX_CONTAINER nginx -s reload

echo "Step 3: New version started up on container: app-blue."
docker compose up -d --build $APP_BLUE
wait_for_healthy $APP_BLUE

echo "Step 4: NGINX proxy switch to container: app-blue."
ln -nsf nginx-blue.conf "$DEFAULT_CONF"
docker exec $NGINX_CONTAINER nginx -s reload

echo "Step 5: Stop and remove container: app-green."
docker stop $APP_GREEN || true
docker rm $APP_GREEN || true

echo "Zero downtime deployment finished successfully. App running on container: app-blue."

