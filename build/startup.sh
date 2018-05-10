#!/usr/bin/env bash

export TUNNEL_ORIGIN_CERT="$(pwd)/argo_tunnel.pem"
openssl enc -d -aes-256-cbc -in argo_tunnel.pem.enc -out argo_tunnel.pem -k $CERT_DECRYPT_KEY
chmod +x ./cloudflared
echo $ARGO_DOMAIN
touch app.log tunnel.log
nohup ./cloudflared --hostname $ARGO_DOMAIN http://localhost:9080 >> tunnel.log 2>&1 & echo $!
nohup ./input_tracker_app >> app.log 2>&1 & echo $!
tail -qF *.log