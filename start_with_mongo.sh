#!/usr/bin/env bash

export PROVIDER=mongodb
export MONGO_CONNECTIONURI=mongodb://localhost
export USER_PASSWORD=complicate_pwd

go run main.go
