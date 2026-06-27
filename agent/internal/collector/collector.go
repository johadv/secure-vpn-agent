package collector

import "time"

// Snapshot on yksi "kuvakaappaus" koneen tilasta yhdellä ajanhetkellä.
type Snapshot struct {
	Timestamp    time.Time `json:"timestamp"`
	ProcessCount int       `json:"process_count"`
}

// Collect kerää nykyisen tilan koneelta ja palauttaa sen Snapshot-muodossa.
// Itse prosessilaskenta (countProcesses) on toteutettu erikseen joka
// käyttöjärjestelmälle - katso collector_linux.go ja collector_darwin.go.
// Go valitsee automaattisesti oikean tiedoston build-aikana sen
// perusteella millä käyttöjärjestelmällä `go build`/`go run` ajetaan.
func Collect() (Snapshot, error) {
	count, err := countProcesses()
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Timestamp:    time.Now().UTC(),
		ProcessCount: count,
	}, nil
}