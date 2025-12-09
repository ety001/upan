#!/bin/bash
docker run -i --user 65534:65534 --rm -v /etc/php7:/etc/php7 -v /data/wwwroot:/data/wwwroot -v /tmp:/tmp -v /data/logs/php7:/var/log/php7 --network lnmp ety001/php:7.2.14 php /data/wwwroot/upan/artisan schedule:run
