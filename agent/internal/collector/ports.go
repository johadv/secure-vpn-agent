package collector

import (
	"net"
	"strconv"
	"time"
)

// commonPorts on lista yleisimmistä porteista joita tarkistetaan.
// Tämä EI ole täydellinen 1-65535 skannaus - tarkoituksella rajattu
// vain tunnettuihin, turvallisuusrelevantteihin palveluportteihin
// (SSH, HTTP/HTTPS, tietokannat, jne). Täydellinen skannaus olisi
// hitaampi ja antaisi vain vähän lisäarvoa tähän käyttötarkoitukseen.
var commonPorts = []int{
	22,   // SSH
	80,   // HTTP
	443,  // HTTPS
	3306, // MySQL
	5432, // PostgreSQL
	6379, // Redis
	8080, // yleinen kehitys-/proxy-portti
	8443, // tämän projektin oma mTLS-portti (Vaihe 3)
}

// scanCommonPorts tarkistaa jokaisen commonPorts-listan portin ja
// palauttaa listan niistä jotka ovat auki PAIKALLISELLA koneella.
//
// TÄRKEÄ TURVALLISUUSPERIAATE: agentti tarkistaa AINOASTAAN omaa
// konettaan (127.0.0.1 / localhost) - se ei koskaan yritä yhdistää
// mihinkään muuhun IP-osoitteeseen verkossa. Toisten koneiden porttien
// skannaaminen olisi portskannausta, joka on monissa verkoissa ja
// lainkäyttöalueilla luvattomana tehtynä epäeettistä tai laitonta.
// Tässä agentti raportoi vain "mitä MINUN koneellani on auki",
// ei tutki tai testaa mitään muuta konetta verkossa.
func scanCommonPorts() []int {
	open := []int{}

	for _, port := range commonPorts {
		address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

		// DialTimeout yrittää avata TCP-yhteyden. Lyhyt timeout (200ms)
		// varmistaa, että agentti ei jää jumiin odottamaan vastausta
		// jos portti ei vastaa mitään - paikallisessa koneessa vastaus
		// tulee normaalisti millisekunneissa.
		conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
		if err != nil {
			// Virhe tässä tarkoittaa yksinkertaisesti "portti ei ole auki",
			// ei mitään vakavampaa - ei tarvitse käsitellä erikseen.
			continue
		}

		conn.Close()
		open = append(open, port)
	}

	return open
}