CREATE DATABASE wordpress;

CREATE USER 'stronk'@'%' IDENTIFIED BY 'pass';
GRANT ALL PRIVILEGES ON wordpress.* TO 'stronk'@'%';
FLUSH PRIVILEGES;

CREATE USER 'stronkest'@'%' IDENTIFIED BY 'stronkpass';
GRANT SELECT ON wordpress.* TO 'stronkest'@'%';
FLUSH PRIVILEGES;

ALTER USER 'root'@'localhost' IDENTIFIED BY 'pass1234';
FLUSH PRIVILEGES;