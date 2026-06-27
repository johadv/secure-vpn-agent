# Etenemispäiväkirja

## Vaihe 0 — Runko
- Aloitettu: 2026-06-27
- Valmis: 
- Opin: 


## Vaihe 1 — node-a pystyssä
- Aloitettu: 2026-06-27
- Valmis: 2026-06-27
- Opin: UTM:n ARM64-VM vaatii arm64-ISO:n (ei amd64), UEFI Shell -kiertotien tarvittaessa, asennus saattaa kaatua satunnaisesti kernel/grub-paketin kohdalla mutta onnistuu uudelleenyrityksellä. node-a IP: 192.168.64.11

## Vaihe 1 — node-a pystyssä, node-b kesken
- Aloitettu: 2026-06-27
- Valmis (node-a): 2026-06-27
- node-a IP: 192.168.64.11, käyttäjätunnus: joni
- Opin: UTM:n ARM64-VM vaatii arm64-ISO:n (ei amd64); asennus voi kaatua satunnaisesti grub/kernel-paketin kohdalla, mutta toimii uudelleenyrityksellä; UEFI Shell -kiertotie tarpeen jos boot ei löydä ISO:a automaattisesti
- TODO seuraavaksi: node-b asennus (levytila/RAM loppui Macilta, harkittava pienempää levykokoa VM:lle tai VPS:ää node-b:n tilalle), sitten SSH-avainsiirtymä molemmille, sitten itse WireGuard-tunneli

## Vaihe 2 — Go-agentti, collector-paketti
- Aloitettu: 2026-06-27
- Valmis: 2026-06-27
- Opin: Go:n build constraints (//go:build linux / darwin) mahdollistavat käyttöjärjestelmäkohtaisen koodin samassa paketissa - Go valitsee oikean tiedoston automaattisesti käännösaikana. /proc toimii vain Linuxissa, macOS:lla pitää käyttää ps-komentoa saman tiedon hakemiseen. Testattu onnistuneesti Macilla: collector.Collect() palauttaa JSON-snapshotin prosessimäärästä.
