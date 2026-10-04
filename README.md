# appleJuice Prometheus Exporter

![](https://img.shields.io/github/license/applejuicenetz/prometheus_exporter.svg)
![](https://github.com/applejuicenetz/prometheus_exporter/workflows/container/badge.svg)

[Prometheus](https://prometheus.io/) Metriken exporter für den appleJuice Client

Go ersetzt PHP/Apache. Bestehende Container-Konfiguration, Port `80` und
Metrikausgabe bleiben erhalten: `applejuice_<Attributname> <Attributwert>`
in der Reihenfolge der XML-Attribute, ohne zusätzliche Metriken, `HELP` oder `TYPE`.
Abruf über `/metrics`, `/metrics/` oder `/metrics/index.php`.

## Prometheus

scraper config:

```yaml
scrape_configs:
  - job_name: 'applejuice'
    scrape_interval: 60s
    scrape_timeout: 10s
    static_configs:
      - targets: [ '192.168.2.1:8383' ]
```

### Exposed Ports

- `80` - HTTP Port (Container läuft als `nonroot` auf Basis von `gcr.io/distroless/static-debian13`)

### Environment Variables

| Variable        | Value                | Description                    |
|-----------------|----------------------|--------------------------------|
| `CORE_HOST`     | `http://192.168.2.1` | IP/HOST where Core is running  |
| `CORE_PORT`     | `9851`               | Core XML Port                  |
| `CORE_PASSWORD` | `md5sum` or plain    | `md5` hashed or plain password |
| `PHP_SOCKET_TIMEOUT` | `10`            | optional, Timeout in Sekunden für den Core-Abruf |
| `LISTEN_ADDR`   | `:80`                | optional, Listen-Adresse des Exporters |

## Docker

### docker run

create and run `applejuice2prometheus` container with the following command

```bash
docker run -d \
        -p 8383:80 \
        -e "CORE_HOST=http://192.168.1.2" \
        -e "CORE_PORT=9851" \
        -e "CORE_PASSWORD=de305845b091d971732a123977e2d816" \
        --name applejuice2prometheus \
        ghcr.io/applejuicenetz/prometheus_exporter:latest
```

### docker-compose.yml

```yaml
services:
  applejuice2prometheus:
    image: ghcr.io/applejuicenetz/prometheus_exporter:latest
    restart: always
    container_name: applejuice2prometheus
    mem_limit: 32MB
    network_mode: bridge
    ports:
      - 8383:80/tcp
    environment:
      CORE_HOST: http://192.168.1.2
      CORE_PORT: 9851
      CORE_PASSWORD: de305845b091d971732a123977e2d816
```
