#!/usr/bin/env bash

#export PROVIDER=mongodb
#export MONGO_CONNECTIONURI=mongodb://localhost

export PROVIDER=postgres
export MONGO_CONNECTIONURI=postgres://conkeys:S0jeje1!@localhost/conkeys?sslmode=disable

export USER_PASSWORD=complicate_pwd

go run main.go
