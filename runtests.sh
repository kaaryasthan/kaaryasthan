#!/bin/bash
set -e
golint db/db.go
golint web/web.go
golint route/route.go
golangci-lint run --deadline=2m
go test ./... -v -race -tags=integration
