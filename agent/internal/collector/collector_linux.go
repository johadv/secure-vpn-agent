//go:build linux

package collector

import (
	"os"
	"path/filepath"
)

// countProcesses laskee /proc-hakemiston alla olevat numeeriset
// alihakemistot. Tämä tiedosto käännetään mukaan VAIN kun build-kohde
// on Linux (eli kun ajat tämän node-a:n/node-b:n sisällä).
func countProcesses() (int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		matched, err := filepath.Match("[0-9]*", entry.Name())
		if err == nil && matched {
			count++
		}
	}

	return count, nil
}