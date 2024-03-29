#!/bin/bash
# Inspired from: https://github.com/grpc/grpc-java/tree/master/examples#generating-self-signed-certificates-for-use-with-grpc

# Output files
# ca.key: Certificate Authority private key file (this shouldn't be shared in real-life)
# ca.crt: Certificate Authority trust certificate (this should be shared with users in real-life)
# server.key: Server private key, password protected (this shouldn't be shared)
# server.csr: Server certificate signing request (this should be shared with the CA owner)
# server.crt: Server certificate signed by the CA (this would be sent back by the CA owner) - keep on server
# server.pem: Conversion of server.key into a format gRPC likes (this shouldn't be shared)

# Summary 
# Private files: ca.key, server.key, server.pem, server.crt
# "Share" files: ca.crt (needed by the client), server.csr (needed by the CA)

# Changes these CN's to match your hosts in your environment if needed.
SERVER_CN=localhost
DIRECTORY=.ssl

if [ ! -d "$DIRECTORY" ]; then
  mkdir "$DIRECTORY"
  echo "creating directory $DIRECTORY"
else
  find $DIRECTORY -name postms_* -delete
  echo "cleaning directory $DIRECTORY/"
fi

# Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ".ssl/ca.key" 4096
openssl req -passin pass:1111 -new -x509 -days 3650 -key ".ssl/ca.key" -out ".ssl/ca.crt" -subj "/CN=${SERVER_CN}" -config ssl.cnf

# Step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out ".ssl/server.key" 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key ".ssl/server.key" -out ".ssl/server.csr" -subj "/CN=${SERVER_CN}"

# Step 4: Sign the certificate with the CA we created (it's called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 3650 -in ".ssl/server.csr" -CA ".ssl/ca.crt" -CAkey ".ssl/ca.key" -set_serial 01 -out ".ssl/server.crt" -extensions req_ext -extfile ssl.cnf

# Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in ".ssl/server.key" -out ".ssl/server.pem"

echo "done!"