#!/bin/bash
set -e

export WP_CLI_PHP_ARGS='-d memory_limit=256M'

# Only install WordPress once
if [ ! -f /var/www/html/.wp-installed ]; then
    wp core download --allow-root --path="/var/www/html" --force
    wp core config \
        --dbname="$WORDPRESS_DB_NAME" \
        --dbuser="$WORDPRESS_DB_USER" \
        --dbpass="$WORDPRESS_DB_PASSWORD" \
        --dbhost="$WORDPRESS_DB_HOST" \
        --allow-root --skip-check
    wp core install \
        --url="https://localhost:9443" \
        --title="cloud-1" \
        --admin_user=user \
        --admin_password=pass \
        --admin_email=user@example.com \
        --allow-root
    touch /var/www/html/.wp-installed
fi

wp option update siteurl "https://localhost" --allow-root
wp option update home "https://localhost" --allow-root

# Start Apache in the foreground to keep the container alive
exec apache2-foreground
