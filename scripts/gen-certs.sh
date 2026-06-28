#!/usr/bin/env bash
set -euo pipefail

# Tämä skripti luo oman, pienen "varmentajan" (CA = Certificate Authority)
# vain tätä projektia varten.
mkdir -p certs && cd certs

echo "1/3: Luodaan juuri-CA..."
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 \
  -out ca.crt -subj "/CN=secure-vpn-agent-CA"

echo "2/3: Luodaan palvelimen sertifikaatti (SAN: 127.0.0.1, localhost)..."
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=node-server"
# -extfile + extensions-tiedosto lisää SAN-kentän (Subject Alternative
# Name) sertifikaattiin. Nykyaikaiset TLS-kirjastot (mukaan lukien Go)
# VAATIVAT tämän kentän osoitteen validointiin - vanha CN-pohjainen
# tarkistus on poistettu käytöstä tietoturvasyistä (CN olisi voitu
# helposti väärinkäyttää, koska sen merkitys oli epäselvä standardissa).
cat > server-ext.cnf << 'CNFEOF'
subjectAltName = IP:127.0.0.1,DNS:localhost,DNS:node-server
CNFEOF
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out server.crt -days 825 -sha256 \
  -extfile server-ext.cnf

echo "3/3: Luodaan agentin (clientin) sertifikaatti..."
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=node-agent"
# Clientin sertifikaatille ei tarvita SAN:ia samalla tavalla, koska
# palvelin ei tarkista MISTÄ osoitteesta client yhdistää - se tarkistaa
# vain ETTÄ client-sertifikaatti on CA:n allekirjoittama.
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out client.crt -days 825 -sha256

echo ""
echo "Valmis. Luodut tiedostot certs/-kansiossa:"
ls -la