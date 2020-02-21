GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/hooky

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

build:
	export GOPATH=$HOME/go/
	go get -v  github.com/NightStr/goHook/hookBot
	go get -v github.com/NightStr/goHook/middleware
	go build -o bin/hooky

heroku: $(DOCKER_CMD)
	heroku container:push web
