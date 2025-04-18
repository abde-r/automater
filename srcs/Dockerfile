FROM wordpress:6.7-php8.3-apache

# Install WP-CLI
RUN curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x wp-cli.phar && \
    mv wp-cli.phar /usr/local/bin/wp

# Increase PHP memory limit
RUN echo "memory_limit = 256M" > /usr/local/etc/php/conf.d/memory.ini

COPY apache/ssl/selfsigned.crt /etc/ssl/certs/selfsigned.crt
COPY apache/ssl/selfsigned.key /etc/ssl/private/selfsigned.key

COPY apache/default-ssl.conf /etc/apache2/sites-available/default-ssl.conf

# Enable SSL modules and site (no reload/start!)
RUN a2enmod ssl && a2ensite default-ssl
