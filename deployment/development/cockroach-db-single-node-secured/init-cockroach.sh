#!/bin/bash

# Create user, database, and grant privileges
cockroach sql --certs-dir=/cockroach/certs --execute="
CREATE USER IF NOT EXISTS user1 WITH PASSWORD 'foobar';
CREATE DATABASE IF NOT EXISTS ospm;
GRANT ALL ON DATABASE ospm TO user1;
"
