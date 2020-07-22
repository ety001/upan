FROM node:10-alpine as build-stage
WORKDIR /app
ADD . /app
RUN npm install && npm run production && rm -rf node_modules

FROM trafex/alpine-nginx-php7
WORKDIR /var/www/html
COPY --chown=nobody --from=build-stage /app /var/www/html
COPY docker-config/nginx/nginx.conf /etc/nginx/nginx.conf
# COPY docker-config/php7 /etc/
COPY --from=composer /usr/bin/composer /usr/bin/composer

VOLUME /var/www/html/storage/app

USER root
RUN apk add --no-cache php7-xmlwriter php7-tokenizer php7-pdo \
    php7-pdo_mysql php7-fileinfo

USER nobody
RUN composer install --optimize-autoloader --no-dev --no-interaction --no-progress && \
    php artisan storage:link
