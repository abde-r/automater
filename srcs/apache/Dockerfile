FROM phpmyadmin/phpmyadmin:latest

USER root

RUN apt-get update && apt-get install -y apache2 openssl && \
    a2enmod ssl rewrite && \
    mkdir -p /etc/apache2/ssl

COPY ssl/selfsigned.crt /etc/ssl/certs/selfsigned.crt
COPY ssl/selfsigned.key /etc/ssl/private/selfsigned.key

COPY default-ssl.conf /etc/apache2/sites-available/default-ssl.conf

RUN a2ensite default-ssl && \
    service apache2 restart

EXPOSE 443

CMD ["apache2ctl", "-D", "FOREGROUND"]
