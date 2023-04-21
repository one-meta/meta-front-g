build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o meta-front-g main.go
