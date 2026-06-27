package collector

import "testing"

// TestCollect varmistaa, että Collect() palauttaa järkeviä arvoja
// eikä kaadu millään käyttöjärjestelmällä jolla testit ajetaan.
// Go:n testitiedostot päättyvät aina "_test.go" - "go test" löytää
// ne automaattisesti eikä niitä käännetä mukaan tavalliseen buildiin.
func TestCollect(t *testing.T) {
	snap, err := Collect()
	if err != nil {
		t.Fatalf("Collect() palautti virheen: %v", err)
	}

	// ProcessCount pitäisi aina olla positiivinen - koneella on aina
	// ainakin muutama prosessi käynnissä (mukaan lukien tämä testi itse).
	if snap.ProcessCount <= 0 {
		t.Errorf("odotettiin ProcessCount > 0, saatiin %d", snap.ProcessCount)
	}

	// DiskUsedPct pitäisi olla järkevällä välillä 0-100.
	if snap.DiskUsedPct < 0 || snap.DiskUsedPct > 100 {
		t.Errorf("odotettiin DiskUsedPct välillä 0-100, saatiin %f", snap.DiskUsedPct)
	}

	// Timestamp ei saa olla tyhjä (nollaodottua aikaa) - se kertoisi
	// että time.Now()-kutsu epäonnistui jollain tavalla.
	if snap.Timestamp.IsZero() {
		t.Error("odotettiin Timestamp olevan asetettu, mutta se oli tyhjä")
	}
}

// TestScanCommonPorts varmistaa, että porttiskannaus palauttaa aina
// listan (mahdollisesti tyhjän) eikä koskaan kaadu tai jää jumiin.
func TestScanCommonPorts(t *testing.T) {
	ports := scanCommonPorts()

	if ports == nil {
		t.Error("odotettiin tyhjää listaa [], saatiin nil")
	}
}