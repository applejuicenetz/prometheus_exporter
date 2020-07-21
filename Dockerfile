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

LABEL org.label-schema.name="appleJuice Prometheus Exporter" \
      org.label-schema.version=${BUILD_VERSION} \
      org.label-schema.vendor="appleJuiceNET" \
      org.label-schema.url="https://applejuicenet.de" \
      org.label-schema.build-date=${BUILD_DATE} \
      org.label-schema.vcs-ref=${VCS_REF} \
      org.label-schema.vcs-url="https://github.com/appleJuiceNET/prometheus_exporter" \
      org.label-schema.schema-version="1.0"
