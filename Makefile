NAME=create-github-app-token

run:
	go run main.go

build:
	go build -o $(NAME) main.go
