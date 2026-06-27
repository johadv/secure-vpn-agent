//go:build linux

package collector

import (
	"os"
	"path/filepath"
	"syscall"
)

// countProcesses laskee /proc-hakemiston alla olevat numeeriset
// alihakemistot - Linuxissa jokainen käynnissä oleva prosessi näkyy
// täällä omana kansionaan, nimettynä prosessin PID:llä (esim. /proc/1234).
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

// diskUsedPct laskee juuriosion (/) käytetyn tilan prosentteina.
// syscall.Statfs on matalan tason käyttöjärjestelmäkutsu joka kysyy
// tiedostojärjestelmältä suoraan lohkojen (blocks) kokonaismäärän ja
// vapaiden lohkojen määrän - samaa tietoa jota "df"-komento käyttää.
func diskUsedPct() (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return 0, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free

	if total == 0 {
		return 0, nil
	}

	return (float64(used) / float64(total)) * 100, nil
}