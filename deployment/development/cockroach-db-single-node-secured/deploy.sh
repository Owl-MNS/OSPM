#!/bin/bash

###############
#   Comment   #
###############
# This script can be used to deploy a single CockroachDB node in the 
# development environment using Docker

#################
#   Dependencies #
#################
# docker
# docker-compose
# openssl

#################
#   Variables   #
#################
# The IP address of the host where the Docker container 
# of the CockroachDB is installed
HOST_IP="127.0.0.1"

CDB_LISTEN_ADDRESS="127.0.0.1"
CDB_LISTEN_PORT="26257"
CDB_ADMIN_UI_PORT="8080"

# The credentials to connect to the DB
CDB_USERNAME="user1" #[A-Z][a-z][0-9] do not use special chars
CDB_PASSWORD="foobar"
CDB_DB_NAME="ospm"

# The directory paths where the cert files will be stored and mounted
# to the CockroachDB container
CERTS_DIR="certs"
CA_KEY_DIR="certs-safe"
CA_KEY_PATH="${CA_KEY_DIR}/ca.key"

CDB_DOCKER_IMAGE="docker.arvancloud.ir/cockroachdb/cockroach:latest"
CDB_CONTAINER_NAME="cockroach-db"

function log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*"
}


# Know more:
# https://www.cockroachlabs.com/docs/v21.2/create-security-certificates-openssl#step-1-create-the-ca-key-and-certificate-pair
function generate_ssl_cert_ca_config_file(){
    cat <<EOL > ca.cnf

# OpenSSL CA configuration file
[ ca ]
default_ca = CA_default

[ CA_default ]
default_days = 365
database = index.txt
serial = serial.txt
default_md = sha256
copy_extensions = copy
unique_subject = no

# Used to create the CA certificate.
[ req ]
prompt=no
distinguished_name = distinguished_name
x509_extensions = extensions

[ distinguished_name ]
organizationName = Cockroach
commonName = Cockroach CA

[ extensions ]
keyUsage = critical,digitalSignature,nonRepudiation,keyEncipherment,keyCertSign
basicConstraints = critical,CA:true,pathlen:1

# Common policy for nodes and users.
[ signing_policy ]
organizationName = supplied
commonName = optional

# Used to sign node certificates.
[ signing_node_req ]
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth,clientAuth

# Used to sign client certificates.
[ signing_client_req ]
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = clientAuth


EOL

}




# Know more:
# https://www.cockroachlabs.com/docs/v21.2/create-security-certificates-openssl#step-1-create-the-ca-key-and-certificate-pair
function generate_ssl_certs_ca(){
    mkdir -p "${CERTS_DIR}" "${CA_KEY_DIR}"

    # Generate CA certificate and key if they don't exist
    if [ ! -f "$CA_KEY_PATH" ]; then
        log "Generating CA key..."
        openssl genrsa -out "$CA_KEY_PATH" 4096
    else
        log "CA key already exists."
    fi

    if [ ! -f "${CERTS_DIR}/ca.crt" ]; then
        log "Generating CA certificate..."
        openssl req \
        -new \
        -x509 \
        -config ca.cnf \
        -key "$CA_KEY_PATH" \
        -out "${CERTS_DIR}/ca.crt" \
        -days 365 -batch
    else
        log "CA certificate already exists."
    fi

    rm -f index.txt serial.txt
    touch index.txt
    echo '01' > serial.txt
}

# Know more:
# https://www.cockroachlabs.com/docs/v21.2/create-security-certificates-openssl#step-2-create-the-certificate-and-key-pairs-for-nodes
function generate_ssl_cert_node_config_file(){
cat <<EOL > node.cnf
# OpenSSL node configuration file
[ req ]
prompt=no
distinguished_name = distinguished_name
req_extensions = extensions

[ distinguished_name ]
organizationName = Cockroach

[ extensions ]
subjectAltName = critical,DNS:${CDB_CONTAINER_NAME},DNS:local,IP:${HOST_IP}

EOL


}

# Know more:
# https://www.cockroachlabs.com/docs/v21.2/create-security-certificates-openssl#step-2-create-the-certificate-and-key-pairs-for-nodes
function generate_ssl_certs_node(){
    if [ ! -f ${CERTS_DIR}/node.key ]; then
        echo "Generating node key..."
        openssl genrsa -out ${CERTS_DIR}/node.key 4096
    else
        echo "Node key already exists."
    fi

    if [ ! -f ${CERTS_DIR}/node.crt ]; then
        echo "Generating node certificate..."
        openssl req \
        -new -config node.cnf \
        -key ${CERTS_DIR}/node.key \
        -out ${CERTS_DIR}/node.csr \
        -batch

        openssl ca \
            -config ca.cnf \
            -keyfile ${CA_KEY_PATH}/ca.key \
            -cert ${CERTS_DIR}/ca.crt \
            -policy signing_policy \
            -extensions signing_node_req \
            -out ${CERTS_DIR}/node.crt \
            -outdir ${CERTS_DIR} \
            -in node.csr -batch

        openssl x509 -in ${CERTS_DIR}/node.crt -text | grep "X509v3 Subject Alternative Name" -A 1

    else
        echo "Node certificate already exists."
    fi
}

# Know more:
# https://www.cockroachlabs.com/docs/v21.2/create-security-certificates-openssl#step-3-create-the-certificate-and-key-pair-for-the-first-user
function generate_ssl_certs_user_root_config(){
    cat <<EOL > root.cnf
    [ req ]
prompt=no
distinguished_name = distinguished_name
req_extensions = extensions

[ distinguished_name ]
organizationName = Cockroach
commonName = root

[ extensions ]
subjectAltName = DNS:root
EOL
}

function generate_ssl_certs_user_root(){
    if [ ! -f "${CERTS_DIR}/client.root.key" ]; then
        log "Generating root client key..."
        openssl genrsa -out "${CERTS_DIR}/client.root.key" 4096
    else
        log "Root client key already exists."
    fi

    if [ ! -f "${CERTS_DIR}/client.root.crt" ]; then
        log "Generating root client certificate..."
        openssl req\
        -new \
        -config root.cnf \
        -key "${CERTS_DIR}/client.root.key" \
        -out "${CERTS_DIR}/client.root.csr" \
        -batch 


        openssl ca \
        -config ca.cnf \
        -keyfile ${CA_KEY_PATH}/ca.key \
        -cert ${CERTS_DIR}/ca.crt \
        -policy signing_policy \
        -extensions signing_client_req \
        -out certs/client.root.crt \
        -outdir ${CERTS_DIR}/ \
        -in client.root.csr \
        -batch

        openssl x509 -in certs/client.root.crt -text | grep CN=

    else
        log "Root client certificate already exists."
    fi
}


function generate_ssl_certs_user_config(){
        cat <<EOL > ${CDB_USERNAME}.cnf
    [ req ]
prompt=no
distinguished_name = distinguished_name
req_extensions = extensions

[ distinguished_name ]
organizationName = Cockroach
commonName = ${CDB_USERNAME}

[ extensions ]
subjectAltName = DNS:${CDB_USERNAME}
EOL

}

function generate_ssl_certs_user(){
    if [ ! -f "${CERTS_DIR}/client.${CDB_USERNAME}.key" ]; then
        log "Generating ${CDB_USERNAME} client key..."
        openssl genrsa -out "${CERTS_DIR}/client.${CDB_USERNAME}.key" 4096
    else
        log "${CDB_USERNAME} client key already exists."
    fi

    if [ ! -f "${CERTS_DIR}/client.root.crt" ]; then
        log "Generating ${CDB_USERNAME} client certificate..."
        openssl req\
        -new \
        -config root.cnf \
        -key "${CERTS_DIR}/client.${CDB_USERNAME}.key" \
        -out "${CERTS_DIR}/client.${CDB_USERNAME}.csr" \
        -batch 


        openssl ca \
        -config ca.cnf \
        -keyfile ${CA_KEY_PATH}/ca.key \
        -cert ${CERTS_DIR}/ca.crt \
        -policy signing_policy \
        -extensions signing_client_req \
        -out certs/client.${CDB_USERNAME}.crt \
        -outdir ${CERTS_DIR}/ \
        -in client.${CDB_USERNAME}.csr \
        -batch

        openssl x509 -in certs/client.${CDB_USERNAME}.crt -text | grep CN=

    else
        log "${CDB_USERNAME} client certificate already exists."
    fi
}

function generate_docker_compose_file(){
    cat <<EOL > docker-compose.yml
version: '3.7'

services:
  ${CDB_CONTAINER_NAME}:
    image: ${CDB_DOCKER_IMAGE}
    container_name: ${CDB_CONTAINER_NAME}
    command: /cockroach/entrypoint.sh
    ports:
      - "${CDB_LISTEN_PORT}:${CDB_LISTEN_PORT}"
      - "${CDB_ADMIN_UI_PORT}:8080"
    volumes:
      - ./${CERTS_DIR}/ca.crt:/cockroach/certs/ca.crt
      - ./${CERTS_DIR}/node.crt:/cockroach/certs/node.crt
      - ./${CERTS_DIR}/node.key:/cockroach/certs/node.key
      - ./${CERTS_DIR}/client.root.key:/cockroach/certs/client.root.key
      - ./${CERTS_DIR}/client.${CDB_USERNAME}.key:/cockroach/certs/client.${CDB_USERNAME}.key
      - ./${CERTS_DIR}/client.root.crt:/cockroach/certs/client.root.crt
      - ./${CERTS_DIR}/client.${CDB_USERNAME}.crt:/cockroach/certs/client.${CDB_USERNAME}.crt
      - ./init-cockroach.sh:/docker-entrypoint-initdb.d/init-cockroach.sh
      - ./entrypoint.sh:/cockroach/entrypoint.sh
    environment:
      - COCKROACH_USER=${CDB_USERNAME}
      - COCKROACH_PASSWORD=${CDB_PASSWORD}
      - COCKROACH_DATABASE=${CDB_DB_NAME}
      - HOST_IP=${HOST_IP}
    restart: unless-stopped
    entrypoint: ["/cockroach/entrypoint.sh"]

volumes:
  cockroach-data:
EOL
}

function generate_cockroachdb_init_script(){
    cat <<EOL > init-cockroach.sh
#!/bin/bash

# Create user, database, and grant privileges
cockroach sql --certs-dir=/cockroach/certs --execute="
CREATE USER IF NOT EXISTS ${CDB_USERNAME} WITH PASSWORD '${CDB_PASSWORD}';
CREATE DATABASE IF NOT EXISTS ${CDB_DB_NAME};
GRANT ALL ON DATABASE ${CDB_DB_NAME} TO ${CDB_USERNAME};
"
EOL

chmod +x init-cockroach.sh
}

function generate_cockroach_db_entrypoint_script(){
    cat <<EOL > entrypoint.sh
#!/bin/bash

# Set permissions on the certificate files
chmod 640 /cockroach/certs/node.key
chmod 640 /cockroach/certs/client.root.key
chmod 640 /cockroach/certs/client.${CDB_USERNAME}.key
chmod 644 /cockroach/certs/ca.crt
chmod 644 /cockroach/certs/node.crt
chmod 644 /cockroach/certs/client.root.crt
chmod 644 /cockroach/certs/client.${CDB_USERNAME}.crt

# Start CockroachDB
cockroach start-single-node --certs-dir=/cockroach/certs --advertise-addr=${HOST_IP} --listen-addr=${CDB_LISTEN_ADDRESS}:${CDB_LISTEN_PORT} --background

# Wait for CockroachDB to start
until cockroach sql --certs-dir=/cockroach/certs --execute="SHOW DATABASES;" &>/dev/null; do
  echo "Waiting for CockroachDB to start..."
  sleep 2
done

# Run the initialization script
/docker-entrypoint-initdb.d/init-cockroach.sh


EOL

chmod +x entrypoint.sh
}

function deploy_cockroach_db(){
    docker-compose up -d --remove-orphans
}

###########
#   Main  #
###########
log "Generating SSL certificates..."
generate_ssl_cert_ca_config_file
generate_ssl_certs_ca
generate_ssl_certs_node
generate_ssl_certs_user_root
generate_ssl_certs_user

log "Generating Docker Compose file..."
generate_docker_compose_file

log "Generating CockroachDB init script..."
generate_cockroachdb_init_script

log "Generating CockroachDB entrypoint script..."
generate_cockroach_db_entrypoint_script

log "Deploying CockroachDB..."
deploy_cockroach_db

log "Deployment complete."
