FROM alpine:latest

ADD build/docker/ci-kubernetes /app/

WORKDIR /app
