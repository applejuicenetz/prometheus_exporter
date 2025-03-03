FROM php:8-apache

ENV CORE_HOST="" \
    CORE_PORT=9851 \
    CORE_PASSWORD=""

COPY index.php /var/www/html/metrics/

RUN echo "<a href=/metrics>/metrics</a>" > /var/www/html/index.html

EXPOSE 80

HEALTHCHECK --interval=60s --start-period=5s CMD curl -I --fail http://localhost:80 || exit 1

