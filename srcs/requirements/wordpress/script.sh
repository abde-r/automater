#! /bin/bash

cd /var/www/html/wordpress
curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar
chmod +x wp-cli.phar
mv wp-cli.phar /usr/local/bin/wp
wp core install --allow-root --url=${WP_URL} --title=${WP_TITLE} --admin_name=${WP_ADMIN_USER} --admin_password=${WP_ADMIN_PASS} --admin_email=${WP_ADMIN_EMAIL}

# wp plugin install redis-cache --activate --allow-root

wp plugin update --all --allow-root

wp theme install twentytwentytwo --activate --allow-root

sed -i 's/listen = \/run\/php\/php7.3-fpm.sock/listen = 9000/g' /etc/php/7.3/fpm/pool.d/www.conf

mkdir /run/php
# wp config set WP_REDIS_HOST redis --allow-root

# wp redis enable --allow-root