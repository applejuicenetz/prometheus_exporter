FROM php:7-apache

ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

ENV CORE_HOST ""
ENV CORE_PORT 9851
ENV CORE_PASSWORD ""

COPY index.php /var/www/html/metrics/

RUN echo "<a href=/metrics>/metrics</a>" > /var/www/html/index.html

EXPOSE 80

HEALTHCHECK --interval=60s --start-period=5s CMD curl -I --fail http://localhost:80 || exit 1

LABEL org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.vendor="appleJuiceNET" \
      org.opencontainers.image.url="https://applejuicenet.cc" \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${VCS_REF} \
      org.opencontainers.image.source="https://github.com/applejuicenetz/prometheus_exporter"
