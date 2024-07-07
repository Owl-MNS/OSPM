#!/bin/bash

###############
#   Comment   #
###############
# This script can be used to deploy a cluster CockroachDB node in the 
# development environment using Docker

#################
#   Dependencies #
#################
# docker
# docker-compose

#################
#   Variables   #
#################
DOCKER_IMAGE="docker.arvancloud.ir/cockroachdb/cockroach:latest"

# The credentials to connect to the DB
CDB_DB_NAME="ospm"
CDB_LISTEN_PORT="26257"
CDB_ADMIN_UI_PORT="8080"
CDB_CONTAINER_NAME="cockroach-db"


#############
# Functions #
#############
function log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*"
}

function generate_docker_compose_file(){
  log "Generating the docker compose to start the DB..."
    cat <<EOL > docker-compose.yml
version: '3.7'

services:
  ${CDB_CONTAINER_NAME}:
    image: ${DOCKER_IMAGE}
    container_name: ${CDB_CONTAINER_NAME}
    hostname: ${CDB_CONTAINER_NAME}
    environment:
      - COCKROACH_DATABASE=${CDB_DB_NAME}
    ports:
      - "${CDB_LISTEN_PORT}:26257"
      - "${CDB_ADMIN_UI_PORT}:8080"
    volumes:
      - ${CDB_CONTAINER_NAME}:/cockroach/cockroach-data
    command: >
      start-single-node
      --insecure 
      --http-addr=${CDB_CONTAINER_NAME}:8080
    restart: unless-stopped	

volumes:
  ${CDB_CONTAINER_NAME}:

EOL

  log "Docker compose file generated at ${pwd}/docker-compose.yml"

}

function run_docker_compose(){
  log "starting the db..."
    docker-compose up -d
}

function create_db(){
  # Wait for CockroachDB to start
  until docker exec -it ${CDB_CONTAINER_NAME} /cockroach/cockroach sql -u root --insecure --execute="SHOW DATABASES;"; do
    echo "Waiting for CockroachDB to start..."
    sleep 2
  done

  docker exec -it ${CDB_CONTAINER_NAME} /cockroach/cockroach sql --insecure -u root --execute="CREATE DATABASE IF NOT EXISTS ${CDB_DB_NAME};"
  echo "DB ${CDB_DB_NAME} created"

}

function print_the_result(){
  clear
  echo "Cockroach DB for development environment has been deployed successfully"
  echo "Access Credentials are:"
  echo -e "\t Username: root"
  echo -e "\t Password: [there is not password required in dev environment. leave it blank!] "
  echo -e "\t DB Name: ${CDB_DB_NAME}"
  echo -e "\t DB listen address: 127.0.0.1:${CDB_LISTEN_PORT}"
  echo -e "\t Admin Web UI: http://127.0.0.1:${CDB_ADMIN_UI_PORT}"
}

#########
# Main  #
#########
generate_docker_compose_file
run_docker_compose
create_db
print_the_result