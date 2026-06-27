package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/johadv/secure-vpn-agent/agent/internal/collector"
)

// main on ohjelman aloituskohta - Go käynnistää aina tämän funktion
// ensimmäisenä kun ajat "go run main.go" tai käynnistät käännetyn binäärin.
func main() {
	// Kutsuu collector.Collect()-funktiota  (collector.go tiedostossa).
	// Tämä palauttaa Snapshot-structin täytettynä nykyisellä prosessimäärällä,
	// sekä mahdollisen virheen (err) jos tiedon kerääminen epäonnistui.
	snap, err := collector.Collect()
	if err != nil {
		// log.Fatalf tulostaa virheviestin JA lopettaa ohjelman heti
		// (vastaa os.Exit(1):tä virheviestin kanssa). Käytetään tätä
		// vain main-funktiossa - syvemmällä koodissa virheet palautetaan
		// kutsujalle eikä koskaan kaadeta ohjelmaa suoraan.
		log.Fatalf("collect epäonnistui: %v", err)
	}

	// json.MarshalIndent muuttaa Snapshot-structin JSON-tekstiksi.
	// "" ja "  " -parametrit määräävät sisennyksen (kaksi välilyöntiä per
	// taso), jotta tuloste on ihmisen luettavaa eikä yhtä pitkää riviä.
	out, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		log.Fatalf("JSON-muunnos epäonnistui: %v", err)
	}

	// Tulostaa valmiin JSON-tekstin konsoliin. "string(out)" muuttaa
	// MarshalIndentin palauttaman tavusarjan ([]byte) tekstiksi (string),
	// jotta Println näyttää sen luettavana tekstinä eikä raakoina tavuina.
	fmt.Println(string(out))
}