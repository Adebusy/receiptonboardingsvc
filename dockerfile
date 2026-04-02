#!/bin/bash
# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app


COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /cartservice

EXPOSE 8080

CMD [ "/receiptonboardingsvc" ] 