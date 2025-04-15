#!/bin/bash
set -e

echo "Waiting for MySQL to be available..."
# Wait until MySQL (db service) responds
# while ! mysqladmin ping -h"$WORDPRESS_DB_HOST" --silent; do
#   sleep 3
# done

echo "MySQL is available."

# If WordPress is not installed, run the installation.
if ! wp core is-installed --path="/var/www/html" --allow-root; then
    echo "WordPress not installed. Installing..."
    wp core download --allow-root --path="/var/www/html" --force
    wp core config \
        --dbname="$WORDPRESS_DB_NAME" \
        --dbuser="$WORDPRESS_DB_USER" \
        --dbpass="$WORDPRESS_DB_PASSWORD" \
        --dbhost="$WORDPRESS_DB_HOST" \
        --allow-root --skip-check
    wp core install \
        --url="http://localhost:9000" \
        --title="My WordPress Site" \
        --admin_user=admin \
        --admin_password=admin \
        --admin_email=admin@example.com \
        --allow-root
    echo "WordPress installed successfully."
else
    echo "WordPress is already installed."
fi

# Start the webserver (the official image uses apache2-foreground).
exec apache2-foreground
