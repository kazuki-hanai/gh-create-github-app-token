CMD=gh-create-github-app-token
CMDPATH=main.go

run:
	go run $(CMDPATH) --debug=true

build:
	go build -o $(CMD) $(CMDPATH)

deps:
	go mod tidy
	go mod vendor

test:
	go test -v ./...
