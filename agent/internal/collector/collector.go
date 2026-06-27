package collector

import "time"

// Snapshot on yksi "kuvakaappaus" koneen tilasta yhdellä ajanhetkellä.
type Snapshot struct {
	Timestamp    time.Time `json:"timestamp"`
	ProcessCount int       `json:"process_count"`
	DiskUsedPct  float64   `json:"disk_used_pct"`
	OpenPorts    []int     `json:"open_ports"`
}

// Collect kerää nykyisen tilan koneelta ja palauttaa sen Snapshot-muodossa.
// countProcesses ja diskUsedPct ovat käyttöjärjestelmäkohtaisia
// (collector_linux.go / collector_darwin.go). countOpenPorts on yhteinen
// kaikille käyttöjärjestelmille, koska se käyttää vain Go:n net-pakettia.
func Collect() (Snapshot, error) {
	count, err := countProcesses()
	if err != nil {
		return Snapshot{}, err
	}

	diskPct, err := diskUsedPct()
	if err != nil {
		return Snapshot{}, err
	}

	ports := scanCommonPorts()

	return Snapshot{
		Timestamp:    time.Now().UTC(),
		ProcessCount: count,
		DiskUsedPct:  diskPct,
		OpenPorts:    ports,
	}, nil
}