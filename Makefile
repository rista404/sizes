build:
	go build -ldflags="-s -w" -o sizes cmd/main.go

set_global:
	cp ./sizes /usr/local/bin