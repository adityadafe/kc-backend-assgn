.PHONY: run build swagger test 
build:
	go build -o bin/app cmd/main.go

run:build
	bin/app

installswag:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	swag init -g cmd/main.go

test:
	go test ./internal/process

dockerbuildxcreate:
	docker buildx create --use --name kc-backend-builder

dockerbuild:dockerbuildxcreate
	docker buildx build --platform linux/amd64,linux/arm64 -t $(USERNAME)/kc-backend:latest --push .
