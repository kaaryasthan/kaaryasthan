build: build-web
	go generate
	go build


build-web:
	cd web && npm install && go generate

test: lint
	go test ./... -v -race -tags=integration

lint:
	golint db/db.go
	golint web/web.go
	golint route/route.go
	golangci-lint run --deadline=2m
