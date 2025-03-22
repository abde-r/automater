#! /bin/bash

docker rmi -f $(docker images -aq)

docker rm -f $(docker ps -aq)

sudo rm -rf "../data/db_data"

sudo rm -rf "../data/wp_data"
