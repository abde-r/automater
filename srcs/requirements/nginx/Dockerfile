FROM debian:buster

RUN apt-get -y  update && apt-get -y upgrade &&\
	apt-get -y install nginx && apt-get install openssl

RUN openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=MA/ST=Denial/L=Springfield/O=Dis/CN=www.ael-asri.com" \
    -keyout /etc/ssl/private/ael-asri.key  -out /etc/ssl/certs/ael-asri.crt

COPY conf/default /etc/nginx/sites-available/default

CMD ["nginx", "-g", "daemon off;"]
