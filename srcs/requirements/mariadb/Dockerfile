FROM debian:buster

RUN apt-get -y update && apt-get -y upgrade &&\
	apt-get -y install default-mysql-server

COPY tools/initdb.sql /

COPY tools/50-server.cnf /etc/mysql/mariadb.conf.d/

RUN service mysql start && mysql < initdb.sql

CMD [ "mysqld" ]
