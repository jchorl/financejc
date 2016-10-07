#! /bin/bash

if [ -n "$DEV" ]; then
        mkdir -p /etc/letsencrypt/live/$DOMAIN
        openssl genrsa -out /etc/letsencrypt/live/$DOMAIN/privkey.pem 2048 && openssl req -new -x509 -sha256 -key /etc/letsencrypt/live/$DOMAIN/privkey.pem -out /etc/letsencrypt/live/$DOMAIN/fullchain.pem -days 3650 -subj '/CN='$DOMAIN':'$PORT'/O=CollabTest/C=CA'
        financejc
else
        # TODO write a script that generates certs with letsencrypt and then starts the server
        go-wrapper run
fi