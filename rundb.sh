#!/usr/bin/env bash

CID=$(docker container ls --filter=name=pg_conkeys --format='{{ .ID }}')
[ -z "$CID" ] || docker stop $CID && docker rm $CID

docker run --name pg_conkeys -p 5432:5432 -e POSTGRES_USER=conkeys -e POSTGRES_PASSWORD=S0jeje1! -d postgres
