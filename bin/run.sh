#!/usr/bin/env bash
set -euo pipefail

echo "[run.sh] Starting service"

if [ -z "${DATABASE_URL:-}" ]; then
  echo "[run.sh] DATABASE_URL is empty, skipping migrations"
else
  echo "[run.sh] Running DB migrations"
  goose -dir ./db/migrations postgres "${DATABASE_URL}" up
fi

echo "[run.sh] Starting Go app"
exec /app/bin/app