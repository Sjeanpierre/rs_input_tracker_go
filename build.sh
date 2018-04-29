#!/usr/bin/env bash

GOOS=linux go build -o bin/input_tracker_app
docker build -t input_tracker .
docker run --rm -it --env-file .env input_tracker