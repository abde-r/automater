#!/bin/bash
set -e

export WP_CLI_PHP_ARGS='-d memory_limit=256M'

set -o allexport
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/.env"
set +o allexport

MACHINE_IP=$(curl -s http://checkip.amazonaws.com)
echo ${MACHINE_IP}

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
        --url="https://${MACHINE_IP}" \
        --title="$WP_TITLE" \
        --admin_user="$WP_ADMIN_USER" \
        --admin_password="$WP_ADMIN_PASS" \
        --admin_email="$WP_ADMIN_EMAIL" \
        --allow-root
    touch /var/www/html/.wp-installed
fi

wp option update siteurl "https://${MACHINE_IP}" --allow-root
wp option update home "https://${MACHINE_IP}" --allow-root

# Start Apache in the foreground to keep the container alive
exec apache2-foreground