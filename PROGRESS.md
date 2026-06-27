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

## Vaihe 2 (laajennus) — portit ja levytila
- Lisätty: 2026-06-27
- Opin: scanCommonPorts tarkistaa AINOASTAAN omaa konetta (127.0.0.1), ei koskaan muita IP-osoitteita - tämä on tietoinen turvallisuussuunnittelu, ei vahinko. Muiden koneiden porttien skannaaminen luvatta olisi epäeettistä/laitonta, joten osoite on kovakoodattu vakioksi sen sijaan että se olisi muutettavissa parametrina. Käytetty rajattua listaa yleisistä porteista (SSH, HTTP/HTTPS, tietokannat) sen sijaan että skannattaisiin koko 1-65535 porttiavaruutta. syscall.Statfs antaa levytilatiedot suoraan käyttöjärjestelmältä, sama tieto jota "df"-komento käyttää. Testattu Macilla: löysi PostgreSQL:n (5432) paikallisesti auki olevana porttina.

## Vaihe 4 — CI/CD-putki (build, test, vet)
- Aloitettu: 2026-06-27
- Valmis: 2026-06-27
- Opin: GitHub Actions -workflow (.github/workflows/ci.yml) ajaa automaattisesti go build, go test ja go vet jokaisella pushilla main-haaraan. defaults/working-directory välttää toistuvan "cd agent" joka komennossa. setup-go:n riippuvuuscache vaatii go.sum-tiedoston - projektilla ei ole ulkoisia riippuvuuksia (vain stdlib), niin cache-varoitus on harmiton. Ensimmäinen workflow-ajo onnistui (success).

## Vaihe 4 (laajennus) — SAST (CodeQL)
- Lisätty: 2026-06-27
- Opin: CodeQL-koodiskannaus on ilmainen vain JULKISILLE repoille GitHub Free/Pro-tileillä - privaateissa repoissa se vaatii GitHub Code Security -lisenssin (Team/Enterprise). Muutin repon julkiseksi tämän takia. Workflow permissions piti olla "Read and write" CodeQL:n security-events-kirjoitusoikeutta varten (vaikka tämä ei itsessään ratkaisi alkuperäistä "repository not found" -virhettä - syy oli CodeQL+privaatti repo -yhdistelmä). Molemmat jobit (build-and-test, sast) menevät nyt läpi, tulokset näkyvät GitHubin Security-välilehdellä.

## Vaihe 4 valmis — koko CI/CD-putki pystyssä
- Valmis: 2026-06-27
- Yhteenveto: build+test+vet, CodeQL (SAST), Trivy (SCA), Dependabot - kaikki toimivat GitHubin Actions-välilehdellä. Repo julkinen.
