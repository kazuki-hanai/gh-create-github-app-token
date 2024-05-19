CMD=create-github-app-token
CMDPATH=cmd/create-github-app-token/main.go

run:
	go run $(CMDPATH) --debug=true

build:
	go build -o $(CMD) $(CMDPATH)
