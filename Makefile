# **************************************************************************** #
#                                                                              #
#                                                         :::      ::::::::    #
#    Makefile                                           :+:      :+:    :+:    #
#                                                     +:+ +:+         +:+      #
#    By: ael-asri <ael-asri@student.42.fr>          +#+  +:+       +#+         #
#                                                 +#+#+#+#+#+   +#+            #
#    Created: 2023/01/11 13:30:38 by ael-asri          #+#    #+#              #
#    Updated: 2023/01/11 13:30:39 by ael-asri         ###   ########.fr        #
#                                                                              #
# **************************************************************************** #

NAME = inception

all:$(NAME)

$(NAME):
	sudo mkdir -p "./data/db_data"
	sudo mkdir -p "./data/wp_data"
	docker compose -f ./srcs/docker-compose.yml up --build -d

down:
	docker compose -f ./srcs/docker-compose.yml down
