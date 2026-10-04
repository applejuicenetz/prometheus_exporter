# Prometheus Exporter

## Kompatibilität

- Go ersetzt `index.php`; nur Standardbibliothek, keine Laufzeitabhängigkeiten.
- Port `80`, `CORE_HOST`, `CORE_PORT`, `CORE_PASSWORD` und
  `PHP_SOCKET_TIMEOUT` erhalten. `LISTEN_ADDR` ist eine optionale Erweiterung.
- Ein 32 Byte langes Passwort bleibt unverändert; andere werden wie bisher
  mit MD5 gehasht. Keine Umstellung des Core-Authentifizierungsprotokolls.
- Nur unpräfixierte Attribute des ersten direkten `information`-Elements
  ausgeben, in XML-Reihenfolge. Werte nicht numerisch konvertieren oder runden.
- Keine zusätzlichen Metriken, `HELP` oder `TYPE`. Fehlerantworten und Logs
  dürfen den Passwortparameter nicht ausgeben.
- Fehlendes `information` ergibt wie PHP eine leere erfolgreiche Ausgabe.
  Ungültiges XML darf keine partiellen Metriken liefern.

- Bekannte Abweichung: PHP normalisiert rohe Zeilenumbrüche in Attributwerten zu
  Leerzeichen, Go nicht. Der Core sendet nur Zahlen, daher ohne Auswirkung.

## Prüfung

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test -race ./...
PHP_COMPATIBILITY=1 go test -run TestPHPCompatibility -v ./...
docker build -t applejuice-exporter:test .
```

Der optionale PHP-Vergleich braucht Docker und `php:8-apache`. Er vergleicht
die Ausgabe derselben XML-Eingaben mit dem ursprünglichen SimpleXML-Ausdruck.

Für einen lokalen Flatpak-Core benötigt der Container Host-Netzwerk, damit
`127.0.0.1:9851` den Core und nicht den Exporter-Container bezeichnet:

```sh
docker run --rm --network host -e CORE_HOST=127.0.0.1 \
  -e CORE_PORT=9851 -e CORE_PASSWORD= -e LISTEN_ADDR=:18383 \
  applejuice-exporter:test
```

Kein Core-Neustart und keine Änderung seiner Konfiguration für Exporter-Tests.
Bei Live-Vergleichen können sich Transferwerte zwischen zwei Abrufen ändern.

## Container und CI

- Multi-Stage-Build mit `CGO_ENABLED=0`, nativer Build-Plattform und
  Cross-Compilation über `TARGETOS`/`TARGETARCH`.
- Laufzeit: `gcr.io/distroless/static-debian13:nonroot`, UID 65532,
  Port 80. Standard-Docker erlaubt diesen Port auch ohne Root; bei explizit
  eingeschränkten Hosts `LISTEN_ADDR=:8080` und Port-Mapping anpassen.
- Healthcheck über `/exporter healthcheck`; keine Shell und kein `curl`.
  Er übernimmt `LISTEN_ADDR` aus derselben Container-Umgebung wie der Server.
  Wildcard-Bindungen werden auf IPv4-/IPv6-Loopback abgebildet; HTTP-Proxys
  werden für diese lokale Prüfung nicht verwendet.
- CI prüft Format, Vet, Race-Tests und PHP-Kompatibilität vor Image-Builds.
  Pull Requests bauen beide Architekturen, veröffentlichen aber nichts.
  Main-Pushes und manuelle Runs erhalten bestehende Registry-Ziele und Tags.
- Root-README ist deutschsprachige Nutzerdokumentation; Entwicklerhinweise
  gehören hierher.
