.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -tags netgo -ldflags "-w -extldflags -static" -o bin/handler handler/main.go
	# env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handler handler/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
