#!/usr/bin/env bash
set -euo pipefail

mkdir -p certs && cd certs

echo "1/3: Luodaan juuri-CA..."
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 \
  -out ca.crt -subj "/CN=secure-vpn-agent-CA"

echo "2/3: Luodaan palvelimen sertifikaatti..."
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=node-server"
# SAN sisältää NYT kaikki osoitteet joissa palvelinta saatetaan ajaa:
# 127.0.0.1 (paikallinen testaus Macilla), 10.10.0.1 (WireGuard-tunnelin
# sisäinen osoite node-b:llä), ja localhost/node-b nimillä. Tämä tekee
# samasta sertifikaatista käyttökelpoisen sekä paikalliseen testaukseen
# että oikeaan WireGuard-yli-tapahtuvaan liikenteeseen, ilman että
# sertifikaatteja tarvitsee generoida erikseen joka tilanteeseen.
cat > server-ext.cnf << 'CNFEOF'
subjectAltName = IP:127.0.0.1,IP:10.10.0.1,DNS:localhost,DNS:node-server,DNS:node-b
CNFEOF
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out server.crt -days 825 -sha256 \
  -extfile server-ext.cnf

echo "3/3: Luodaan agentin (clientin) sertifikaatti..."
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=node-agent"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out client.crt -days 825 -sha256

echo ""
echo "Valmis. Luodut tiedostot certs/-kansiossa:"
ls -la