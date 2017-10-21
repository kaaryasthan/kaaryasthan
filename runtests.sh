#!/bin/bash
set -e
golint db/db.go
golint web/web.go
gometalinter $(glide nv)  --disable-all --enable=errcheck --enable=ineffassign \
	--enable=gofmt --enable=vet --enable=deadcode --enable=varcheck \
	--enable=structcheck --enable=maligned --enable=unconvert --enable=unused \
	--enable=goconst --enable=gas --enable=unparam --enable=staticcheck \
	--enable=interfacer --enable=vetshadow --enable=megacheck --enable=golint \
	--skip=db --skip=web
#gometalinter $(glide nv)  --disable-all --enable=safesql --enable=misspell --skip=db --skip=web
go test $(glide nv) -v -race
