#!/usr/bin/env bash

export POSTGRES_CONNECTIONURI=postgres://conkeys:S0jeje1!@localhost/conkeys?sslmode=disable

export ADMIN_PASSWORD=complicate_pwd

go run main.go
