FROM golang:1.13

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

ARG DEBIAN_FRONTEND=noninteractive

# Update and install curl
RUN apt-get update

# Creating work directory
RUN mkdir /code

# Adding project to work directory
ADD . /code

# Choosing work directory
WORKDIR /code

# build project
RUN go build -o movie.night.api.server .

EXPOSE 9002

CMD ["./movie.night.api.server"]