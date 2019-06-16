build: build-web
	go generate
	go build


build-web:
	cd web && npm install && go generate
