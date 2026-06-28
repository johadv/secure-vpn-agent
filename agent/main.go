package main

import (
	"encoding/json"
	"log"

	"github.com/johadv/secure-vpn-agent/agent/internal/collector"
	"github.com/johadv/secure-vpn-agent/agent/internal/transport"
)

// main on ohjelman aloituskohta.
func main() {
	// Kutsuu collector.Collect()-funktiota - tämä palauttaa Snapshot-
	// structin täytettynä koneen nykyisellä tilalla (prosessit, levy,
	// avoimet portit).
	snap, err := collector.Collect()
	if err != nil {
		log.Fatalf("collect epäonnistui: %v", err)
	}

	// json.Marshal muuttaa Snapshot-structin JSON-tavusarjaksi
	// lähetystä varten - ei tarvita sisennystä (MarshalIndent) tässä,
	// koska tätä ei enää tulosteta ihmiselle luettavaksi vaan
	// lähetetään suoraan verkon yli.
	payload, err := json.Marshal(snap)
	if err != nil {
		log.Fatalf("JSON-muunnos epäonnistui: %v", err)
	}

	// NewClient rakentaa mTLS-yhteyden osaavan HTTP-clientin käyttäen
	// agentin omaa sertifikaattia/avainta ja CA:n sertifikaattia.
	// Polut ovat suhteessa agent/-kansioon, sama tapa kuin server-puolella.
	client, err := transport.NewClient(
		"../certs/ca.crt",
		"../certs/client.crt",
		"../certs/client.key",
	)
	if err != nil {
		log.Fatalf("mTLS-clientin luonti epäonnistui: %v", err)
	}

	// Send lähettää JSON-payloadin palvelimelle POST-pyyntönä.
	// HUOM: tässä paikallisessa testausvaiheessa osoite on
	// 127.0.0.1:8443 - kun viedään tämä node-a/node-b-väliseksi,
	// tämä vaihtuu WireGuard-tunnelin sisäiseen IP-osoitteeseen.
	err = transport.Send(client, "https://127.0.0.1:8443/ingest", payload)
	if err != nil {
		log.Fatalf("lähetys epäonnistui: %v", err)
	}

	log.Println("Snapshot lähetetty onnistuneesti palvelimelle.")
}