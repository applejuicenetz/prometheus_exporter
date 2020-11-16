# appleJuice Prometheus Exporter

![](https://img.shields.io/github/license/applejuicenetz/prometheus_exporter.svg)
![](https://github.com/applejuicenetz/prometheus_exporter/workflows/docker/badge.svg)

[Prometheus](https://prometheus.io/) Metriken exporter für den appleJuice Client

## Prometheus

scraper config:

```yaml
scrape_configs:
    -   job_name: 'applejuice'
        scrape_interval: 60s
        scrape_timeout: 10s
        static_configs:
            -   targets: ['192.168.2.1:8383']
```

### Exposed Ports

- `80` - HTTP Port

### Environment Variables

| Variable                | Value                | Description                              |
|-------------------------|----------------------|------------------------------------------|
| `CORE_HOST`             | `http://192.168.2.1` | IP/HOST where Core is running, with http |
| `CORE_PORT`             | `9851`               | Core XML Port                            |
| `CORE_PASSWORD`         | `md5sum`             | `md5` hashed password                    |


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
version: '2.4'

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
