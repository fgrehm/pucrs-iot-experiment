IMAGE = fgrehm/rpi-arduino-buttons
TARGET = app

source_go = $(shell find -L "src/" -type f -name '*.go')
bindata = src/cmd/app/assets.go

default: build

.PHONY: hack
hack: docker.build
	@mkdir -p .docker-dev/gradle
	docker run -ti --rm -v `pwd`:/code -v `pwd`/.docker-dev/gradle:/home/developer/.gradle --privileged $(IMAGE)

.PHONY: serve
serve: docker.build build/app
	docker run -ti --rm -p 35729:35729 -p 8080:8080 -v `pwd`/build:/code $(IMAGE) ./app

.PHONY: deploy.rpi
deploy.rpi: build/app-arm
	scp $< pi@10.32.143.201:~/buttons-app

.PHONY: build
build: build/app

.PHONY: build.client
build.client:
	docker run -ti --rm -v `pwd`:/code $(IMAGE) sh -c 'cd client && make'

.PHONY: build.mobile
build.mobile: build/android.apk

.PHONY: docker.build
docker.build:
	docker build -t $(IMAGE) .

.PHONY: clean
clean:
	rm -rf bin/* build/*

$(bindata): docker.build build.client
	docker run -ti --rm -v `pwd`:/code $(IMAGE) sh -c 'cd client && go-bindata-assetfs -nomemcopy www/...'
	mv client/bindata_assetfs.go $(@)

build/app: $(source_go) $(bindata)
	docker run -ti --rm -v `pwd`:/code $(IMAGE) gb build cmd/app
	mv bin/app $(@)

build/app-arm: $(source_go)
	docker run -ti --rm -e GOOS=linux -e GOARCH=arm -v `pwd`:/code $(IMAGE) gb build cmd/app
	mv bin/app-linux-arm $(@)

build/android.apk:
	@mkdir -p .docker-dev/gradle
	docker run -ti --rm -v `pwd`:/code -v `pwd`/.docker-dev/gradle:/home/developer/.gradle $(IMAGE) sh -c 'cd client && make build.android'
