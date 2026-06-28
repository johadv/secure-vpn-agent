#!/usr/bin/env bash
set -euo pipefail

# Tämä skripti luo oman, pienen "varmentajan" (CA = Certificate Authority)
# vain tätä projektia varten. Oikeassa tuotantokäytössä sertifikaatit
# tulisivat oikealta CA:lta (esim. Let's Encrypt), mutta mTLS:n
# OPETTELUUN oma CA on täysin riittävä ja yleinen tapa.
mkdir -p certs && cd certs

echo "1/3: Luodaan juuri-CA..."
# genrsa luo 4096-bittisen RSA-avainparin CA:lle. Tämä on CA:n OMA
# avain, jota käytetään allekirjoittamaan kaikki muut sertifikaatit -
# jos tämä avain vuotaa, kuka tahansa voisi väärentää sertifikaatteja
# joita "palvelin" ja "agentti" luottavat, niin .gitignore estää
# tämän koko certs/-kansion päätymisen Git-historiaan.
openssl genrsa -out ca.key 4096

# req -x509 luo itseallekirjoitetun juurisertifikaatin (CA:n oman
# "henkilökortin"). -days 3650 = kelpoisuusaika noin 10 vuotta,
# riittää reilusti tämän oppimisprojektin elinkaarelle.
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 \
  -out ca.crt -subj "/CN=secure-vpn-agent-CA"

echo "2/3: Luodaan palvelimen (node-a/node-b) sertifikaatti..."
openssl genrsa -out server.key 2048
# CSR (Certificate Signing Request) on "pyyntö" saada sertifikaatti -
# CN (Common Name) kertoo kenelle sertifikaatti on tarkoitettu.
openssl req -new -key server.key -out server.csr -subj "/CN=node-server"
# CA allekirjoittaa pyynnön omalla avaimellaan -> syntyy server.crt,
# jonka kuka tahansa CA:n julkisen sertifikaatin (ca.crt) tunteva
# osapuoli voi vahvistaa aidoksi.
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out server.crt -days 825 -sha256

echo "3/3: Luodaan agentin (clientin) sertifikaatti..."
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=node-agent"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out client.crt -days 825 -sha256

echo ""
echo "Valmis. Luodut tiedostot certs/-kansiossa:"
ls -la
