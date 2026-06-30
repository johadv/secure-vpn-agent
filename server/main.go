package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	caCert, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatalf("CA-sertifikaatin luku epäonnistui: %v", err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS13,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) > 0 {
			clientCN := r.TLS.PeerCertificates[0].Subject.CommonName
			log.Printf("Yhteys vahvistetulta clientilta: %s", clientCN)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "body-luku epäonnistui", http.StatusBadRequest)
			return
		}

		fmt.Printf("Saatu snapshot: %s\n", body)
		w.WriteHeader(http.StatusAccepted)
	})

	server := &http.Server{
		// PÄIVITETTY: kuunnellaan nyt WireGuard-tunnelin sisäisessä
		// osoitteessa (10.10.0.1) sen sijaan että kuunneltaisiin
		// julkisessa LAN-verkossa (192.168.64.x) tai vain paikallisesti
		// (127.0.0.1). Tämä tarkoittaa: palvelu EI ole tavoitettavissa
		// MITENKÄÄN ilman WireGuard-yhteyttä - sama "defense in depth"
		// -periaate jota suunniteltiin jo projektin alussa.
		Addr:      "10.10.0.1:8443",
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	log.Println("Palvelin kuuntelee mTLS:llä WireGuard-osoitteessa 10.10.0.1:8443")
	log.Fatal(server.ListenAndServeTLS("../certs/server.crt", "../certs/server.key"))
}