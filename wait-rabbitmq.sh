#!/bin/sh
# wait-rabbitmq.sh

# Wait for RabbitMQ service to be up and running.

set -e

host="$1"
cmd="$2"

until curl -I "http://guest:guest@$host/api/overview"; do
  >&2 echo "RabbitMQ is unavailable - sleeping..."
  sleep 3
done

>&2 echo "RabbitMQ is up - executing command $cmd"
exec $cmd
