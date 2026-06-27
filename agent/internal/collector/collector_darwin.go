//go:build darwin

package collector

import (
	"os/exec"
	"syscall"
)

// countProcesses macOS:lla ei voi lukea /proc-hakemistoa (sitä ei ole),
// joten käytetään `ps`-komentoa ja lasketaan rivien määrä.
func countProcesses() (int, error) {
	out, err := exec.Command("ps", "-A", "-o", "pid=").Output()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, b := range out {
		if b == '\n' {
			count++
		}
	}

	return count, nil
}

// diskUsedPct laskee juuriosion (/) käytetyn tilan prosentteina.
// macOS:lla syscall.Statfs on käytettävissä samalla nimellä kuin
// Linuxissa, mutta sen sisäinen struct-rakenne (Statfs_t) on hieman
// erilainen Applen BSD-pohjaisen kernelin takia - siksi tämä on
// kirjoitettu omana, erillisenä tiedostona Linux-version sijaan.
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