#/bin/bash
go get github.com/pilu/fresh
go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
go get github.com/alecthomas/gometalinter
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.16.0
npm install -g @angular/cli