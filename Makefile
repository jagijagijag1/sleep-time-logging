.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/bedtime-delay bedtime-delay/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wakeup-delay wakeup-delay/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
