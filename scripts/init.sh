#!/bin/bash

cp .env.sample .env

docker network create wiser

docker compose up -d --build
