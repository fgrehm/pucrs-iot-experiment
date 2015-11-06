IMAGE = fgrehm/alpine-go-web:1.5.1
TARGET = home-automation

source_go = $(shell find -L $(source) ! -path "src/*" -type f -name '*.go')

default: build

.PHONY: hack
hack:
	docker run -ti --rm -v `pwd`:/code $(IMAGE)

.PHONY: deploy
deploy: bin/home-automation-linux-arm
	scp bin/home-automation-linux-arm pi@10.32.143.201:~/home-automation

.PHONY: build
build: bin/home-automation

bin/home-automation: $(source_go)
	docker run -ti --rm -v `pwd`:/code $(IMAGE) gb build

bin/home-automation-linux-%: $(source_go)
	docker run -ti --rm -e GOOS=linux -e GOARCH=$(*) -v `pwd`:/code $(IMAGE) gb build

.PHONY: clean
clean:
	rm -rf bin/*
