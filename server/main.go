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

// main käynnistää HTTP-palvelimen, joka vaatii mTLS:n eli sekä
// palvelimen ETTÄ clientin täytyy todistaa identiteettinsä
// sertifikaatilla, jotta yhteys hyväksytään.
func main() {
	// Luetaan CA:n julkinen sertifikaatti levyltä. Tätä käytetään
	// SEN TARKISTAMISEEN, että client esittää sertifikaatin jonka
	// JUURI tämä CA on allekirjoittanut - ei CA:n yksityistä avainta,
	// joten tämä tiedosto (ca.crt) on turvallinen lukea tässä.
	caCert, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatalf("CA-sertifikaatin luku epäonnistui: %v", err)
	}

	// CertPool on Go:n tapa pitää listaa "luotetuista" CA-sertifikaateista.
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	// tls.Config määrittää TÄSMÄLLEEN miten TLS-yhteys käsitellään.
	tlsConfig := &tls.Config{
		// ClientCAs kertoo: "näitä CA:ita luotetaan kun tarkistetaan
		// clientin esittämä sertifikaatti".
		ClientCAs: caPool,

		// RequireAndVerifyClientCert on TÄRKEIN rivi koko tiedostossa -
		// tämä PAKOTTAA mTLS:n. Ilman tätä Go tekisi tavallisen
		// yksisuuntaisen TLS:n (vain palvelin todistaa identiteettinsä).
		// Tämän kanssa: jos client ei esitä kelvollista sertifikaattia,
		// yhteys hylätään AUTOMAATTISESTI ennen kuin sovelluskoodi
		// handler-funktio ehtii ajaa ollenkaan.
		ClientAuth: tls.RequireAndVerifyClientCert,

		// MinVersion estää vanhojen, haavoittuvien TLS-versioiden
		// (1.0, 1.1, 1.2) käytön - vaaditaan tuorein standardi.
		MinVersion: tls.VersionTLS13,
	}

	mux := http.NewServeMux()

	// /ingest on ainoa endpoint - agentti lähettää tähän JSON-snapshotin.
	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		// r.TLS.PeerCertificates sisältää clientin esittämän sertifikaatin -
		// Go on JO tässä vaiheessa varmistanut sen kelvollisuuden
		// (RequireAndVerifyClientCert-asetuksen ansiosta), niin tässä
		// voidaan turvallisesti olettaa että puhuja on aito.
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
		// HUOM: tässä paikallisessa testausvaiheessa kuunnellaan
		// localhost:8443 - kun viedään tämä oikeasti node-a/node-b
		// väliseksi, tämä osoite vaihtuu WireGuard-tunnelin sisäiseen
		// IP-osoitteeseen 10.10.0.1:8443 EI julkiseen verkkoon.
		Addr:      "127.0.0.1:8443",
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	log.Println("Palvelin kuuntelee mTLS:llä osoitteessa 127.0.0.1:8443")
	// ListenAndServeTLS lataa palvelimen OMAN sertifikaatin (server.crt)
	// ja yksityisen avaimen (server.key) - näitä käytetään todistamaan
	// CLIENTILLE että palvelin on aito, samalla kun tlsConfig yllä
	// vaatii clientilta vastaavan todistuksen toiseen suuntaan.
	log.Fatal(server.ListenAndServeTLS("../certs/server.crt", "../certs/server.key"))
}