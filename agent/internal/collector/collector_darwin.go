//go:build darwin

package collector

import "os/exec"

// countProcesses macOS:lla ei voi lukea /proc-hakemistoa (sitä ei ole),
// joten käytetään `ps`-komentoa ja lasketaan rivien määrä. Tämä on vain
// kehitys-/testauskäyttöön Macilla - tuotantoagentti pyörii lopulta
// Linux-VM:issä, joissa collector_linux.go:n /proc-toteutus on tarkempi.
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