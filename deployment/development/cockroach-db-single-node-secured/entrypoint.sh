#!/bin/bash

# Set permissions on the certificate files
chmod 640 /cockroach/certs/node.key
chmod 640 /cockroach/certs/client.root.key
chmod 640 /cockroach/certs/client.user1.key
chmod 644 /cockroach/certs/ca.crt
chmod 644 /cockroach/certs/node.crt
chmod 644 /cockroach/certs/client.root.crt
chmod 644 /cockroach/certs/client.user1.crt

# Start CockroachDB
cockroach start-single-node --certs-dir=/cockroach/certs --advertise-addr=127.0.0.1 --listen-addr=127.0.0.1:26257 --background

# Wait for CockroachDB to start
until cockroach sql --certs-dir=/cockroach/certs --execute="SHOW DATABASES;" &>/dev/null; do
  echo "Waiting for CockroachDB to start..."
  sleep 2
done

# Run the initialization script
/docker-entrypoint-initdb.d/init-cockroach.sh


