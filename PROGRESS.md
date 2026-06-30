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

## Vaihe 3 — mTLS toimii paikallisesti (Macilla, server+agent samalla koneella)

- Aloitettu: 2026-06-28
- Valmis (paikallinen testaus): 2026-06-28
- Opin: Nykyaikaiset TLS-kirjastot (Go mukaan lukien) eivät enää luota sertifikaatin CN-kenttään osoitteen validoinnissa - SAN (Subject Alternative Name) -kenttä on pakollinen. Ensimmäinen sertifikaattigenerointi epäonnistui juuri tästä syystä ("doesn't contain any IP SANs"), korjattu lisäämällä -extfile openssl-komentoon. mTLS-yhteys toimii molempiin suuntiin: palvelin vaatii ja varmistaa client-sertifikaatin (RequireAndVerifyClientCert), client varmistaa palvelimen sertifikaatin CA:ta vasten. Testattu onnistuneesti: agentti lähetti JSON-snapshotin palvelimelle TLS 1.3:n yli, palvelin tunnisti clientin CN:n (node-agent).
- Seuraava: viedä tämä node-a/node-b-väliseksi WireGuard-tunnelin yli kun node-b on asennettu

## Vaihe 3 (laajennus) — mTLS-turvatestit

- Testattu: 2026-06-28
- Testi 1: yhteys ilman client-sertifikaattia -> palvelin hylkäsi (TLS alert: certificate required)
- Testi 2: yhteys CA:n allekirjoittamattomalla, itse luodulla sertifikaatilla -> palvelin hylkäsi (TLS alert: unknown ca)
- Opin: RequireAndVerifyClientCert tarkistaa kahdessa vaiheessa - ensin ETTÄ sertifikaatti esitetään ylipäätään, sitten KUKA sen on allekirjoittanut. Kahden testin avulla todistettu molemmat suojakerrokset toimivat erikseen. Tämä on konkreettinen, esittelykelpoinen todiste mTLS:n toimivuudesta CV:tä/haastattelua varten.

## Vaihe 1 — node-b asennettu, molemmat näkevät toisensa

- Valmis: 2026-06-30
- node-a IP: 192.168.64.11
- node-b IP: 192.168.64.17
- Opin: Ubuntu-ISO-latauksissa kannattaa AINA varmistaa SHA256-checksum ennen VM-asennusta - vanha/väärä ISO-versio (26.04 sekoittui 24.04:n kanssa) aiheutti toistuvan asennuskaatumisen joka näytti muistiongelmalta mutta oli oikeasti väärä/vioittunut tiedosto. cdimage.ubuntu.com:n suorat URL:t muuttuvat point release -numeron mukana (24.04 -> 24.04.4) - vanha linkki ilman tarkkaa versionumeroa johtaa 404:ään. Ping-testi ennen WireGuardia kannattaa aina tehdä ensin - varmistaa peruskonnektiviteetti ennen monimutkaisempaa konfiguraatiota.

## Vaihe 1 valmis — WireGuard-tunneli toimii node-a ja node-b välillä

- Valmis: 2026-06-30
- node-a: 10.10.0.2 (WireGuard-sisäinen), 192.168.64.11 (LAN)
- node-b: 10.10.0.1 (WireGuard-sisäinen), 192.168.64.17 (LAN), kuuntelee porttia 51820
- Opin: apt-asennus saattoi epäonnistua UTM:n virtuaaliverkossa IPv6-reitityksen takia ("Network is unreachable") - korjattu pakottamalla apt käyttämään IPv4:ää (-o Acquire::ForceIPv4=true). Avainten generointi täytyy tehdä juuri /etc/wireguard-kansiossa rootina, muuten Permission denied. wg-quick up wg0 tekee neljä asiaa automaattisesti: luo virtuaalisen verkkokortin, asettaa avaimet, antaa IP-osoitteen, ja nostaa rajapinnan ylös. Tunneli vahvistettu toimivaksi: wg show näyttää "latest handshake" molemmilla koneilla, ping 10.10.0.x onnistui 0% pakettihäviöllä.
- Seuraava: vie mTLS-yhteys (server+agent) tämän WireGuard-tunnelin sisään, 127.0.0.1-osoitteiden sijaan käytä 10.10.0.1/10.10.0.2

## MERKKIPAALU — koko ketju toimii päästä päähän, node-a -> WireGuard -> mTLS -> node-b

- Valmis: 2026-06-30
- node-a agentti lähetti onnistuneesti JSON-snapshotin WireGuard-tunnelin (10.10.0.2 -> 10.10.0.1) läpi, mTLS:n suojaamana, node-b:n palvelimelle.
- Palvelin vahvisti clientin identiteetin sertifikaatista (CN=node-agent) ja vastaanotti datan onnistuneesti.
- Macin rajallisista resursseista huolimatta (levytila/muisti tiukoilla koko projektin ajan) sain silti rakennettua toimivan yhteyden kahden virtuaalikoneen välille SSH:n kautta - tässä opin käytännössä miten VPN, verkko ja WireGuard toimivat yhdessä, en vain teoriassa vaan itse pystyttämällä ja debugaamalla koko ketjun alusta loppuun.
- Opin matkalla: SSH/scp saattaa kokea hetkellisiä "No route to host" -katkoja UTM:n virtuaaliverkossa vaikka ping toimii - kannattaa yrittää uudelleen ennen syvempää debuggausta. PATH-muutokset .profile-tiedostoon eivät periydy automaattisesti uusiin terminaali-istuntoihin samalla käyttäjällä - source ~/.profile tarvitaan tai uusi kirjautuminen. Sertifikaattien siirto eri koneille onnistuu turvallisesti scp:llä suoraan, koska certs/ on aina .gitignoressa eikä koskaan kulje Git-historian kautta.
