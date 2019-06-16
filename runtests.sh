#!/bin/bash
set -e
golint db/db.go
golint web/web.go
golint route/route.go
golangci-lint run --deadline=2m
#gometalinter $(glide nv)  --disable-all --enable=safesql --enable=misspell --skip=db --skip=web
go test ./... -v -race -tags=integration
