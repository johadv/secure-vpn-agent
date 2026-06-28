package transport

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// NewClient rakentaa HTTP-clientin joka osaa todistaa identiteettinsä
// mTLS:llä. caPath, certPath ja keyPath ovat polkuja levyllä oleviin
// PEM-tiedostoihin (CA:n sertifikaatti, oma sertifikaatti, oma avain).
func NewClient(caPath, certPath, keyPath string) (*http.Client, error) {
	// Luetaan CA:n sertifikaatti - tällä clientti tarkistaa, että
	// PALVELIMEN esittämä sertifikaatti on aidosti tämän CA:n
	// allekirjoittama (eikä esim. joku väliin asettunut hyökkääjä).
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("CA-sertifikaatin luku epäonnistui: %w", err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	// LoadX509KeyPair lataa CLIENTIN OMAN sertifikaatin ja yksityisen
	// avaimen yhdeksi pariksi - tätä paria käytetään PALVELIMELLE
	// todistamiseen "minä olen oikeasti node-agent, tässä todiste".
	clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("client-sertifikaatin lataus epäonnistui: %w", err)
	}

	tlsConfig := &tls.Config{
		// RootCAs: mihin CA:han LUOTETAAN kun tarkistetaan palvelimen
		// sertifikaatti (vastaa palvelimen ClientCAs-kenttää, mutta
		// toiseen suuntaan).
		RootCAs: caPool,

		// Certificates: oma sertifikaatti/avain-pari, jonka Go lähettää
		// AUTOMAATTISESTI palvelimelle TLS-handshaken aikana, kun
		// palvelin pyytää client-sertifikaattia (RequireAndVerifyClientCert).
		Certificates: []tls.Certificate{clientCert},

		MinVersion: tls.VersionTLS13,
	}

	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}, nil
}

// Send lähettää JSON-payloadin annettuun URL-osoitteeseen POST-pyyntönä.
func Send(client *http.Client, url string, payload []byte) error {
	resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("POST-pyyntö epäonnistui: %w", err)
	}
	// defer varmistaa että yhteys suljetaan AINA funktion lopussa,
	// vaikka tapahtuisi virhe alempana - estää resurssivuodon.
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("palvelin vastasi statuksella %d", resp.StatusCode)
	}

	return nil
}